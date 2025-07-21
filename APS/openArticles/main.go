package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	functions "aps/openArticles/functions"
	generals "aps/generalFunctions"
	"time"
)


func getPageData(browser_path, user_data_dir, mainDir, url string) (data []functions.ArticleLink, OldData bool) {
	defer func() {
		if r := recover(); r != nil {

			log.Fatal("Error in gettings the page data in Nature :", r)

		}
	}()
	defer func() {
		if err := generals.Delete_temp_files(mainDir); err != nil {
		}
	}()
	browser, err := generals.Open_browser(browser_path, user_data_dir)
	defer browser.MustClose()
	if err != nil {
		log.Fatal("Could not Open the Browser, error is :", err)

	}
	
	page, err := generals.Open_page(browser, url)
	defer page.MustClose()
	if err != nil {
		log.Fatal("Could not Open the Page trying agian")

	}
	start_time := time.Now()
	

	data, oldData := functions.Get_data(page, url)

	min, max := 2, 6
	// Generate a random number between min and max (inclusive)
	randomNumber := rand.Intn(max-min+1) + min

	if crawlTime := time.Since(start_time); crawlTime < time.Duration(randomNumber)*time.Second {
		fmt.Println("crawl time ", crawlTime)
		time.Sleep(time.Duration(randomNumber)*time.Second - crawlTime)
	}
	if oldData {

		return data, true
	}
	return data, false

}

func routine(browser_path, mainDir, filesDir, webPage string, index,userDataindex int, doneChan chan bool) {
	defer func() {
		if r := recover(); r != nil {

			log.Fatal("Error in gettings the Article Links in :", r)

		}
	}()

	defer func(){
		doneChan <- true
	}()

	
	var all_data []functions.ArticleLink
	for i:= 1 ;i <= 10000; i++ {
		user_data_dir1 := filepath.Join(mainDir, "user-dir", fmt.Sprintf("%d_%d_%d", index, i,userDataindex))
		os.Mkdir(user_data_dir1, 0777)

		url := webPage + fmt.Sprintf("?page=%d", i)
		data, OldData := getPageData(browser_path, user_data_dir1, mainDir, url)

		all_data = append(all_data, data...)
		if OldData {
			log.Println("Other remaining data was old")
			break
		}
		log.Println("Total number of links are for routine:",userDataindex + 1, len(all_data))
		log.Println("we are in page ",i, "in url of ",webPage)
		if err := generals.Delete_temp_files(mainDir); err != nil {
		}

	}

	err := functions.Save_data(filesDir, all_data)
	if err != nil {
		log.Fatal("Saving data failed, error is:", err)
	}
}

func main() {
	var numberOfRoutines = 1

	defer func() {
		if r := recover(); r != nil {

			log.Fatal("Error in gettings the Article Links  :", r)

		}
	}()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)

	}
	mainDir := filepath.Join(currentDir, "..")
	journalsDir := filepath.Join(mainDir, "files")
	filesDir := filepath.Join(mainDir, "files")
	
	urls, startIndex, endIndex, err := functions.GetJournalLinks(filesDir, journalsDir)
	if err != nil {
		log.Fatal("Could not get Issues links, error is :", err)

	}

	browser_path := generals.Get_browser_path(mainDir)

	for index := 0; index < len(urls); index += numberOfRoutines {
		if index + numberOfRoutines > len(urls){
			numberOfRoutines = numberOfRoutines - 1
		}
		doneChan := make(chan bool, numberOfRoutines)
		for i := 0; i < numberOfRoutines; i++ {
			// if index + i > len(urls){

			// }
			webPage := urls[index+i]
			go routine(browser_path, mainDir, filesDir, webPage, index,i,doneChan)
		}
		
		for j := 0; j < numberOfRoutines; j++ {

			<-doneChan
			log.Println("checking the routine number", j+1)
		}
		if err := generals.Delete_temp_files(mainDir); err != nil {
		}
		err := functions.Update_slices(startIndex+index+numberOfRoutines, endIndex, filesDir)
		if err != nil {
			fmt.Println("error in updating slice")
		}

	}
	generals.ClearScreen()
	log.Println("Got All the Article links")

}

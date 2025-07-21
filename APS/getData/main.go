package main

import (
	"fmt"

	functions "aps/getData/functions"
	general "aps/generalFunctions"

	"os"

	"math/rand"
	"path/filepath"

	"time"

	"log"
)

func routine(url, browser_path, user_data_dir, files_dir string, doneChan chan bool) {
	start_time := time.Now()
	
	defer func() {
		if r := recover(); r != nil {
			log.Println("Program encountered an error and recovered:", r)
		}
	}()

	browser, err := general.Open_browser(browser_path, user_data_dir)
	defer browser.MustClose()
	if err != nil {
		fmt.Println("Error Oppening the browser: ", err)
		return
	}
	done := make(chan bool)
	
	go general.PrintActiveTimer(start_time, done)
	page, err := general.Open_page(browser, url)
	defer page.MustClose()
	defer func() {
		done <- true
		doneChan <- true
	}()
	if err != nil {
		fmt.Println("Error Oppening the page : ", err)

	}
	all_data, err := functions.Get_data(page, url)
	if err != nil {
		fmt.Println("erRor in getting data 1")
	} else {
		if err := functions.Save_data(files_dir, all_data); err != nil {
			fmt.Println("Error save file :", err)
		}
	}
	fmt.Println("Got all the data of ")
	// Define the range
	min, max := 3, 8
	// Generate a random number between min and max (inclusive)
	randomNumber := rand.Intn(max-min+1) + min

	if crawlTime := time.Since(start_time); crawlTime < time.Duration(randomNumber)*time.Second {
		fmt.Println("crawl time ", crawlTime)
		time.Sleep(time.Duration(randomNumber)*time.Second - crawlTime)
	}

}

func main() {
	

	const numberOfRoutines = 2
	fmt.Println("Starting the crawler ")
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("Program encountered an error in Getting Nature Data:", r)
			if err := general.Delete_temp_files(currentDir); err != nil {
				log.Println("deleted remain temps")
			}
		}
	}()
	if err := general.Delete_temp_files(currentDir); err != nil {
		fmt.Println("Error:", err)
	}

	
	mainDir := filepath.Join(currentDir, "..")
	filesDir := filepath.Join(mainDir, "files")
	browser_path := general.Get_browser_path(mainDir)
	urls, start_slice, end_slice, err := functions.GetArticleLinks(filesDir)
	if err != nil {
		fmt.Println("Error getting urls:", err)
		return
	}
	

	for index := 0; index < len(urls); index += numberOfRoutines {

		doneChan := make(chan bool, numberOfRoutines)
		for i := 0; i < numberOfRoutines; i++ {
			log.Println("starting routine",i)
			user_data_dir := filepath.Join(currentDir, "user-dir", fmt.Sprint(index+i))
			os.Mkdir(user_data_dir, 0777)
			go routine(urls[index+i], browser_path, user_data_dir, filesDir, doneChan)

		}

		for j := 0; j < numberOfRoutines; j++ {
			<-doneChan
		}
	

		if err := general.Delete_temp_files(currentDir); err != nil {
			log.Println("deleted some temps")
		}
		err = functions.Update_slices(start_slice+index+numberOfRoutines, end_slice, filesDir)
		if err != nil {
			fmt.Println("error in updating slice")
		}
		general.ClearScreen()
		fmt.Printf("\nCrawl Count Is %v", index+numberOfRoutines)

	}
	log.Println("Web crawling has been finished")
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	functions "aps/openJournals/functions"
	general "aps/generalFunctions"
)

func routine(url, mainDir, browserPath, filesDir string)(journalsCounts int) {
	
	log.Println("Getting journal links from ", url)
	user_data_dir := filepath.Join(mainDir, "user-dir", fmt.Sprint(time.Now().Format("20060102_150405_000000")))
	os.Mkdir(user_data_dir, 0777)
	browser, err := general.Open_browser(browserPath, user_data_dir)
	defer browser.MustClose()
	if err != nil {
		log.Fatal("Could not Open the Browser, error is :", err)

	}
	page, err := general.Open_page(browser, url)
	defer page.MustClose()

	if err != nil {
		log.Fatal("Could not Open the Page, error is:", err)

	}
	data := functions.Get_data(page, url)
	
	err = functions.Save_data(filesDir, data)
	if err != nil {
		log.Fatal("Saving data failed, error is:", err)

	}
	journalsCounts = len(data)
	return

}

func main() {
	var urls = []string{
		"https://target-website",

	}
	

	defer func() {
		if r := recover(); r != nil {

			log.Fatal("Error in gettings the Journal Links :", r)

		}
	}()
	currentDir, err := os.Getwd()
	
	if err != nil {
		log.Fatal("Error getting current directory:", err)

	}
	mainDir := filepath.Join(currentDir, "..")
	defer func() {
		if err := general.Delete_temp_files(mainDir); err != nil {
			log.Println("Some temp files from user-dir were open and didn't deleted")
		}
	}()
	filesDir := filepath.Join(mainDir, "files")
	functions.CheckIfFileExistsAndDelete(mainDir)
	browserPath := general.Get_browser_path(mainDir)
	journalCounts := 0
	for _, url := range urls {
 
		journalCounts += routine(url, mainDir, browserPath, filesDir)
		general.ClearScreen()
	}
	log.Println("number of journals are :",journalCounts)
}

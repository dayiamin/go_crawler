package generalFunctions

import (
	"fmt"
	"log"
	"time"
	"bufio"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"os"
	"path/filepath"
)

// takes browser and url, opens a new page
func Open_page(browser *rod.Browser, url string) (page *rod.Page, err error) {
	page, err = stealth.Page(browser) // Replace with your target URL
	log.Println(url)

	custom_user_agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: custom_user_agent})
	if err != nil {
		fmt.Println("Failed to Open:", err)
		return
	}
	err = page.Navigate(url)
	page.Mouse.MustScroll(300.1,200.2)
	if err != nil {
		fmt.Println("Failed to navigate:", err)
		return
	}
	return
}

// Open a browser with given browser path , store data in user data path
func Open_browser(browser_path, user_data_path string) (browser *rod.Browser, err error) {
	log.Println("oppening the browser")
	new_launcher := launcher.New().HeadlessNew(true).Bin(browser_path).UserDataDir(user_data_path)
	url, err := new_launcher.Launch()
	if err != nil {
		// Create a new browser instance with a default client that sets headers
		launcher := launcher.New().HeadlessNew(true).UserDataDir(user_data_path)
		url = launcher.MustLaunch()
	}
	browser = rod.New().ControlURL(url).Timeout(3 * time.Minute)
	err = browser.Connect()
	return browser, err
}





// extract browser path from files/browser_path.txt
func Get_browser_path(mainDir string) (path string) {
	filename := filepath.Join(mainDir,"files", "browser_path.txt")
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error Browser Path does not exist in files folder")
		return
	}
	defer file.Close() // Ensure the file is closed when the function ends

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line
	if scanner.Scan() {
		path = scanner.Text() // Get the text of the line

	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}
	return
}

// Extract urls from files/article_links.jsonl

// Try to delete tempory files from user-dir
func Delete_temp_files(mainDir string) error {
	mainDir = filepath.Join(mainDir, "user-dir")
	entries, err := os.ReadDir(mainDir)
	if err != nil {
		return err
	}
	// Loop through each item in the directory
	for _, entry := range entries {
		// Get the full path of the item
		path := filepath.Join(mainDir, entry.Name())

		// Check if it is a directory
		if entry.IsDir() {
			// Remove the directory and its contents
			if err := os.RemoveAll(path); err != nil {
				
			}
		} else {
			// If it is a file, remove it
			if err := os.Remove(path); err != nil {
				
			}
		}
	}
	return nil
}

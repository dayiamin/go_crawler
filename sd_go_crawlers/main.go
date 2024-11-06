package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type ArticleData struct {
	Title       string   `json:"title"`
	Publisher   string   `json:"publisher"`
	DOI         string   `json:"doi"`
	Website     string   `json:"website"`
	PublishDate string   `json:"publish_date"`
	Journal     string   `json:"journal"`
	Authors     []string `json:"authors"`
	Keywords    []string `json:"keywords"`
	Abstract    string   `json:"abstract"`
	References  []string `json:"references"`
}

func save_data(file_path string, data ArticleData) (err error) {
	filename := filepath.Join(file_path, "article_data.jsonl")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the JSON data to the file followed by a newline
	_, err = file.WriteString(string(jsonData) + "\n")
	return err
}

func fixAuthors(authors []string) (fullNames []string) {

	tempName := ""

	for _, namePart := range authors {
		namePart = strings.TrimSpace(namePart) // Trim whitespace
		if namePart != "" {                    // Check if not empty
			if tempName != "" { // We have a first name
				fullName := tempName + " " + namePart
				fullNames = append(fullNames, fullName)
				tempName = "" // Reset tempName
			} else {
				tempName = namePart // Store the first name part
			}
		}
	}
	return
}

func get_ref(page *rod.Page) (all_references []string) {
	bibReferences, err := page.Elements(`li.bib-reference.u-margin-s-bottom`)
	if err == nil {
		for _, ref := range bibReferences {
			text, _ := ref.Text()
			all_references = append(all_references, text)
		}
		return
	}
	// 2. Select all spans with the class "reference"
	spanReferences, err := page.Elements(`span.reference`)
	if err == nil {
		for _, span := range spanReferences {
			text, _ := span.Text()
			all_references = append(all_references, text)
		}

		return
	}
	return
}

func get_data(page *rod.Page, url string) (articleData ArticleData, err error) {
	fmt.Println("Getting Data")
	time_out := 30 * time.Second
	if title, err := page.Timeout(time_out).Element(`span.title-text`); err == nil {
		articleData.Title = title.MustText()
	} else {
		page.Reload()
	}
	// Journal
	if journal, err := page.Timeout(time_out).Element(`meta[name="citation_journal_title"]`); err == nil {
		articleData.Journal = *journal.MustAttribute("content")
	}
	// Authors
	author_elements, err := page.Elements(`span.react-xocs-alternative-link`)
	if err == nil {
		for _, el := range author_elements {
			if author, err := el.Text(); err == nil {
				articleData.Authors = append(articleData.Authors, author)
			}
		}
		articleData.Authors = fixAuthors(articleData.Authors)
	}

	// DOI
	if doi, err := page.Element(`meta[name="citation_doi"]`); err == nil {
		articleData.DOI = *doi.MustAttribute("content")
	}
	// Keywords
	keyword_elements, err := page.Elements(`div.keywords-section div`)
	if err == nil {
		for _, el := range keyword_elements {
			if keyword, err := el.Text(); err == nil {
				articleData.Keywords = append(articleData.Keywords, keyword)
			}
		}
	}

	// Abstract
	abstractElements, _ := page.Elements(`div.abstract.author *`)
	var articleAbstract []string
	for _, node := range abstractElements {
		text, _ := node.Text()
		articleAbstract = append(articleAbstract, text)
	}
	articleData.Abstract = strings.Join(articleAbstract, " ")
	// Publish Date
	if publishDate, err := page.Element(`meta[name="citation_publication_date"]`); err == nil {
		articleData.PublishDate = *publishDate.MustAttribute("content")
	}
	articleData.Website = url
	articleData.References = get_ref(page)
	articleData.Publisher = "Science Direct"
	fmt.Print("\033[H\033[2J")
	return
}

func open_page(browser *rod.Browser, url string) (page *rod.Page, err error) {
	page = browser.MustPage() // Replace with your target URL
	err = page.Navigate(url)
	fmt.Println("Oppening the page")
	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	})
	return
}

func open_browser(browser_path, user_data_path string) (browser *rod.Browser, err error) {

	launcher := launcher.New().Headless(true).Bin(browser_path).UserDataDir(user_data_path)
	url := launcher.MustLaunch()

	// Create a new browser instance with a default client that sets headers
	browser = rod.New().ControlURL(url).MustConnect()

	return browser, err
}

func get_browser_path(files_path string) (path string) {
	filename := filepath.Join(files_path, "browser_path.txt")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
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
		fmt.Println("Error reading file:", err)
	}
	return
}

type Slices struct {
	StartSlice int `json:"start_slice"`
	EndSlice   int `json:"end_slice"`
}

func update_slices(startSlice, endSlice int, filesPath string) (err error) {
	data := Slices{StartSlice: startSlice, EndSlice: endSlice}
	slicesPaths := filepath.Join(filesPath, "article_slices.json")
	file, err := os.Create(slicesPaths)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Encode the struct as JSON and write to the file
	encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ") // Pretty-print the JSON with indentation
	if err = encoder.Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	return nil
}
func get_slices(filesPath string) (start_slice, end_slice int, err error) {
	slicesPaths := filepath.Join(filesPath, "article_slices.json")
	file, err := os.Open(slicesPaths)
	if err != nil {
		fmt.Println("Error opening slices file: ", err)
		return
	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading slices:", err)
		return
	}

	var slices Slices
	if err = json.Unmarshal(data, &slices); err != nil {
		fmt.Println("Error parsing slices JSON:", err)
		return
	}

	return slices.StartSlice, slices.EndSlice, err
}

type ArticleLINKS struct {
	ArticleLinks string `json:"article_links"`
}

func get_urls(filesPath string) (urls []string, startSlice, endSlice int, err error) {
	startSlice, endSlice, err = get_slices(filesPath)
	if err != nil {
		fmt.Println("Error parsing slices JSON:", err)
		return
	}

	urlsPaths := filepath.Join(filesPath, "article_links.jsonl")
	file, err := os.Open(urlsPaths)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Slice to hold each line from the file
	var websites []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var article ArticleLINKS
		if err := json.Unmarshal(scanner.Bytes(), &article); err != nil {
			fmt.Printf("failed to unmarshal JSON line: %v", err)
			continue
		}

		websites = append(websites, article.ArticleLinks) // Append it to the slice
	}

	urls = make([]string, len(websites[startSlice:endSlice]))
	copy(urls, websites[startSlice:endSlice])
	// Check for errors during scanning
	if err1 := scanner.Err(); err1 != nil {
		fmt.Println("Error sacnner reading file:", err1)
		return
	}

	return
}

func delete_temp_files(dir string) error {
	// Read all items in the directory
	dir = filepath.Join(dir, "user-dir")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	// Loop through each item in the directory
	for _, entry := range entries {
		// Get the full path of the item
		path := filepath.Join(dir, entry.Name())

		// Check if it is a directory
		if entry.IsDir() {
			// Remove the directory and its contents
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			// If it is a file, remove it
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Program encountered an error and recovered:", r)
			fmt.Println("Press Enter to exit...")
			fmt.Scanln()
		}
	}()
	fmt.Println("Starting the crawler ")
	current_dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	user_data_dir := filepath.Join(current_dir, "user-dir", time.Now().Format("2006-01-02_15-04-05"))
	os.Mkdir(user_data_dir, 0777)
	files_dir := filepath.Join(current_dir, "files")
	browser_path := get_browser_path(files_dir)

	urls, start_slice, end_slice, err := get_urls(files_dir)
	if err != nil {
		fmt.Println("Error getting urls:", err)
		return
	}
	crawl_counter := 0
	delete_counter := 1
	for _, url := range urls {
		new_url := "https://www.sciencedirect.com" + url
		start_time := time.Now()
		browser, err := open_browser(browser_path, user_data_dir)
		if err != nil {
			fmt.Println("Error Oppening the browser: ", err)
			return
		}

		page, err := open_page(browser, new_url)
		if err != nil {
			fmt.Println("Error Oppening the page: ", err)
			return
		}

		all_data, err := get_data(page, new_url)
		if err != nil {
			fmt.Printf("error in getting the data")
		}

		if err := save_data(files_dir, all_data); err != nil {
			fmt.Println("Error:", err)
		}

		page.MustClose()
		browser.MustClose()
		crawl_counter++
		delete_counter++
		if delete_counter%100 == 0 {
			if err := delete_temp_files(current_dir); err != nil {
				fmt.Println("Error:", err)
			}

		}
		err = update_slices(start_slice+crawl_counter, end_slice, files_dir)
		end_time := time.Since(start_time)
		fmt.Printf("\nCrawl Count Is %v", crawl_counter)
		fmt.Printf("\nCrawl Time is %v", end_time)

		if err != nil {
			fmt.Println("error in ", err)
		}
	}
}

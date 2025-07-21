package functions

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// save data to file path with json format
func Save_data(file_path string, data ArticleData) (err error) {
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


type Slices struct {
	StartSlice int `json:"start_slice"`
	EndSlice   int `json:"end_slice"`
}

// Extract urls from files/article_links.jsonl
func Update_slices(startSlice, endSlice int, filesPath string) (err error) {
	data := Slices{StartSlice: startSlice, EndSlice: endSlice}
	slicesPaths := filepath.Join(filesPath, "articles_slices.json")
	file, err := os.Create(slicesPaths)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Encode the struct as JSON and write to the file
	encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ") // Pretty-print the JSON with indentation
	if err = encoder.Encode(data); err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}
	return nil
}

func renameData(filePath string) (err error) {
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	oldFile := filepath.Join(filePath, "article_data.jsonl")
	// Get the current date and time
	newFile := filepath.Join(filePath, "article_data"+"_"+currentTime+".jsonl")
	// Rename the file
	err = os.Rename(oldFile, newFile)
	return

}

func get_slices(filesPath string) (start_slice, end_slice int, err error) {
	slicesPaths := filepath.Join(filesPath, "articles_slices.json")
	file, err := os.Open(slicesPaths)
	if err != nil {
		log.Println("Error opening slices file: ", err)
		return

	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading slices:", err)
		return
	}

	var slices Slices
	if err = json.Unmarshal(data, &slices); err != nil {
		log.Println("Error parsing slices JSON:", err)
		return
	}

	return slices.StartSlice, slices.EndSlice, err
}

// Extract urls from files/article_links.jsonl

type ArticleLINKS struct {
	ArticleLinks string `json:"article_links"`
}

func GetArticleLinks(folderFilesPath string) (urls []string, startSlice, endSlice int, err error) {

	articlePaths := filepath.Join(folderFilesPath, "article_links.jsonl")
	file, err := os.Open(articlePaths)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	// Slice to hold each line from the file
	var articles []string
	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var article ArticleLINKS
		if err := json.Unmarshal(scanner.Bytes(), &article); err != nil {
			log.Printf("failed to unmarshal JSON line: %v", err)
			continue
		}

		articles = append(articles, article.ArticleLinks) // Append it to the slice
	}

	startSlice, endSlice, err = get_slices(folderFilesPath)
	if err != nil {
		log.Println("Error parsing slices JSON:", err)
		err = Update_slices(0, len(articles), folderFilesPath)
		startSlice, endSlice = 0, len(articles)
		if err != nil {
			log.Fatal("Error in updating slices")
			return
		}
	}
	if startSlice == endSlice {
		err = renameData(folderFilesPath)
		if err != nil {
			log.Println("Article Data does not exist")
		}
		err = Update_slices(0, len(articles), folderFilesPath)
		startSlice, endSlice = 0, len(articles)
		if err != nil {
			log.Fatal("Error in updating slices")
			return
		}
	}
	if len(articles) != endSlice{

		endSlice = len(articles)
		err = Update_slices(startSlice, len(articles), folderFilesPath)
		if err != nil {
			log.Fatal("Error in updating slices")
			return
		}
	}
	urls = make([]string, len(articles[startSlice:endSlice]))
	copy(urls, articles[startSlice:endSlice])
	// Check for errors during scanning
	if err1 := scanner.Err(); err1 != nil {
		log.Println("Error sacnner reading file:", err1)
		return
	}

	return
}

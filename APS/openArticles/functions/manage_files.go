package functions

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

// save data to file path with json format




func Save_data(file_path string, data []ArticleLink) (err error) {
	fileName := filepath.Join(file_path, "article_links.jsonl")
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal the struct to JSON
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false) // Prevent escaping of &, <, >

	// Marshal and write the struct to JSON
	for _, link := range data {
		if err := encoder.Encode(link); err != nil {
			return err
		}
	}

	log.Println("Number of save issues is:",len(data))
	return err
}


type Slices struct {
	StartSlice int `json:"start_slice"`
	EndSlice   int `json:"end_slice"`
}

// Extract urls from files/article_links.jsonl
func Update_slices(startSlice, endSlice int, filesPath string) (err error) {
	data := Slices{StartSlice: startSlice, EndSlice: endSlice}
	slicesPaths := filepath.Join(filesPath, "journal_slices.json")
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
func get_slices(filesPath string) (start_slice, end_slice int, err error) {
	slicesPaths := filepath.Join(filesPath, "journal_slices.json")
	file, err := os.Open(slicesPaths)
	if err != nil {
		log.Println("Error opening slices file: ", err)
		file, err = os.Create(slicesPaths)
		if err != nil {
			log.Println(err)
			return
		}
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

type JournalLINKS struct {
	JournalLinks string `json:"journal_links"`
}

func GetJournalLinks(folderFilesPath, journalsDir string) (urls []string, startSlice, endSlice int, err error) {

	journalPaths := filepath.Join(journalsDir, "journal_links.jsonl")
	file, err := os.Open(journalPaths)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	// Slice to hold each line from the file
	var journals []string
	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var journal JournalLINKS
		if err := json.Unmarshal(scanner.Bytes(), &journal); err != nil {
			log.Printf("failed to unmarshal JSON line: %v", err)
			continue
		}

		journals = append(journals, journal.JournalLinks) // Append it to the slice
	}

	startSlice, endSlice, err = get_slices(folderFilesPath)
	if err != nil {
		log.Println("Error parsing slices JSON:", err)
		err = Update_slices(0, len(journals), folderFilesPath)
		startSlice, endSlice = 0, len(journals)
		if err != nil {
			log.Fatal("Error in updating slices")
			return
		}
	}
	if startSlice == endSlice {
		err = Update_slices(0, len(journals), folderFilesPath)
		startSlice, endSlice = 0, len(journals)
		if err != nil {
			log.Fatal("Error in updating slices")
			return
		}
	}

	// new
	if len(journals) != endSlice{
		log.Println(len(journals),endSlice)
		endSlice = len(journals)
	}
	urls = make([]string, len(journals[startSlice:endSlice]))
	copy(urls, journals[startSlice:endSlice])
	// Check for errors during scanning
	if err1 := scanner.Err(); err1 != nil {
		log.Println("Error sacnner reading file:", err1)
		return
	}

	return
}

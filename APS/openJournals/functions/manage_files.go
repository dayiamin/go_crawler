package functions

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func CheckIfFileExistsAndDelete(mainDir string) {
	// Check if the file exists
	filePath := filepath.Join(mainDir,"files","journal_links.jsonl")
	if _, err := os.Stat(filePath); err == nil {
		// File exists, delete it
		if err := os.Remove(filePath); err != nil {
		}
		log.Println("Old file deleted")
	} else if os.IsNotExist(err) {
		// File does not exist, nothing to do
		log.Println("Old file does not exist")
	} 
	 
}

// save data to file path with json format

func Save_data(file_path string, data []JournalLink) (err error) {
	log.Println("saving data")
	fileName := filepath.Join(file_path, "journal_links.jsonl")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal the struct to JSON
	for _, link := range data {
		jsonData, err := json.Marshal(link)
		if err != nil {
			return err
		}
		_, err = file.WriteString(string(jsonData) + "\n")
		if err != nil {
			return err
		}
	}

	log.Println("Journal Links saved to", fileName)
	return err
}

package functions

import (
	"log"
	"strings"
	"github.com/go-rod/rod"
)

type JournalLink struct {
	JournalLinks string `json:"journal_links"`
}

func Get_data(page *rod.Page, url string) (links []JournalLink) {
	log.Println("Oppening the Page And Getting the journal links")
	// page.MustWaitLoad()
	//a[@class="journal-thumbnail-link"]
	journalXpath := `//a[@class="journal-thumbnail-link"]`
	page.WaitElementsMoreThan(`a`, 40)
	log.Println("Using this Xpath to get the journal links \n:", journalXpath)
	elements := page.MustElementsX(`//a[@class="journal-thumbnail-link"]`)
	for _, el := range elements {
		link := *el.MustAttribute("href")

		if link != "" {
			link = strings.TrimSuffix(link, "/")
			links = append(links, JournalLink{JournalLinks: "https://journals.aps.org"+ link + `/recent` })
		}


	}
	return
}

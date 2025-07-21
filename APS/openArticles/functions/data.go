package functions

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	
)

type ArticleLink struct {
	ArticleLinks string `json:"article_links"`
}


func normalizeText(text string) string {
	text = strings.ToLower(text) // Normalize to lowercase for comparison
	// Extract the portion inside parentheses, if any
	newtext := strings.Split(text,"published ")
	if len(newtext) > 0 {
		text = newtext[1]
	}
	// Return the full text if no parentheses
	return strings.TrimSpace(text)
}

func CheckIfThreeMonthsAgo(dateStr string) (bool, error) {
	
	// Split the date string to get the months and year
	normalizedDate := normalizeText(dateStr)
	if normalizedDate == "" {
		fmt.Println("  No date found.")

	}

	// Try parsing with different formats
	
	date, err := time.Parse("2 January, 2006", normalizedDate)
	if err != nil {
		return false, err
	}
	threeMonthsAgo := time.Now().AddDate(0, -3, -31)

	// Check if the given date is before three months ago
	return date.After(threeMonthsAgo), nil

}


func Get_data(page *rod.Page, url string) (links []ArticleLink, oldData bool) {

	log.Println("Oppening the Page And Getting the issues links")
	err := page.Timeout(15*time.Second).WaitElementsMoreThan(`div.content`, 0)
	if err != nil{
		return links, true
	}
	issuesElements, err := page.ElementsX(`//div[@class="content"]`)
	if err == nil {
		for _, el := range issuesElements {
			href, err := el.ElementX(`./section[@class="headline"]/h4/a`)
			if err != nil || href == nil {
				continue
			}
			yearItems, err := el.ElementX(`./section[@class="description-container"]/p/span[2]`)
			
			if err != nil || yearItems == nil {

				continue
			}

			
			year := yearItems.MustText()
			
			update, err := CheckIfThreeMonthsAgo(year)
			if err != nil {
				log.Println("Error in checking the date in Oppening the issues :", err)
			}
			if update {

				// Append the link as a JournalLink struct
				links = append(links, ArticleLink{ArticleLinks: "https://journals.aps.org"+*href.MustAttribute("href")})
			} else {

				return links, true
			}

		}
	}

	log.Println("Number of Issues are :", len(links))
	return
}

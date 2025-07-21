package functions

import (
	"fmt"
	
	"time"

	"github.com/go-rod/rod"
	
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



// get article data of the page
func Get_data(page *rod.Page, url string) (articleData ArticleData, err error) {

	err = page.WaitElementsMoreThan(`div#abstract-section-content > p`,0)
	if err!=nil{

		return
	}
	if title := page.MustElement(`meta[property="og:title"]`); title != nil {
		articleData.Title = *title.MustAttribute("content")
	}
	// Journal
	if journal := page.MustElement(`meta[property="og:site_name"]`); journal != nil {
		articleData.Journal = *journal.MustAttribute("content")
	}
	// Authors
	author_elements, err := page.Elements(`meta[name="citation_author"]`)
	if err == nil {
		for _, el := range author_elements {
			if author := el.MustAttribute("content"); author != nil {
				articleData.Authors = append(articleData.Authors, *author)
			}
		}
	}
	// DOI
	if doi, err := page.Element(`meta[name="citation_doi"]`); err == nil {
		articleData.DOI = *doi.MustAttribute("content")
	}

	// Abstract
	if abstract, err := page.ElementX(`//div[@id="abstract-section-content"]/p`); err == nil {
		articleData.Abstract = abstract.MustText()
	}

	// Publish Date
	if publishDate, err := page.Element(`meta[name="citation_date"]`); err == nil {
		articleData.PublishDate = *publishDate.MustAttribute("content")
	}
	articleData.Website = url
	// references
	referenceElements, err := page.Timeout(5 * time.Second).ElementsX(`//ol[@class="references"]/li/span`)
	if len(referenceElements)>0 {
		var articleReferences []string
		for _, node := range referenceElements {
			text := node.MustText()
			

			articleReferences = append(articleReferences, text)
			}
		
		articleData.References = articleReferences

	
	}else {
		fmt.Println("no reference")
		articleData.References = nil
	}
	
	// keywords
	//ul[@class="physh-tagging"]/li/a
	keywordsElements, err := page.Timeout(5 * time.Second).ElementsX(`//ul[@class="physh-tagging"]/li/a`)
	if len(keywordsElements)>0 {
		var Keywords []string
		for _, node := range keywordsElements {
			text := node.MustText()
			

			Keywords = append(Keywords, text)
			}
		
		articleData.Keywords = Keywords}
	articleData.Publisher = "APS"

	return
}

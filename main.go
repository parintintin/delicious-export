package main

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type bookmark struct {
	title string
	url   string
	tags  []string
}

var bookmarks = []*bookmark{}

func main() {
	doc, _ := goquery.NewDocument("https://del.icio.us/hendrikwill")
	pages, _ := strconv.Atoi(doc.Find("ul.pagination li").Eq(4).Text())

	fmt.Println("Total pages to parse: ", pages)

	for page := 1; page <= pages; page++ {
		findBookmarks(page)
	}

	fmt.Printf("Found bookmarks: %d", len(bookmarks))
	fmt.Printf("Found bookmarks: %d", cap(bookmarks))
}

func findBookmarks(page int) {
	fmt.Printf("Parsing page %d", page)
	doc, _ := goquery.NewDocument("https://del.icio.us/hendrikwill?&page=" + strconv.Itoa(page))
	doc.Find(".articleThumbBlockOuter").Each(func(i int, s *goquery.Selection) {
		b := new(bookmark)
		b.title = s.Find(".articleTitlePan a.title").Text()
		url, _ := s.Find(".articleInfoPan a").Attr("href")
		b.url = url

		tags := s.Find(".thumbTBriefTxt ul li")
		for i := range tags.Nodes {
			b.tags = append(b.tags, tags.Eq(i).Text())
		}
		bookmarks = append(bookmarks, b)
	})
}

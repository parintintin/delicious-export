package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, _ := goquery.NewDocument("https://del.icio.us/hendrikwill")

	doc.Find(".articleThumbBlockOuter").Each(func(i int, s *goquery.Selection) {
		fmt.Println("Title: ", s.Find(".articleTitlePan a.title").Text())
		url, _ := s.Find(".articleInfoPan a").Attr("href")
		fmt.Println("URL: ", url)

		tags := s.Find(".thumbTBriefTxt ul li")

		for i := range tags.Nodes {
			fmt.Println("Tag: ", tags.Eq(i).Text())
		}

	})
}

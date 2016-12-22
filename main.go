package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	"net/http"
	"net/url"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

type bookmark struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
}

const baseURL = "https://del.icio.us"

var bookmarks = []*bookmark{}
var cookie = http.Cookie{}
var password string
var username string

func main() {
	flag.StringVar(&username, "username", "", "Your del.icio.us username")
	flag.StringVar(&password, "password", "", "Your del.icio.us password")
	flag.Parse()

	getSessionCookie()
	r := doPost(baseURL + "/" + username)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(r))
	pages, _ := strconv.Atoi(doc.Find("ul.pagination li").Eq(4).Text())

	fmt.Printf("Parsing %d pages\n", pages)
	for page := 1; page <= pages; page++ {
		findBookmarks(page)
	}

	exportBookmarks()
	fmt.Printf("Exported bookmarks: %d\n", len(bookmarks))
}

func findBookmarks(page int) {
	r := doPost(baseURL + "/" + username + "?&page=" + strconv.Itoa(page))
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(r))
	doc.Find(".articleThumbBlockOuter").Each(func(i int, s *goquery.Selection) {
		b := new(bookmark)
		b.Title = s.Find(".articleTitlePan a.title").Text()
		url, _ := s.Find(".articleInfoPan a").Attr("href")
		b.URL = url

		tags := s.Find(".thumbTBriefTxt ul li")
		for i := range tags.Nodes {
			b.Tags = append(b.Tags, tags.Eq(i).Text())
		}
		bookmarks = append(bookmarks, b)
	})
}

func getSessionCookie() {
	uv := url.Values{"username": {username}, "password": {password}, "next": {""}}
	resp, _ := http.PostForm(baseURL+"/login", uv)
	htmlResp, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile("(delavid)=(.+=)")
	kv := re.FindStringSubmatch(string(htmlResp))
	defer resp.Body.Close()
	cookie.Name = kv[1]
	cookie.Value = kv[2]
}

func doPost(url string) (doc string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(&cookie)
	var client = &http.Client{}
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(b)
}

func exportBookmarks() {
	dir := "./export"
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		} else {
			log.Println(err)
		}
	}

	jsondata, err := json.Marshal(bookmarks)
	if err != nil {
		log.Println(err)
	}

	jsonFile, err := os.Create(dir + "/data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsondata)
	jsonFile.Close()

}

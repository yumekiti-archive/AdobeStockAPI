package main

import (
	"fmt"
	"strings"
	"log"
	"regexp"
	"net/url"
	"encoding/json"
	"io/ioutil"

	"github.com/gocolly/colly"
)

type Body struct {
	name string
	number string
	path string
}

type Bodys []Body

func Scraping() {
	bodys := Bodys{}

	c := colly.NewCollector()

	path := strings.Split("http://10.201.10.133/0配布用サーバ/_AdobeStock/", "/")
	previous := strings.ToLower(url.QueryEscape(path[len(path) - 3]))

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// 以下の条件を満たす場合はスキップする
		if strings.Contains(e.Attr("href"), previous) { return }
		if strings.Contains(e.Attr("href"), "http") { return }
		if !strings.Contains(e.Attr("href"), "/") { return }

		// リンクをたどる
		e.Request.Visit(e.Attr("href"))
	})
	
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// 以下の条件を満たす場合はスキップする
		if !strings.Contains(e.Attr("href"), ".") { return }
		if !strings.Contains(e.Attr("href"), "229944361") { return }

		// URLをデコードする
		u, err := url.QueryUnescape(e.Request.URL.String())
		if err != nil {
			log.Fatal(err)
		}

		// 数字のみを抽出する
		re := regexp.MustCompile(`[0-9]+`)
		number := re.FindString(e.Attr("href"))

		// 構造体に格納する
		body := Body{
			number: number,
			name: e.Attr("href"),
			path: u,
		}
		bodys = append(bodys, body)

		fmt.Println(body)
	})

	// c.Visit("http://example/0配布用サーバ/_AdobeStock/")
	c.Visit("http://10.201.10.133/0配布用サーバ/_AdobeStock/")

	// --- ここからjsonに変換する処理 ---
	jsonData, err := json.MarshalIndent(bodys, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("info.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
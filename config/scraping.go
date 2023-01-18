package config

import (
	"strings"
	"log"
	"regexp"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gocolly/colly"
)

type Body struct {
	Name string
	Number string
	Path string
	Tag string
}

type Bodys []*Body

func Scraping(targetHost string) {
	start := time.Now()
	bodys := Bodys{}

	c := colly.NewCollector()

	path := strings.Split(targetHost, "/")
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

		// URLをデコードする
		u, err := url.QueryUnescape(e.Request.URL.String())
		if err != nil {
			log.Fatal(err)
		}

		// 数字のみを抽出する
		re := regexp.MustCompile(`[0-9]+`)
		number := re.FindString(e.Attr("href"))

		// _adobeStockの次の文字列を抽出する
		tag := strings.Split(e.Attr("href"), "/")[3]

		// 構造体に格納する
		body := &Body{
			Name: e.Attr("href"),
			Number: number,
			Path: u,
			Tag: tag,
		}
		bodys = append(bodys, body)
	})

	c.Visit(targetHost)

	// --- ここからjsonに変換する処理 ---
	jsonData, err := json.MarshalIndent(bodys, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("info.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// timeを出力する
	elapsed := time.Since(start)
	log.Printf("Scraping took %s", elapsed)
}
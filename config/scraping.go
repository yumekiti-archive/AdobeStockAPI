package config

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"

	"AdobeStockAPI/domain"
)

// 2023/01/18 16:25:08 Scraping took 4m24.452563285s
// 2023/01/18 16:54:58 Scraping took 24.863113706s
func Scraping(targetHost string) domain.Bodies {
	bodies := domain.Bodies{}

	c := colly.NewCollector()

	path := strings.Split(targetHost, "/")
	previous := strings.ToLower(url.QueryEscape(path[len(path)-3]))

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// 以下の条件を満たす場合はスキップする
		if strings.Contains(e.Attr("href"), previous) {
			return
		}
		if strings.Contains(e.Attr("href"), "http") {
			return
		}
		if !strings.Contains(e.Attr("href"), "/") {
			return
		}

		// リンクをたどる
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// 以下の条件を満たす場合はスキップする
		if !strings.Contains(e.Attr("href"), ".") {
			return
		}

		// URLをデコードする
		u, err := url.QueryUnescape(e.Request.URL.String())
		if err != nil {
			log.Fatal(err)
		}

		// 数字のみを抽出する
		re := regexp.MustCompile(`[0-9]+`)
		number := re.FindString(e.Attr("href"))

		// _adobeStockの次の文字列を抽出する
		tag := strings.Split(u, "/")[5]

		// 構造体に格納する
		body := &domain.Body{
			Name:   e.Attr("href"),
			Number: number,
			Path:   u,
			Tag:    tag,
		}
		bodies = append(bodies, body)
	})

	c.Visit(targetHost)

	if (len(bodies)) == 0 {
		log.Fatal("No data was found")
	}

	return bodies
}

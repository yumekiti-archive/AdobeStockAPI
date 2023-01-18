package main

import (
	"net/http"
	"io/ioutil"
	"log"
	
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// jsonを返す
	e.GET("/", func(c echo.Context) error {
		jsonFile, err := ioutil.ReadFile("links.json")
		if err != nil {
			log.Fatal(err)
		}

		return c.String(http.StatusOK, string(jsonFile))
	})

	// スクレイピングを再実行する
	e.GET("/scrape", func(c echo.Context) error {
		Scraping()
		return c.String(http.StatusOK, "success")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

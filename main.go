package main

import (
	"net/http"
	"io/ioutil"
	"log"

	"github.com/labstack/echo/v4"
	"AdobeStockAPI/config"
)

func main() {
	e := echo.New()

	// jsonを返す
	e.GET("/", func(c echo.Context) error {
		jsonFile, err := ioutil.ReadFile("info.json")
		if err != nil {
			log.Println(err)
		}

		return c.String(http.StatusOK, string(jsonFile))
	})

	// スクレイピングを再実行する
	e.GET("/scrape", func(c echo.Context) error {
		config.Scraping(config.GetEnv("TARGETHOST", "http://localhost"))
		return c.String(http.StatusOK, "success")
	})

	e.Logger.Fatal(e.Start(":" + config.GetEnv("PORT", "1323")))
}

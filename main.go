package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"AdobeStockAPI/config"
	"AdobeStockAPI/domain"
)

func main() {
	e := echo.New()
	db := config.NewDB()

	e.GET("/", func(c echo.Context) error {
		bodys := []domain.Body{}
		db.Find(&bodys)
		return c.JSON(http.StatusOK, bodys)
	})

	// スクレイピングを実行する
	e.GET("/scrape", func(c echo.Context) error {
		start := time.Now()
		bodys := config.Scraping(config.GetEnv("TARGETHOST", "http://localhost"))
		db.Create(&bodys)

		elapsed := time.Since(start)
		log.Printf("Scraping took %s", elapsed)
		return c.String(http.StatusOK, `{"message": "Scraping took `+elapsed.String()+`", "status": "success"}`)
	})

	e.Logger.Fatal(e.Start(":" + config.GetEnv("PORT", "1323")))
}

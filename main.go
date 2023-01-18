package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"AdobeStockAPI/config"
	"AdobeStockAPI/domain"
)

func init() {
	config.LoadEnv()
}

func main() {
	e := echo.New()
	db := config.NewDB()

	e.GET("/", func(c echo.Context) error {
		bodies := domain.Bodies{}
		db.Find(&bodies)
		return c.JSON(http.StatusOK, bodies)
	})

	// スクレイピングを実行する
	e.GET("/scrape", func(c echo.Context) error {
		start := time.Now()
		bodies := config.Scraping(config.GetEnv("TARGETHOST", ""))

		// ドロップしてから作成
		db.Migrator().DropTable(&domain.Body{})
		db.AutoMigrate(&domain.Body{})

		tx := db.Begin()
		for _, body := range bodies {
			if err := tx.Create(&body).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()

		elapsed := time.Since(start)
		log.Printf("Scraping took %s", elapsed)
		return c.String(http.StatusOK, `{"status": "ok", "message": "Scraping took `+elapsed.String()+`"}`)
	})

	e.Logger.Fatal(e.Start(":" + config.GetEnv("PORT", "1323")))
}

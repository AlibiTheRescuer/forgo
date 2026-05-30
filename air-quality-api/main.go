package main

import (
	"air-quality-api/handlers"
	"air-quality-api/models"
	"air-quality-api/parser"
	"air-quality-api/repository"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const AQICN_TOKEN = "6178b1a1d8a796f74e1df0b10887ff0e72239b65"

func main() {

	repo, err := repository.InitDB("air_quality.db") //initialize database
	if err != nil {
		log.Fatalf("Ошибка создания бд: %v", err)
	}

	seedCities(repo)

	apiClient := parser.NewParserClient(AQICN_TOKEN)
	go startBackgroundParser(repo, apiClient, 1*time.Hour)

	r := gin.Default()
	h := handlers.NewAPIHandler(repo)
	api := r.Group("/api") //all handlers
	{
		api.POST("/cities", h.CreateCity)
		api.GET("/cities", h.GetAllCities)
		api.GET("/cities/:id", h.GetCity)
		api.PUT("/cities/:id", h.UpdateCity)
		api.DELETE("/cities/:id", h.DeleteCity)

		api.POST("/air-quality", h.AddAirQualityRecord)
		api.GET("/cities/:id/air-quality", h.GetCityAirQuality)
		api.DELETE("/air-quality/:id", h.DeleteAirQuality)
	}

	log.Println("Запуск HTTP сервера на порту :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

func seedCities(repo *repository.Database) {
	existing, _ := repo.GetAllCities()
	if len(existing) > 0 { //if cities are empty add default
		return
	}

	defaultCities := []models.City{
		{Name: "Москва", Slug: "moscow", Country: "Russia"},
		{Name: "Токио", Slug: "tokyo", Country: "Japan"},
		{Name: "Сеул", Slug: "seoul", Country: "South Korea"},
		{Name: "Пағиж", Slug: "paris", Country: "France"},
		{Name: "Нью-Йорк", Slug: "newyork", Country: "USA"},
		{Name: "Пекин", Slug: "beijing", Country: "China"},
	}

	for _, city := range defaultCities {
		_ = repo.CreateCity(&city)
	}
	log.Println("Первые 6 городов по умолчанию добавлены")
}

func startBackgroundParser(repo *repository.Database, client *parser.ParserClient, interval time.Duration) {
	executeParsing(repo, client)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		executeParsing(repo, client)
	}
}

func executeParsing(repo *repository.Database, client *parser.ParserClient) {
	log.Println("Запуск парсинга показателей на фоне")
	cities, err := repo.GetAllCities()
	if err != nil {
		log.Printf("Ошибка парсера: ошибка получения списка городов: %v", err)
		return
	}

	for _, city := range cities {
		aqi, pm25, pm10, temp, src, err := client.FetchAirQuality(city.Slug)
		if err != nil {
			log.Printf("Ошибка парсера: ошибка обработки города %s (%s): %v", city.Name, city.Slug, err)
			continue
		}

		record := models.AirQuality{
			CityID:      city.ID,
			AQI:         aqi,
			PM25:        pm25,
			PM10:        pm10,
			Temperature: temp,
			Source:      src,
		}

		if err := repo.SaveAirQuality(&record); err != nil {
			log.Printf("Ошибка парсера: не удалось сохранить запись для %s: %v", city.Name, err)
		} else {
			log.Printf("Успешно сохранены данные для %s (AQI: %d)", city.Name, aqi)
		}

		time.Sleep(1 * time.Second)
	}
}

/*
GET localhost:8080/api/cities
get all cities in db

POST localhost:8080/api/cities
	{"name": "cityName", "slug": "citySlug","country": "countryName"}
add new city

GET /api/cities/1?include_records=true
get all info of city with id:1

PUT localhost:8080/api/cities/1
	{"name": "newName", "slug":"newSlug", "country": "newName"}
update info of a city id:1

DELETE localhost:8080/api/cities/1
delete all info

GET localhost:8080/api/cities/1/air-quality
get only air quality record fo a city

POST localhost:8080/api/air-quality
	{"city_id": 1, "aqi": 42, "pm25": 10.2}
update air quality manually

*/

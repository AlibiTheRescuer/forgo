package repository

import (
	"air-quality-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDB(filepath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	//initialize sqlite with automigration to update data automatically
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.City{}, &models.AirQuality{})
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

// crud operations for the cities
func (r *Database) CreateCity(city *models.City) error {
	return r.DB.Create(city).Error
}

func (r *Database) GetAllCities() ([]models.City, error) {
	var cities []models.City
	err := r.DB.Find(&cities).Error
	return cities, err
}

func (r *Database) GetCityByID(id uint, includeRecords bool) (models.City, error) {
	var city models.City
	db := r.DB
	if includeRecords {
		db = db.Preload("Records")
	}
	err := db.First(&city, id).Error
	return city, err
}

func (r *Database) UpdateCity(city *models.City) error {
	return r.DB.Save(city).Error
}

func (r *Database) DeleteCity(id uint) error {
	return r.DB.Delete(&models.City{}, id).Error
}

// crud operations for air quality
func (r *Database) SaveAirQuality(record *models.AirQuality) error {
	return r.DB.Create(record).Error
}

func (r *Database) GetAirQualityRecords(cityID uint) ([]models.AirQuality, error) {
	var records []models.AirQuality
	err := r.DB.Where("city_id = ? ", cityID).Order("created_at DESC").Find(&records).Error
	return records, err
}

func (r *Database) DeleteAirQualityRecord(id uint) error {
	return r.DB.Delete(&models.AirQuality{}, id).Error
}

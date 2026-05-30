package models

import (
	"time"

	"gorm.io/gorm"
)

type City struct { //description of city qualities
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:varchar(100);not null;unique" json:"name" binding:"required"`
	Slug      string         `gorm:"type:varchar(100);not null;unique" json:"slug" binding:"required"`
	Country   string         `gorm:"type:varchar(100)" json:"country"`
	Records   []AirQuality   `gorm:"constraint:OnDelete:CASCADE;" json:"records,omitempty"`
}

type AirQuality struct { // description of air quality properties, that are contained at exact moment
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	CityID      uint      `gorm:"not null;index" json:"city_id"`
	AQI         int       `gorm:"not null" json:"aqi"`
	PM25        float64   `json:"pm25"`
	PM10        float64   `json:"pm10"`
	Temperature float64   `json:"temperature"`
	Source      string    `json:"source"`
}

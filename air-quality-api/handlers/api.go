package handlers

import (
	"air-quality-api/models"
	"air-quality-api/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// connect network with db
type APIHandler struct {
	Repo *repository.Database
}

// initializing handler and give pointer
func NewAPIHandler(repo *repository.Database) *APIHandler {
	return &APIHandler{Repo: repo}
}

func (h *APIHandler) CreateCity(c *gin.Context) {
	//read raw text and convert to object
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //format error message and exit
		return
	}
	//insert new city
	if err := h.Repo.CreateCity(&city); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error	": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, city)
}

func (h *APIHandler) GetAllCities(c *gin.Context) {
	cities, err := h.Repo.GetAllCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cities) //return status and array of cities
}

func (h *APIHandler) GetCity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id")) //retrieve and convert id to actual number
	includeRecords := c.Query("include_records") == "true"
	city, err := h.Repo.GetCityByID(uint(id), includeRecords)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Город не найден"})
		return
	}
	c.JSON(http.StatusOK, city)
}

func (h *APIHandler) UpdateCity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	city, err := h.Repo.GetCityByID(uint(id), false)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Город не найден"})
		return
	}
	//unpack new data over existing
	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//sql update
	if err := h.Repo.UpdateCity(&city); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error	": err.Error()})
		return
	}
	c.JSON(http.StatusOK, city)
}

func (h *APIHandler) DeleteCity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Repo.DeleteCity(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error	": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Город успешно удален"})
}

func (h *APIHandler) AddAirQualityRecord(c *gin.Context) {
	var rec models.AirQuality
	if err := c.ShouldBindJSON(&rec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error	": err.Error()})
		return
	}
	if err := h.Repo.SaveAirQuality(&rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error	": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rec)
}

func (h *APIHandler) GetCityAirQuality(c *gin.Context) {
	cityID, _ := strconv.Atoi(c.Param("id"))
	records, err := h.Repo.GetAirQualityRecords(uint(cityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *APIHandler) DeleteAirQuality(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Repo.DeleteAirQualityRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Запись показателей удалена"})
}

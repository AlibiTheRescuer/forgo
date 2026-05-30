package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// separating aqi envelope to status and raw data
type AQICNResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

// mapping inner data of api response
type AQIData struct {
	AQI  int `json:"aqi"`
	IAQI struct {
		PM25 struct {
			V float64 `json:"v"`
		} `json:"pm25"`
		PM10 struct {
			V float64 `json:"v"`
		} `json:"pm10"`
		T struct {
			V float64 `json:"v"`
		} `json:"t"`
	} `json:"iaqi"`
	Attributions []struct {
		Name string `json:"name"`
	} `json:"attributions"`
}

// holds configuration for connection
type ParserClient struct {
	Token string
}

// return pointer to created client
func NewParserClient(token string) *ParserClient {
	return &ParserClient{Token: token}
}

// receive current air quality data by slug of the city
func (p *ParserClient) FetchAirQuality(citySlug string) (int, float64, float64, float64, string, error) {
	url := fmt.Sprintf("https://api.waqi.info/feed/%s/?token=%s", citySlug, p.Token)

	client := &http.Client{Timeout: 10 * time.Second} //timelimit
	resp, err := client.Get(url)
	if err != nil {
		return 0, 0, 0, 0, "", err
	}
	defer resp.Body.Close() //ensure further streamline

	//handle http errors idk
	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, 0, "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var apiResp AQICNResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, 0, 0, 0, "", err
	}

	if apiResp.Status != "ok" { //if error status then exit
		return 0, 0, 0, 0, "", fmt.Errorf("API error status: %s", apiResp.Status)
	}

	var data AQIData
	if err := json.Unmarshal(apiResp.Data, &data); err != nil {
		return 0, 0, 0, 0, "", fmt.Errorf("несуществующий слаг или станция не найдена (API вернуло: %s)", string(apiResp.Data))
	}
	//if source is unnamed then default
	source := "AQICN Data Platform"
	if len(data.Attributions) > 0 {
		source = data.Attributions[0].Name
	}

	//return typed answer
	return data.AQI, data.IAQI.PM25.V, data.IAQI.PM10.V, data.IAQI.T.V, source, nil
}

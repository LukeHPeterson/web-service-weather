package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// this is the response for the forecast endpoint from PointsResponse
type ForecastResponse struct {
	ForecastProperties struct {
		Updated         string `json:"updated"`
		Units           string `json:"units"`
		ForecastPeriods []struct {
			ID          int    `json:"number"`
			Name        string `json:"name"`
			StartTime   string `json:"startTime"`
			EndTime     string `json:"endTime"`
			Summary     string `json:"shortForecast"`
			Details     string `json:"detailedForecast"`
			Temperature int    `json:"temperature"`
		} `json:"periods"`
	} `json:"properties"`
}

type EndpointUrlResponse struct {
	PointsProperties struct {
		ID                     string `json:"@id"`
		CWA                    string `json:"cwa"`
		Office                 string `json:"forecastOffice"`
		GridX                  int64  `json:"gridX"`
		GridY                  int64  `json:"gridY"`
		EndpointForecast       string `json:"forecast"`
		EndpointForecastHourly string `json:"forecastHourly"`
	} `json:"properties"`
}

type ForecastData struct {
	Temperature string
	Summary     string
}

func main() {
	//const lat, long = 33.5139, -81.957

	router := gin.Default()
	router.GET("/forecast/:lat/:long", getForecast)
	router.Run("localhost:8080")
}

func getForecast(c *gin.Context) {
	lat := c.Param("lat")
	long := c.Param("long")

	url, err := getWeatherApiUrl(lat, long)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	data, err := getForecastData(url)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, data)
}

func getWeatherApiUrl(lat, long string) (string, error) {
	requestUrl := fmt.Sprint("https://api.weather.gov/points/", lat, ",", long)

	response, err := http.Get(requestUrl)
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var rd EndpointUrlResponse
	err = json.Unmarshal(responseData, &rd)
	if err != nil {
		return "", err
	}

	return rd.PointsProperties.EndpointForecast, nil
}

func getForecastData(url string) (ForecastData, error) {

	response, err := http.Get(url)
	if err != nil {
		return ForecastData{}, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return ForecastData{}, err
	}

	var fr ForecastResponse
	err = json.Unmarshal(responseData, &fr)
	if err != nil {
		return ForecastData{}, err
	}

	latestPeriod := fr.ForecastProperties.ForecastPeriods[0]
	temp := latestPeriod.Temperature

	if err != nil {
		return ForecastData{}, err
	}

	data := ForecastData{Temperature: getTemperatureString(temp), Summary: latestPeriod.Summary}

	return data, nil
}

func getTemperatureString(temp int) string {
	if temp > 80 {
		return "hot"
	} else if temp > 60 && temp < 80 {
		return "moderate"
	} else {
		return "cold"
	}
}

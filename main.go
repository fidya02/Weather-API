package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Condition struct {
	Text string `json:"text"`
}

type Current struct {
	TempC     float64   `json:"temp_c"`
	Humidity  int       `json:"humidity"`
	Condition Condition `json:"condition"`
	WindKph   float64   `json:"wind_kph"`
}

type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

func fetchWeatherData(location string) (WeatherResponse, error) {
	apiKey := "cc0a364fa16d48d0941135837230311"
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, location)

	response, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer response.Body.Close()

	var weatherData WeatherResponse
	err = json.NewDecoder(response.Body).Decode(&weatherData)
	if err != nil {
		return WeatherResponse{}, err
	}

	return weatherData, nil
}

func weatherHandler(c echo.Context) error {
	
	location := c.QueryParam("location")

	
	if location == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Parameter Lokasi dibutuhkan"})
	}

	weatherData, err := fetchWeatherData(location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data cuaca"})
	}

	return c.JSON(http.StatusOK, weatherData)
}

func main() {
	e := echo.New()

	e.Static("/", "views") 

	e.GET("/weather", weatherHandler)

	port := 8080
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

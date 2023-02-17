package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CitiesJson struct {
	Areaname  string      `json:"areaname"`
	Origin    []float64   `json:"origin"`
	ModelInfo []ModelInfo `json:"modelinfo"`
}

type ModelInfo struct {
	Filename string    `json:"filename"`
	Lowerpos []float64 `json:"lowerpos"`
	Upperpos []float64 `json:"upperpos"`
}

type ResponseModelInfo struct {
	Url      string    `json:"url"`
	Lowerpos []float64 `json:"lowerpos"`
	Upperpos []float64 `json:"upperpos"`
}

type SearchCityModelParams struct {
	Longitude float64 `query:"longitude"`
	Latitude  float64 `query:"latitude"`
	Alt       float64 `query:"alt"`
	Radius    float64 `query:"radius"`
}

type ResponseSearchCityModel struct {
	Status       string              `json:"status"`
	ErrorMessage string              `json:"errorMessage"`
	Items        []ResponseModelInfo `json:"items"`
}

var sapporoCities CitiesJson

const baseURL = "https://snow-globe.almikan.com"

func main() {
	initSapporoCities()

	e := echo.New()
	e.Static("/public", "public")

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	g := e.Group("/api")
	g.GET("/search-city-model", func(c echo.Context) error {
		var searchParams SearchCityModelParams
		err := c.Bind(&searchParams)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseSearchCityModel{
				Status:       "error",
				ErrorMessage: "bad request",
				Items:        nil,
			})
		}

		resultCities := searchCityModel(searchParams)

		return c.JSON(http.StatusOK, ResponseSearchCityModel{
			Status:       "ok",
			ErrorMessage: "",
			Items:        resultCities,
		})
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func initSapporoCities() {
	// 札幌の街データの読み込み
	raw, err := ioutil.ReadFile("./sapporo.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &sapporoCities)
	if err != nil {
		panic(err)
	}
}

func searchCityModel(searchParams SearchCityModelParams) []ResponseModelInfo {
	urlList := make([]ResponseModelInfo, 0)
	for _, modelInfo := range sapporoCities.ModelInfo {
		llposDist := calcDist(
			modelInfo.Lowerpos[0], modelInfo.Lowerpos[1], searchParams.Alt,
			searchParams.Latitude, searchParams.Longitude, searchParams.Alt,
		)
		ulposDist := calcDist(
			modelInfo.Upperpos[0], modelInfo.Lowerpos[1], searchParams.Alt,
			searchParams.Latitude, searchParams.Longitude, searchParams.Alt,
		)
		uuposDist := calcDist(
			modelInfo.Upperpos[0], modelInfo.Upperpos[1], searchParams.Alt,
			searchParams.Latitude, searchParams.Longitude, searchParams.Alt,
		)
		luposDist := calcDist(
			modelInfo.Lowerpos[0], modelInfo.Upperpos[1], searchParams.Alt,
			searchParams.Latitude, searchParams.Longitude, searchParams.Alt,
		)

		if llposDist < searchParams.Radius &&
			ulposDist < searchParams.Radius &&
			uuposDist < searchParams.Radius &&
			luposDist < searchParams.Radius {
			urlList = append(urlList, ResponseModelInfo{
				Url:      baseURL + "/public/model/sapporo_256/" + modelInfo.Filename,
				Lowerpos: modelInfo.Lowerpos,
				Upperpos: modelInfo.Upperpos,
			})
		}
	}

	return urlList
}

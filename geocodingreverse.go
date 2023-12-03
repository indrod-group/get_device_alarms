package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func UnmarshalGeoapifyResponse(data []byte) (GeoapifyResponse, error) {
	var r GeoapifyResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GeoapifyResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GeoapifyResponse struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
	Query    Query     `json:"query"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
	Bbox       []float64  `json:"bbox"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	Datasource   Datasource `json:"datasource"`
	Name         string     `json:"name"`
	Country      string     `json:"country"`
	CountryCode  string     `json:"country_code"`
	State        string     `json:"state"`
	City         string     `json:"city"`
	Postcode     string     `json:"postcode"`
	District     string     `json:"district"`
	Suburb       string     `json:"suburb"`
	Street       string     `json:"street"`
	Housenumber  string     `json:"housenumber"`
	Lon          float64    `json:"lon"`
	Lat          float64    `json:"lat"`
	Distance     float64    `json:"distance"`
	ResultType   string     `json:"result_type"`
	Formatted    string     `json:"formatted"`
	AddressLine1 string     `json:"address_line1"`
	AddressLine2 string     `json:"address_line2"`
	Category     string     `json:"category"`
	Timezone     Timezone   `json:"timezone"`
	PlusCode     string     `json:"plus_code"`
	Rank         Rank       `json:"rank"`
	PlaceID      string     `json:"place_id"`
}

type Datasource struct {
	Sourcename  string `json:"sourcename"`
	Attribution string `json:"attribution"`
	License     string `json:"license"`
	URL         string `json:"url"`
}

type Rank struct {
	Importance float64 `json:"importance"`
	Popularity float64 `json:"popularity"`
}

type Timezone struct {
	Name             string `json:"name"`
	OffsetSTD        string `json:"offset_STD"`
	OffsetSTDSeconds int64  `json:"offset_STD_seconds"`
	OffsetDST        string `json:"offset_DST"`
	OffsetDSTSeconds int64  `json:"offset_DST_seconds"`
	AbbreviationSTD  string `json:"abbreviation_STD"`
	AbbreviationDST  string `json:"abbreviation_DST"`
}

type Query struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	PlusCode string  `json:"plus_code"`
}

const GEOAPIFY_URL_REVERSE_GEOCODING = "https://api.geoapify.com/v1/geocode/reverse?lat=%s&lon=%s&apiKey=%s"

func GetAddress(lat string, lng string) *string {
	apiKey := os.Getenv("GEOAPIFY_KEY")
	url := fmt.Sprintf(GEOAPIFY_URL_REVERSE_GEOCODING, lat, lng, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making the request to Geoapify: %s\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body: %s\n", err)
		return nil
	}

	var data GeoapifyResponse
	data, err = UnmarshalGeoapifyResponse(body)
	if err != nil {
		fmt.Printf("Error unmarshalling the JSON response: %s\n", err)
		return nil
	}

	if len(data.Features) > 0 {
		return &data.Features[0].Properties.Formatted
	}

	return nil
}

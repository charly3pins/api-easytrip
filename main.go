package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
)

const (
	apiKey              = "ca722829566893739318779257539996"
	apiKeyTest          = "prtl6749387986743898559646983194"
	flightsURL          = "http://partners.api.skyscanner.net/apiservices/pricing/v1.0/"
	hotelsURLByEntity   = "https://gateway.skyscanner.net/hotels/v1/prices/search/entity/"
	hotelsURLByLocation = "https://gateway.skyscanner.net/hotels/v1/prices/search/location/"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET, PUT, POST, DELETE, OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/flights", flights)
	router.GET("/hotels", hotels)

	router.Run(":8080")
}

func flights(c *gin.Context) {
	form := url.Values{}
	form.Add("country", "UK")
	form.Add("currency", "GBP")
	form.Add("locale", "en-GB")
	form.Add("locationSchema", "iata")
	form.Add("apikey", apiKey)
	form.Add("grouppricing", "on")
	form.Add("originplace", "EDI")
	form.Add("destinationplace", "LHR")
	form.Add("outbounddate", "2018-02-24")
	form.Add("inbounddate", "2018-03-03")
	form.Add("adults", "1")
	form.Add("children", "0")
	form.Add("infants", "0")
	form.Add("cabinclass", "Economy")

	req, err := http.NewRequest("POST", flightsURL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[create session] Error making request for flights in Skyscanner", err)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		log.Println("[create session] Error in response. Code", resp.StatusCode)
		return
	}

	requestPoll := fmt.Sprintf("%s?apiKey=%s", resp.Header.Get("Location"), apiKey)
	req, err = http.NewRequest("GET", requestPoll, nil)
	req.Header.Add("Accept", "application/json")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println("[poll results] Error making request for flights in Skyscanner", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("[poll results] Error in response. Code", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("[poll results] Error reading the body: ", err)
		return
	}
	log.Println(string(body))
}

func hotels(c *gin.Context) {
	// location := "55.95,-3.37-latlong"
	// market := "ES"
	// locale := "es-ES"
	// checkin := "2018-03-15"
	// checkout := "2018-03-20"
	// currency := "EUR"
	// rooms := "1"
	// adults := "2"
	// finalHotelsURLByLocation := fmt.Sprintf("%s%s?market=%s&locale=%s&checkin_date=%s&checkout_date=%s&curreny=%s&rooms=%s&adults=%s", hotelsURLByLocation, location, market, locale, checkin, checkout, currency, rooms, adults)

	entity := "29475375"
	market := "ES"
	locale := "es-ES"
	checkin := "2018-03-15"
	checkout := "2018-03-20"
	currency := "EUR"
	rooms := "1"
	adults := "2"
	finalHotelsURLByEntity := fmt.Sprintf("%s%s?market=%s&locale=%s&checkin_date=%s&checkout_date=%s&curreny=%s&rooms=%s&adults=%s&apikey=%s", hotelsURLByEntity, entity, market, locale, checkin, checkout, currency, rooms, adults, apiKeyTest)

	req, err := http.NewRequest("GET", finalHotelsURLByEntity, nil)
	req.Header.Add("x-user-agent", "D;B2B")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[create session] Error making request for flights in Skyscanner", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("[poll results] Error in response. Code", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("[poll results] Error reading the body: ", err)
		return
	}
	log.Println(string(body))
}

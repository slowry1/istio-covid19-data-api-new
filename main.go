package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var currentCovid19Data string
var counter int
var covidUpdateTimeInSeconds time.Duration

func getCovid19Data() {
	response, err := http.Get("https://api.covid19api.com/all") //"https://api.covid19api.com/all") //"https://api.covid19api.com/live/country/us/status/confirmed") //"https://api.covid19api.com/all") //http.Get("https://api.covid19api.com/") //"https://api.covid19api.com/dayone/country/us/status/confirmed")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Close the response body when finished with it
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			currentCovid19Data = string(responseData)
			counter += 1
			fmt.Println("Number of times the currentCovid19Data has been updated: ", counter)
		}
	}
}

func main() {
	// // Testing json.marshal
	// type FruitBasket struct {
	// 	Name    string
	// 	Fruit   []string
	// 	Id      int64  `json:"ref"`
	// 	private string // An unexported field is not encoded.
	// 	Created time.Time
	// }

	// // Testing json.marshal
	// basket := FruitBasket{
	// 	Name:    "Standard",
	// 	Fruit:   []string{"Apple", "Banana", "Orange"},
	// 	Id:      999,
	// 	private: "Second-rate",
	// 	Created: time.Now(),
	// }

	// // Testing json.marshal
	// var jsonData []byte
	// jsonData, err := json.Marshal(basket)
	// //jsonData, err := json.MarshalIndent(basket, "", "   ")
	// if err != nil {
	// 	log.Println(err)
	// }

	// Set the update timer time so the covid19 data gets updated every X seconds
	covidUpdateTimeInSeconds = 3600 // 3600 seconds in an hour
	// Populate currentCovid19Data for the first time to make sure it not empty if someone requests the data
	getCovid19Data()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-Type", "application/json")
		fmt.Fprintf(w, currentCovid19Data) //string(jsonData))
	})

	// This is a goroutine for an anonoymous function call that will run asynchronously.
	// This will start the timer for getting the most current Covid19 data and saving it.
	// covidUpdateTimeInSeconds is passed to it
	go func(covidUpdateTime time.Duration) {
		for true {
			fmt.Println("Getting covid data     Hello !!")
			getCovid19Data()
			time.Sleep(covidUpdateTime * time.Second)
		}
	}(covidUpdateTimeInSeconds)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

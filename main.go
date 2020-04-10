package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var listOfCountries = []string{"united-states", "italy", "spain", "germany", "china"} // by their Slug
var listOfStatus = []string{"confirmed", "recovered", "deaths"}

var covidDataMap = make(map[string]string)

// UpdateTime 3600 seconds in an hour
const UpdateTime = 3600

var covidUpdateTime time.Duration
var currentCovid19SummaryData string
var currentCovid19CountriesData string

func getCovid19SummaryData() {
	response, err := http.Get("https://api.covid19api.com/summary")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Close the response body when finished with it
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			currentCovid19SummaryData = string(responseData)
			msg := "Current Covid19 Summary data updated!"
			fmt.Println(msg)
			log.Print(msg)
		}
	}
}

func getCovid19CountriesData() {
	response, err := http.Get("https://api.covid19api.com/countries")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Close the response body when finished with it
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			currentCovid19CountriesData = string(responseData)
			msg := "Current Covid19 Countries data updated!"
			fmt.Println(msg)
			log.Print(msg)
		}
	}
}

// We are cheating right now hitting API2 endpoing to get the heroku timeseries for the US since we couldn't get the API1 US data since it is to large.
func getCovid19USTimeseriesByStatus() {
	for _, status := range listOfStatus {
		// Skip recovered becuase the api doesnt have this data
		if status == "recovered" {
			continue
		}
		response, err := http.Get("http://covid2019-api.herokuapp.com/v2/timeseries/US/" + status)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			log.Printf("The HTTP request failed with error %s\n", err)
		} else {
			// Close the response body when finished with it
			defer response.Body.Close()
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			} else {
				mapKeyString := "us_timeseries_" + status
				covidDataMap[mapKeyString] = string(responseData)
				msg := "Current Covid19 timeseries " + status + " data updated!"
				fmt.Println(msg)
				log.Print(msg)
			}
		}
	}
}

func getCovid19LiveDataByCounrtyAndStatus() {
	for _, country := range listOfCountries {
		// Skip united states for this one since it might take a while. test the other ones first before breaking the program
		if country == "united-states" {
			continue
		}
		for _, status := range listOfStatus {

			response, err := http.Get("https://api.covid19api.com/country/" + country + "/status/" + status + "/live")
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				log.Printf("The HTTP request failed with error %s\n", err)
			} else {
				// Close the response body when finished with it
				defer response.Body.Close()
				responseData, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				} else {
					mapKeyString := country + "_live_data_" + status
					covidDataMap[mapKeyString] = string(responseData)
					msg := "Current Covid19 " + country + " live " + status + " data updated!"
					fmt.Println(msg)
					log.Print(msg)
				}
			}

		}
	}
}

func covid19Summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, currentCovid19SummaryData)
}

func covid19Countries(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, currentCovid19CountriesData)
}

func covid19USTimesriesConfirmed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["us_timeseries_confirmed"])
}

func covid19USTimesriesDeaths(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["us_timeseries_deaths"])
}

// Group of confirmed cases
func covid19USLiveConfirmed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["united-states_live_data_confirmed"])
}

func covid19ItalyLiveConfirmed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["italy_live_data_confirmed"])
}

func covid19SpainLiveConfirmed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["spain_live_data_confirmed"])
}

func covid19GermanyLiveConfirmed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["germany_live_data_confirmed"])
}

func covid19ChinaLiveConfirmed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["china_live_data_confirmed"])
}

// Group of recovered cases
func covid19USLiveRecovered(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["united-states_live_data_recovered"])
}

func covid19ItalyLiveRecovered(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["italy_live_data_recovered"])
}

func covid19SpainLiveRecovered(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["spain_live_data_recovered"])
}

func covid19GermanyLiveRecovered(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["germany_live_data_recovered"])
}

func covid19ChinaLiveRecovered(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["china_live_data_recovered"])
}

// Group of deaths cases
func covid19USLiveDeaths(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["united-states_live_data_deaths"])
}

func covid19ItalyLiveDeaths(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["italy_live_data_deaths"])
}

func covid19SpainLiveDeaths(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["spain_live_data_deaths"])
}

func covid19GermanyLiveDeaths(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["germany_live_data_deaths"])
}

func covid19ChinaLiveDeaths(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-Type", "application/json")
	fmt.Fprintf(w, covidDataMap["china_live_data_deaths"])
}

func handleRequests() {
	http.HandleFunc("/summary", covid19Summary)
	http.HandleFunc("/countries", covid19Countries)
	http.HandleFunc("/country/timeseries/united-states/confirmed", covid19USTimesriesConfirmed)
	http.HandleFunc("/country/timeseries/united-states/deaths", covid19USTimesriesDeaths)
	http.HandleFunc("/country/united-states/status/confirmed/live", covid19USLiveConfirmed)
	http.HandleFunc("/country/italy/status/confirmed/live", covid19ItalyLiveConfirmed)
	http.HandleFunc("/country/spain/status/confirmed/live", covid19SpainLiveConfirmed)
	http.HandleFunc("/country/germany/status/confirmed/live", covid19GermanyLiveConfirmed)
	http.HandleFunc("/country/china/status/confirmed/live", covid19ChinaLiveConfirmed)
	http.HandleFunc("/country/united-states/status/recovered/live", covid19USLiveRecovered)
	http.HandleFunc("/country/italy/status/recovered/live", covid19ItalyLiveRecovered)
	http.HandleFunc("/country/spain/status/recovered/live", covid19SpainLiveRecovered)
	http.HandleFunc("/country/germany/status/recovered/live", covid19GermanyLiveRecovered)
	http.HandleFunc("/country/china/status/recovered/live", covid19ChinaLiveRecovered)
	http.HandleFunc("/country/united-states/status/deaths/live", covid19USLiveDeaths)
	http.HandleFunc("/country/italy/status/deaths/live", covid19ItalyLiveDeaths)
	http.HandleFunc("/country/spain/status/deaths/live", covid19SpainLiveDeaths)
	http.HandleFunc("/country/germany/status/deaths/live", covid19GermanyLiveDeaths)
	http.HandleFunc("/country/china/status/deaths/live", covid19ChinaLiveDeaths)
}

func main() {
	// setup logging file
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	// Set the update timer time so the covid19 data gets updated every X seconds
	covidUpdateTime = UpdateTime
	log.Print("Data is set to update every " + string(UpdateTime) + " seconds.")
	// Populate currentCovid19Data for the first time to make sure it not empty if someone requests the data
	getCovid19SummaryData()
	// Populate currentCovid19CountriesData for the first time to make sure it not empty if someone requests the data
	getCovid19CountriesData()

	getCovid19LiveDataByCounrtyAndStatus()

	getCovid19USTimeseriesByStatus()

	handleRequests()

	// This is a goroutine for an anonoymous function call that will run asynchronously.
	// This will start the timer for getting the most current Covid19 data and saving it.
	// covidUpdateTime is passed to it
	go func(covidUpdateTime time.Duration) {
		for true {
			msg := "Updating stored covid data"
			fmt.Println(msg)
			log.Print(msg)
			getCovid19SummaryData()
			getCovid19CountriesData()
			getCovid19LiveDataByCounrtyAndStatus()
			getCovid19USTimeseriesByStatus()
			time.Sleep(covidUpdateTime * time.Second)
		}
	}(covidUpdateTime)

	log.Fatal(http.ListenAndServe(":9001", nil))
}

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func fetchGeocode(l Location) Coordinate {
	baseUrl := "https://maps.googleapis.com/maps/api/geocode/json?address="
	address := l.Address + " " + l.City + " " + l.State
	authUrl := "&key=AIzaSyCv2TRAAF2y9GokOKq3UTovq5KYT2R2qDA"
	url := strings.Replace(baseUrl+address+authUrl, " ", "%20", -1)
	//	fmt.Println("URL:", url)
	println("Hitting url:", baseUrl, address, authUrl)
	res, err := http.Get(url)
	defer res.Body.Close()
	checkErr(err)
	//	fmt.Println(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)
	resp := make(map[string]interface{})
	err = json.Unmarshal(body, &resp)
	checkErr(err)
	//	fmt.Println(resp)
	results := resp["results"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"]
	lat := results.(map[string]interface{})["lat"]
	lng := results.(map[string]interface{})["lng"]

	//	fmt.Println(lat, lng)
	return Coordinate{lat.(float64), lng.(float64)}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchRoute(from_location, to_location Location, c chan UberEstimate) {
	from := from_location.Coordinate
	to := to_location.Coordinate
	baseUrl := "https://api.uber.com/v1/estimates/price?"
	params := "start_longitude=" + strconv.FormatFloat(from.Lng, 'f', -1, 64) +
		"&end_longitude=" + strconv.FormatFloat(to.Lng, 'f', -1, 64) +
		"&start_latitude=" + strconv.FormatFloat(from.Lat, 'f', -1, 64) +
		"&end_latitude=" + strconv.FormatFloat(to.Lat, 'f', -1, 64)
	url := strings.Replace(baseUrl+params, " ", "%20", -1)
	println("Hitting url:", baseUrl, params)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "token CFfWssVgYg6uHUPMO3M9VHH9W3jRF--7mTzy8zxN")
	res, _ := client.Do(req)
	defer res.Body.Close()
	checkErr(err)
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	type Prices struct {
		EstimateList []UberEstimate `json:"prices"`
	}
	estimate_list := Prices{}
	err = json.Unmarshal(body, &estimate_list)
	checkErr(err)

	for i := 0; i < len(estimate_list.EstimateList); i++ {
		if strings.Compare(estimate_list.EstimateList[i].Display_name, "uberX") == 0 {
			estimate_list.EstimateList[i].Destination = to_location
			c <- estimate_list.EstimateList[i]
			return
		}
	}
}

func find_shortest_path(start_location Location, destinations []Location) []UberEstimate {
	best_path := find_shortest_path_recurse(start_location, destinations)
	
	// Add last leg of trip to make it round trip
	messages := make(chan UberEstimate)
	go fetchRoute(best_path[len(best_path)-1].Destination, start_location, messages)
	estimate := <-messages
	best_path = append(best_path, estimate)

	//	fmt.Println(best_path)
	return best_path
}

// Use Goroutine to find shortest path recursively
func find_shortest_path_recurse(start_location Location, destinations []Location) (best_path []UberEstimate) {

	if len(destinations) == 0 {
		return best_path
	}
	messages := make(chan UberEstimate, len(destinations))
	for i := 0; i < len(destinations); i++ {
		go fetchRoute(start_location, destinations[i], messages)
	}

	minimum_price := 1000
	var best_estimation UberEstimate

	for i := 0; i < len(destinations); i++ {
		estimate := <-messages
		//				fmt.Println(estimate)
		if estimate.Low_estimate < minimum_price {
			best_estimation = estimate
			minimum_price = best_estimation.Low_estimate
		}
	}
	//	fmt.Println(best_estimation)

	var new_destinations []Location
	// Remove destination from destinations list and call this function again
	for i, v := range destinations {
		if v.Id == best_estimation.Destination.Id {
			new_destinations = append(destinations[:i], destinations[i+1:]...)
		}
	}
	//	fmt.Println("NEW DESTINATIONS::::")
	//	fmt.Println(new_destinations)
	best_path = append(best_path, best_estimation)
	return append(best_path, find_shortest_path_recurse(best_estimation.Destination, new_destinations)...)
}

const REQUEST_BASE_URL = "https://sandbox-api.uber.com/v1/requests?"
const OAUTH_ACCESS_TOKEN = "access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiJhN2FhNmVhNS0zNzE1LTRiYWUtYTI0Ny04ZGEzZjhlZTkxYzkiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6ImFhYmU2YmNmLTg2OWQtNDc2Ni04OWNiLTg1M2QwOTdmYWU5ZCIsImV4cCI6MTQ1MDczODMyNCwiaWF0IjoxNDQ4MTQ2MzIzLCJ1YWN0IjoiYmZHRjFSY3BITW91em56elNvMGRRSkFMQ0Z6OFRVIiwibmJmIjoxNDQ4MTQ2MjMzLCJhdWQiOiJfRUZaOTdEdkZUejQyOVVoOTFiSlpwTjFQN2k0cV81ViJ9.LqPlWmuCLV2YS_63qccDrG8m6Ugp2DPJ6iuKJRSAa1pefIZu5w6Izi7-JoHYn28ZyjQ-utWfr1VzGw4C6GEKlFssQ8pwBzMOVq9D90NJ10h7EIMihuojOJWZop24qkQi5EGXFYavejD3Vxa2jsO-0ASs1SQWITIIDhGw3457CSPle2BgQnF1FFfpc3Y5GPib-qxrgKPU4Z1t8zi9BAySmDi8_DxwHflAV84QffI6loLH0f062cR_UmKfhpA6ENYuoQXcG9t0f2f5OBgFIeRzzuwPEKJZSpBWrF7Zo5L6ile5aa_nuXw6QAebEhxlbYKZhFV2EGi289vBqWULlUjspw"

func request_uber(from_location, to_location *Location, product_id string) (UberRequest, error) {
	response := UberRequest{}
	from := from_location.Coordinate
	to := to_location.Coordinate

	params := UberRequestParams{}

	params.Start_latitude = from.Lat
	params.Start_longitude = from.Lng
	params.End_latitude = to.Lat
	params.End_longitude = to.Lng
	params.Product_id = product_id

	url := strings.Replace(REQUEST_BASE_URL+OAUTH_ACCESS_TOKEN, " ", "%20", -1)
	println("Hitting url:", url)
	client := &http.Client{}
	json_params, err := json.Marshal(params)
	if err != nil {
		return response, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_params))
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	defer res.Body.Close()
	checkErr(err)
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	err = json.Unmarshal(body, &response)
	checkErr(err)
	//	fmt.Println("UBER::::::::",response)
	return response, nil
}

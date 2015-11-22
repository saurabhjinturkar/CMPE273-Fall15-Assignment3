package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func createLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	var l Location
	err = json.Unmarshal(body, &l)
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(err)
		return
	}
	//	fmt.Println(l)

	if !l.validate() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	locationDAO := LocationController{}
	loc, err := locationDAO.CreateLocation(l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(loc)
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	location_id := vars["location_id"]
	i, err := strconv.Atoi(location_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	locationDAO := LocationController{}
	loc, err := locationDAO.GetLocation(i)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(loc)
}

func deleteLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	location_id := vars["location_id"]
	i, err := strconv.Atoi(location_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	locationDAO := LocationController{}
	err = locationDAO.DeleteLocation(i)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(err)
		return
	}
}

func updateLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	location_id := vars["location_id"]
	i, err := strconv.Atoi(location_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var l Location
	err = json.Unmarshal(body, &l)
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(err)
	}
	//	fmt.Println(l)
	locationDAO := LocationController{}
	loc, err := locationDAO.UpdateLocation(i, l)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loc)
}

func planTrip(w http.ResponseWriter, r *http.Request) {
	log.Println("/trips")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	var request TripsRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	//	fmt.Println(request)
	uber := UberService{}
	plan, err := uber.findBestRoute(request)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())
	}
	t := TripsResponse{}
	t.prepareResponse(plan)
	json.NewEncoder(w).Encode(t)
}

func getPlan(w http.ResponseWriter, r *http.Request) {
	log.Println("/trips/{trip_id}")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	plan_id := vars["trip_id"]
	id, err := strconv.Atoi(plan_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	uber := UberService{}
	plan, err := uber.GetPlan(id)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(err)
		return
	}
	t := TripsResponse{}
	t.prepareResponse(plan)
	json.NewEncoder(w).Encode(t)
}

func startTrip(w http.ResponseWriter, r *http.Request) {
	log.Println("/trips/{trip_id}/request")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	plan_id := vars["trip_id"]
	id, err := strconv.Atoi(plan_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	uber := UberService{}
	reply, err := uber.StartTrip(id)
	if err != nil {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(reply)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/locations", createLocation).Methods("POST")
	r.HandleFunc("/locations/{location_id}", getLocation).Methods("GET")
	r.HandleFunc("/locations/{location_id}", deleteLocation).Methods("DELETE")
	r.HandleFunc("/locations/{location_id}", updateLocation).Methods("PUT")

	r.HandleFunc("/trips", planTrip).Methods("POST")
	r.HandleFunc("/trips/{trip_id}", getPlan).Methods("GET")
	r.HandleFunc("/trips/{trip_id}/request", startTrip).Methods("PUT")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", r))
}

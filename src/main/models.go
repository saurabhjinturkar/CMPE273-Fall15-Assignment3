package main

type Coordinate struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json: "lng" bson:"lng"`
}

type Location struct {
	Id         int        `json:"id" bson:"_id,omitempty"`
	Name       string     `json:"name" bson:"name"`
	Address    string     `json:"address" bson:"address"`
	City       string     `json:"city" bson:"city"`
	State      string     `json:"state" bson:"state"`
	Zip        string     `json:"zip" bson:"zip"`
	Coordinate Coordinate `json:"coordinates" bson:"coordinates"`
}

func (l *Location) validate() bool {
	if len(l.Address) == 0 || len(l.City) == 0 || len(l.State) == 0 {
		return false
	}
	return true
}

/*
Structure for ID. It can be location ID or plan ID.
*/
type ID struct {
	Type string `json:"type" bson:"type"`
	Id   int    `json:"id" bson:"id"`
}

type TripsRequest struct {
	Starting_from_location_id int   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	Location_ids              []int `json:"location_ids" bson:"location_ids"`
}

type UberEstimate struct {
	Product_id       string  `json:"product_id"`       // : "a27a867a-35f4-4253-8d04-61ae80a40df5",
	Currency_code    string  `json:"currency_code"`    //: "USD",
	Display_name     string  `json:"display_name"`     //: "uberX",
	Estimate         string  `json:"estimate"`         //: "$15",
	Low_estimate     int     `json:"low_estimate`      //": 15,
	High_estimate    int     `json:"high_estimate"`    //: 15,
	Surge_multiplier float64 `json:"surge_multiplier"` //: 1,
	Duration         int     `json:"duration"`         //: 640,
	Distance         float64 `json:"distance"`         //: 5.34
	Destination      Location
}

type Plan struct {
	Id                           int     `json:"id" bson:"_id,omitempty"`
	Status                       string  `json:"status" bson:"status"`
	Starting_from_location_id    int     `json:"starting_from_location_id" bson: "starting_from_location_id"`
	Next_destination_location_id int     `json:"next_destination_location_id,omitempty" bson:"next_destination_location_id,omitempty"`
	Best_route_location_ids      []int   `json:"best_route_location_ids" bson:"best_route_location_ids"`
	Total_uber_costs             int     `json:"total_uber_costs" bson:"total_uber_costs"`
	Total_uber_duration          int     `json:"total_uber_duration" bson: "total_uber_duration"`
	Total_distance               float64 `json:"total_distance" bson:"total_distance"`
	Uber_wait_time_eta           int     `json:"uber_wait_time_eta,omitempty" bson:"uber_wait_time_eta,omitempty"`
	Next_destination_index       int     `json:"-" bson:"next_destination_index,omitempty"`
	Product_id                   string  `json:"-" bson:"product_id,omitempty"`
}

type UberRequest struct {
	Status           string  `json:"status"`
	Request_id       string  `json:"request_id"`
	Driver           string  `json:"driver"`
	ETA              int     `json:"eta"`
	Location         string  `json:"location"`
	Vehicle          string  `json:"vehicle"`
	Surge_multiplier float64 `json:"surge_multiplier"`
}

type UberRequestParams struct {
	Start_longitude float64 `json:"start_longitude"`
	Start_latitude  float64 `json:"start_latitude"`
	End_longitude   float64 `json:"end_longitude"`
	End_latitude    float64 `json:"end_latitude"`
	Product_id      string  `json:"product_id"`
}

type TripsResponse struct {
	Id                           int     `json:"id" bson:"_id,omitempty"`
	Status                       string  `json:"status" bson:"status"`
	Starting_from_location_id    int     `json:"starting_from_location_id" bson: "starting_from_location_id"`
	Best_route_location_ids      []int   `json:"best_route_location_ids" bson:"best_route_location_ids"`
	Total_uber_costs             int     `json:"total_uber_costs" bson:"total_uber_costs"`
	Total_uber_duration          int     `json:"total_uber_duration" bson: "total_uber_duration"`
	Total_distance               float64 `json:"total_distance" bson:"total_distance"`
}

func (t *TripsResponse) prepareResponse(plan Plan) {
	t.Id = plan.Id
	t.Status = plan.Status
	t.Starting_from_location_id = plan.Starting_from_location_id
	t.Best_route_location_ids = plan.Best_route_location_ids
	t.Total_uber_costs = plan.Total_uber_costs
	t.Total_uber_duration = plan.Total_uber_duration
	t.Total_distance = plan.Total_distance
}
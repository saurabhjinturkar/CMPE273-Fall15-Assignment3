package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LocationController struct{}

var url string = "test:test@ds059888.mongolab.com:59888/cmpe273" // "mongodb://localhost"

func generateID(session mgo.Session, id_type string) (int, error) {
	c := session.DB("cmpe273").C("id")
	result := ID{}
	err := c.Find(bson.M{"type": id_type}).One(&result)
	if err != nil {
		return -1, err
	}
	//	fmt.Println(result)
	id := result.Id + 1
	colQuerier := bson.M{"type": id_type}
	change := bson.M{"$set": bson.M{"id": id}}
	err = c.Update(colQuerier, change)
	return id, nil
}

func (l *LocationController) CreateLocation(location Location) (loc Location, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return loc, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("location")
	location.Id, err = generateID(*session, "location")
	if err != nil {
		return loc, err
	}
	g := fetchGeocode(location)
	location.Coordinate = g
	err = c.Insert(location)
	if err != nil {
		return loc, err
	}
	return location, nil
}

func (l *LocationController) GetLocation(id int) (loc Location, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return loc, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("location")
	result := Location{}
	err = c.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		return loc, err
	}
	//	fmt.Println(result)
	return result, nil
}

func (l *LocationController) GetLocationByIds(ids []int) (locs []Location, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return locs, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("location")
	result := []Location{}
	err = c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&result)
	if err != nil {
		return locs, err
	}
	return result, nil
}

func (l *LocationController) DeleteLocation(id int) error {
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("location")
	err = c.Remove(bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (l *LocationController) UpdateLocation(id int, location Location) (loc Location, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return loc, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("location")
	colQuerier := bson.M{"_id": id}
	location.Coordinate = fetchGeocode(location)
	change := bson.M{"$set": bson.M{"address": location.Address, "city": location.City, "state": location.State, "zip": location.Zip, "coordinates": location.Coordinate}}
	err = c.Update(colQuerier, change)
	if err != nil {
		return loc, err
	}
	result := Location{}
	err = c.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		return loc, err
	}
	return result, nil
}

// Services related to Uber used for Trip planning and requesting
type UberService struct {
}

func (u *UberService) findBestRoute(request TripsRequest) (Plan, error) {
	l := LocationController{}

	ids := request.Location_ids
	ids = append(ids, request.Starting_from_location_id)

	// Fetch all locations from database
	output, err := l.GetLocationByIds(ids)
	if err != nil {
		return Plan{}, err
		//		fmt.Println(err)
	}

	// Get start_location and destinations.
	// Mongodb may return result in random order. Need to find exact match.
	var start_location Location
	var destinations []Location
	for i := 0; i < len(output); i++ {
		if output[i].Id == request.Starting_from_location_id {
			start_location = output[i]
			destinations = append(output[:i], output[i+1:]...)
			break
		}
	}

	//	fmt.Println(start_location)
	//	fmt.Println(destinations)

	// Find shortest path
	best_path := find_shortest_path(start_location, destinations)

	// Calculations for best path
	best_route_location_ids := []int{}
	total_uber_costs := 0
	total_uber_duration := 0
	total_distance := 0.0
	for i := 0; i < len(best_path)-1; i++ {
		best_route_location_ids = append(best_route_location_ids, best_path[i].Destination.Id)
		total_uber_costs += best_path[i].Low_estimate
		total_uber_duration += best_path[i].Duration
		total_distance += best_path[i].Distance
	}

	// Calculations for last leg of trip to make it round trip
	total_uber_costs += best_path[len(best_path)-1].Low_estimate
	total_uber_duration += best_path[len(best_path)-1].Duration
	total_distance += best_path[len(best_path)-1].Distance

	plan := Plan{}
	plan.Best_route_location_ids = best_route_location_ids
	plan.Starting_from_location_id = request.Starting_from_location_id
	plan.Status = "Planning"
	plan.Total_distance = total_distance
	plan.Total_uber_costs = total_uber_costs
	plan.Total_uber_duration = total_uber_duration
	plan.Next_destination_index = -1
	plan.Next_destination_location_id = -1

	// Keeping San Jose UberX as default product id. This ID should be stored in DB with every plan.
	plan.Product_id = "04a497f5-380d-47f2-bf1b-ad4cfdcb51f2"
	plan, err = u.StorePlan(plan)
	return plan, err
}

func (u *UberService) StorePlan(plan Plan) (Plan, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return plan, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("plan")
	plan.Id, err = generateID(*session, "plan")
	if err != nil {
		return plan, err
	}
	err = c.Insert(plan)
	if err != nil {
		return plan, err
	}
	return plan, nil
}

func (u *UberService) GetPlan(id int) (Plan, error) {
	result := Plan{}
	session, err := mgo.Dial(url)
	if err != nil {
		return result, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("plan")
	err = c.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		return result, err
	}
	//	fmt.Println(result)
	return result, nil
}

func (u *UberService) StartTrip(id int) (Plan, error) {
	plan, err := u.GetPlan(id)
	if err != nil {
		return plan, err
	}
	if plan.Status == "Completed" {
		return plan, errors.New("Trip already completed!")
	}
	l := LocationController{}
	locationIds := []int{-1, -1}

	//	fmt.Println(plan.Next_destination_index)
	//	fmt.Println(plan.Best_route_location_ids)
	status := ""
	if plan.Next_destination_index == -1 {
		locationIds[0] = plan.Starting_from_location_id
		locationIds[1] = plan.Best_route_location_ids[0]
		status = "Requesting"
	} else if plan.Next_destination_index == len(plan.Best_route_location_ids)-1 {
		locationIds[0] = plan.Best_route_location_ids[plan.Next_destination_index-1]
		locationIds[1] = plan.Starting_from_location_id
		status = "Completed"
	} else {
		locationIds[0] = plan.Best_route_location_ids[plan.Next_destination_index]
		locationIds[1] = plan.Best_route_location_ids[plan.Next_destination_index+1]
		status = "Requesting"
	}
	index := plan.Next_destination_index
	locs, err := l.GetLocationByIds(locationIds)

	if err != nil {
		return plan, err
	}
	var start_loc, end_loc Location
	if locs[0].Id == locationIds[0] {
		start_loc = locs[0]
		end_loc = locs[1]
	} else {
		start_loc = locs[1]
		end_loc = locs[0]
	}
	//	fmt.Println("STARTING FROM:::::: ", start_loc)
	//	fmt.Println("END LOCATION::::::", end_loc)

	// Request UBER
	uberRequest, err := request_uber(&start_loc, &end_loc, plan.Product_id)
	if err != nil {
		return plan, err
	}
	plan.Next_destination_location_id = index + 1

	//	fmt.Println(plan)

	// Update stored Plan
	session, err := mgo.Dial(url)
	if err != nil {
		return plan, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("plan")
	colQuerier := bson.M{"_id": plan.Id}
	change := bson.M{"$set": bson.M{"next_destination_index": index + 1, "next_destination_location_id": end_loc.Id, "status": status}}
	err = c.Update(colQuerier, change)
	if err != nil {
		return plan, err
	}
	err = c.Find(bson.M{"_id": id}).One(&plan)
	if err != nil {
		return plan, err
	}
	//	fmt.Println("PLAN:::::::::::::::", plan)
	plan.Uber_wait_time_eta = uberRequest.ETA
	return plan, nil
}

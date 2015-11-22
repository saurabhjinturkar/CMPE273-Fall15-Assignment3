# CMPE273-Fall15-Assignment3
###Problem Statement
Please check CMPE273-Fall15-Assignment3.pdf
###Input
The service supports trips API along with already implemented Location API. 
- Trips can be created by calling ```POST http://localhost:8081/trips/```
- Trip information for particular trip id can be get by calling ```GET http://localhost:8081/trips/{trip_id}```
- Uber API request for current location(from plan) to next destination from plan can be placed by ```PUT http://localhost:8081/trips/{trip_id}/request```

package main

import (
	"testing"
)

func TestFetchRoute(t *testing.T) {
	from := Coordinate{37.335877, -121.887775}
	to := Coordinate{37.363947, -121.928938}
	fetchRoute(from, to)
	t.Error("Fetch Route Error");
}

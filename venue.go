package pickem

import "cloud.google.com/go/firestore"

// A Venue represents a place where a Matchup is played.
type Venue struct {
	Name            string                   `json:"name" firestore:"name"`
	Capacity        int                      `json:"capacity" firestore:"capacity"`
	Grass           bool                     `json:"grass" firestore:"grass"`
	City            string                   `json:"city" firestore:"city"`
	State           string                   `json:"state" firestore:"state"`
	Zip             string                   `json:"zip" firestore:"zip"`
	LatLonAlt       []float64                `json:"lat_lon_alt" firestore:"lat_lon_alt"`
	YearConstructed int                      `json:"year_constructed" firestore:"year_constructed"`
	Dome            bool                     `json:"dome" firestore:"dome"`
	HomeTeams       []*firestore.DocumentRef `json:"home_teams" firestore:"home_teams"`
}

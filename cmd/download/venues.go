package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/reallyasi9/pickem"
)

var venuesFlagSet flag.FlagSet

func init() {
	commands["venues"] = venues

	venuesFlagSet.BoolVar(&dryRunFlag, "dryrun", false, "download and print actions only (do not upload to Firestore)")
	venuesFlagSet.BoolVar(&overwriteFlag, "overwrite", false, "overwrite documents in Firestore if they already exist")
}

type cfbdVenue struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Capacity    int     `json:"capacity"`
	Grass       bool    `json:"grass"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Zip         *string `json:"zip"`
	CountryCode string  `json:"country_code"`
	Location    *struct {
		Latitude  float64 `json:"x"`
		Longitude float64 `json:"y"`
	} `json:"location"`
	Elevation       *string `json:"elevation"` // needs to be converted to float64
	YearConstructed int     `json:"year_constructed"`
	Dome            bool    `json:"dome"`
}

func (v cfbdVenue) pickem() (*pickem.Venue, error) {
	var pv pickem.Venue
	pv.Name = v.Name
	pv.Capacity = v.Capacity
	pv.Grass = v.Grass
	pv.City = v.City
	pv.State = v.State
	if v.Zip == nil {
		pv.Zip = ""
	} else {
		pv.Zip = *v.Zip
	}
	pv.YearConstructed = v.YearConstructed
	pv.HomeTeams = make([]*firestore.DocumentRef, 0)

	if v.Location == nil {
		return &pv, nil
	}

	var ele float64
	var err error
	if v.Elevation == nil {
		ele = 0.
	} else if ele, err = strconv.ParseFloat(*v.Elevation, 64); err != nil {
		return nil, err
	}
	pv.LatLonAlt = make([]float64, 3)
	pv.LatLonAlt[0] = v.Location.Latitude
	pv.LatLonAlt[1] = v.Location.Longitude
	pv.LatLonAlt[2] = ele

	return &pv, nil
}

func venues(ctx context.Context, args []string) error {
	if err := venuesFlagSet.Parse(args); err != nil {
		return err
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return err
	}
	u.Path = "/venues"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var venues []cfbdVenue
	err = json.Unmarshal(body, &venues)
	if err != nil {
		return err
	}

	toWrite := newFSCommitter(fs, 500)
	collection := fs.Collection("xvenues")

	var v cfbdVenue
	for _, v = range venues {
		venue, err := v.pickem()
		if err != nil {
			return err
		}
		ref := collection.Doc(strconv.Itoa(v.ID))
		if dryRunFlag {
			fmt.Printf("%s <- %v\n", ref.ID, venue)
			continue
		}

		if overwriteFlag {
			if err := toWrite.Set(ctx, ref, &venue); err != nil {
				return err
			}
		} else {
			if err := toWrite.Create(ctx, ref, &venue); err != nil {
				return err
			}
		}
	}

	if !dryRunFlag {
		if err := toWrite.Commit(ctx); err != nil {
			return err
		}
	}

	return nil

}

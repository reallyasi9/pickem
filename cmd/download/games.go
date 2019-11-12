package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/reallyasi9/pickem"
	"google.golang.org/api/iterator"
)

var gamesFlagSet flag.FlagSet
var gamesYearFlag int
var gamesWeekFlag int
var gamesTeamFlag string
var gamesConferenceFlag string
var gamesSeasonTypeFlag seasonType
var gamesUpdateTeamVenueFlag bool

type seasonType string

const (
	regularSeason seasonType = "regular"
	postseason    seasonType = "postseason"
)

// String is the method to format the flag's value, part of the flag.Value interface.
func (s *seasonType) String() string {
	return string(*s)
}

// Set is the method to set the flag value, part of the flag.Value interface.
func (s *seasonType) Set(value string) error {
	v := seasonType(value)
	switch v {
	case regularSeason:
	case postseason:
	default:
		return fmt.Errorf("'%s' is not a season type", value)
	}
	*s = v
	return nil
}

func init() {
	commands["games"] = games

	gamesFlagSet.StringVar(&gamesConferenceFlag, "conference", "", "conference download filter")
	gamesFlagSet.IntVar(&gamesYearFlag, "year", time.Now().Year(), "year to download")
	gamesFlagSet.IntVar(&gamesWeekFlag, "week", 0, "week download filter (starting with 1, any number < 1 will download all weeks in the season)")
	gamesFlagSet.StringVar(&gamesTeamFlag, "team", "", "team download filter")
	gamesFlagSet.Var(&gamesSeasonTypeFlag, "type", "season type download filter (regular or postseason)")
	gamesFlagSet.BoolVar(&gamesUpdateTeamVenueFlag, "updateVenues", false, "update home team venues using game information")
	gamesFlagSet.BoolVar(&dryRunFlag, "dryrun", false, "download and print actions only (do not upload to Firestore)")
	gamesFlagSet.BoolVar(&overwriteFlag, "overwrite", false, "overwrite documents in Firestore if they already exist")
}

type cfbdGame struct {
	ID             int        `json:"id"`
	Season         int        `json:"season"`
	Week           int        `json:"week"`
	SeasonType     seasonType `json:"season_type"`
	StartDate      time.Time  `json:"start_date"`
	NeutralSite    bool       `json:"neutral_site"`
	ConferenceGame bool       `json:"conference_game"`
	Attendance     *int       `json:"attendance"`
	VenueID        int        `json:"venue_id"`
	Venue          string     `json:"venue"`
	HomeTeam       string     `json:"home_team"`
	HomeConference *string    `json:"home_conference"`
	HomePoints     *int       `json:"home_points"`
	HomeLineScores []int      `json:"home_line_scores"`
	AwayTeam       string     `json:"away_team"`
	AwayConference *string    `json:"away_conference"`
	AwayPoints     *int       `json:"away_points"`
	AwayLineScores []int      `json:"away_line_scores"`
}

func (g cfbdGame) pickem() (*pickem.Game, error) {
	var pg pickem.Game

	pg.Season = fs.Collection("seasons").Doc(strconv.Itoa(g.Season))
	pg.Week = g.Week
	pg.Postseason = g.SeasonType == postseason
	pg.StartTime = g.StartDate
	pg.NeutralSite = g.NeutralSite
	pg.ConferenceGame = g.ConferenceGame
	pg.Attendance = g.Attendance
	pg.Venue = fs.Collection("xvenues").Doc(strconv.Itoa(g.VenueID))
	var ok bool
	if pg.HomeTeam, ok = bySchool[g.HomeTeam]; !ok {
		return nil, fmt.Errorf("ID of team '%s' not found", g.HomeTeam)
	}
	pg.HomePoints = g.HomePoints
	pg.HomeLineScores = g.HomeLineScores
	if pg.AwayTeam, ok = bySchool[g.AwayTeam]; !ok {
		return nil, fmt.Errorf("ID of team '%s' not found", g.AwayTeam)
	}
	pg.AwayPoints = g.AwayPoints
	pg.AwayLineScores = g.AwayLineScores

	return &pg, nil
}

var bySchool map[string]*firestore.DocumentRef

func fillSchools(ctx context.Context) error {
	bySchool = make(map[string]*firestore.DocumentRef)
	itr := fs.Collection("xteams").Documents(ctx)
	defer itr.Stop()
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		school, err := doc.DataAt("school_name")
		if err != nil {
			return err
		}
		if _, ok := bySchool[school.(string)]; ok {
			return fmt.Errorf("school '%s' appears in data more than once", school.(string))
		}
		bySchool[school.(string)] = doc.Ref
	}
	return nil
}

func games(ctx context.Context, args []string) error {
	if err := gamesFlagSet.Parse(args); err != nil {
		return err
	}

	if err := fillSchools(ctx); err != nil {
		return err
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return err
	}
	u.Path = "/games"

	q := u.Query()
	q.Set("conference", gamesConferenceFlag)
	q.Set("year", strconv.Itoa(gamesYearFlag))
	q.Set("seasonType", string(gamesSeasonTypeFlag))
	if gamesWeekFlag >= 1 {
		q.Set("week", strconv.Itoa(gamesWeekFlag))
	}
	q.Set("team", gamesTeamFlag)
	u.RawQuery = q.Encode()

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

	var games []cfbdGame
	err = json.Unmarshal(body, &games)
	if err != nil {
		return err
	}

	toWrite := newFSCommitter(fs, 250)
	collection := fs.Collection("xgames")
	byHomeTeam := make(map[string]*firestore.DocumentRef)

	var g cfbdGame
	for _, g = range games {
		game, err := g.pickem()
		if err != nil {
			return err
		}
		ref := collection.Doc(strconv.Itoa(g.ID))

		if gamesUpdateTeamVenueFlag && !g.NeutralSite {
			if v, ok := byHomeTeam[game.HomeTeam.ID]; ok && v.ID != game.Venue.ID {
				log.Printf("warning: team %v has multiple home venues\n", game.HomeTeam.ID)
			}
			byHomeTeam[game.HomeTeam.ID] = game.Venue
		}

		if dryRunFlag {
			fmt.Printf("%s <- %v\n", ref.ID, game)
			continue
		}

		if overwriteFlag {
			if err := toWrite.Set(ctx, ref, &game); err != nil {
				return err
			}
		} else {
			if err := toWrite.Create(ctx, ref, &game); err != nil {
				return err
			}
		}
	}

	if err := toWrite.Commit(ctx); err != nil {
		return err
	}

	if gamesUpdateTeamVenueFlag {
		itr := fs.Collection("xteams").Documents(ctx)
		defer itr.Stop()
		for {
			doc, err := itr.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var venue *firestore.DocumentRef
			var ok bool
			if venue, ok = byHomeTeam[doc.Ref.ID]; !ok {
				continue
			}
			if dryRunFlag {
				fmt.Printf("%s <- %s\n", doc.Ref.ID, venue.ID)
				continue
			}
			if err := toWrite.Update(ctx, doc.Ref, []firestore.Update{{Path: "home_venue", Value: venue}}); err != nil {
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

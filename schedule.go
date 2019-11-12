package pickem

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// byWeek sorts a collection of games by week
type byWeek []*Game

func (b byWeek) Len() int           { return len(b) }
func (b byWeek) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byWeek) Less(i, j int) bool { return b[i].Week < b[j].Week }

// Matchups converts a team's schedule into Matchups.
func Matchups(ctx context.Context, fs *firestore.Client, team string, season int, startWeek int) ([]*Matchup, error) {
	var t *Team
	var err error
	if t, err = LookupTeam(ctx, fs, team); err != nil {
		return nil, err
	}

	seasonRef := fs.Collection("seasons").Doc(strconv.Itoa(season))
	teamRef := fs.Collection("xteams").Doc(t.SchoolName)

	games := make([]*Game, 0)
	fmt.Println(seasonRef.Path)
	fmt.Println(startWeek)
	fmt.Println(teamRef.Path)
	homeGameItr := fs.Collection("xgames").Where("season", "==", seasonRef).Where("week", ">=", startWeek).Where("home_team", "==", teamRef).Documents(ctx)
	defer homeGameItr.Stop()
	for {
		doc, err := homeGameItr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var game Game
		if err := doc.DataTo(&game); err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	awayGameItr := fs.Collection("xgames").Where("season", "==", seasonRef).Where("week", ">=", startWeek).Where("away_team", "==", teamRef).Documents(ctx)
	defer awayGameItr.Stop()
	for {
		doc, err := awayGameItr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var game Game
		if err := doc.DataTo(&game); err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	sort.Sort(byWeek(games))
	matchups := make([]*Matchup, len(games))
	for i, game := range games {
		var t1, t2 Team
		doc, err := game.HomeTeam.Get(ctx)
		if err != nil {
			return nil, err
		}
		if err := doc.DataTo(&t1); err != nil {
			return nil, err
		}
		doc, err = game.AwayTeam.Get(ctx)
		if err != nil {
			return nil, err
		}
		if err := doc.DataTo(&t2); err != nil {
			return nil, err
		}
		loc := Home
		if game.NeutralSite {
			loc = Neutral
		}
		matchups[i] = &Matchup{Team1: &t1, Team2: &t2, Location: loc}
	}

	return matchups, nil

}

func loc(locTeam string) RelativeLocation {
	// Note: this is relative to the schedule team, not the team given here.
	switch locTeam[0] {
	case '@':
		return Away
	case '>':
		return Far
	case '<':
		return Near
	case '!':
		return Neutral
	default:
		return Home
	}
}

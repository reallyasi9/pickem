package pickem

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// TeamSchedule is a Team's schedule for a Season.
type TeamSchedule struct {
	Team      *firestore.DocumentRef   `json:"team",firestore:"team"`
	Opponents []*firestore.DocumentRef `json:"opponents",firestore:"opponents"`
	Locations []RelativeLocation       `json:"locations",firestore:"locations"`
}

// Schedule is a complete schedule for all teams for a Season.
type Schedule struct {
	Season    *firestore.DocumentRef `json:"season",firestore:"season"`
	Schedules []TeamSchedule         `json:"schedules",firestore:"schedules"`
}

// Matchups converts a TeamSchedule into Matchups.
func (ts *TeamSchedule) Matchups(ctx context.Context, fs *firestore.Client) ([]*Matchup, error) {
	if len(ts.Opponents) != len(ts.Locations) {
		return nil, fmt.Errorf("opponents and locations must be the same length")
	}
	teamDoc, err := ts.Team.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get team %s from firestore: %v", ts.Team.ID, err)
	}
	var t1 Team
	err = teamDoc.DataTo(&t1)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal team %s: %v", ts.Team.ID, err)
	}
	oppDocs, err := fs.GetAll(ctx, ts.Opponents)
	if err != nil {
		return nil, fmt.Errorf("unable to get opponents from firestore: %v", err)
	}
	mus := make([]*Matchup, len(oppDocs))
	for i, doc := range oppDocs {
		var t2 Team
		err = doc.DataTo(&t2)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal team %s: %v", doc.Ref.ID, err)
		}
		mus[i] = &Matchup{team1: &t1, team2: &t2, location: ts.Locations[i]}
	}
	return mus, nil
}

// MatchupMap converts a Schedule into a team => Matchup map.
func (s *Schedule) MatchupMap(ctx context.Context, fs *firestore.Client) (map[*Team][]*Matchup, error) {
	mm := make(map[*Team][]*Matchup)

	for _, ts := range s.Schedules {
		mus, err := ts.Matchups(ctx, fs)
		if err != nil {
			return nil, err
		}
		if len(mus) == 0 {
			continue
		}
		t1 := mus[0].team1
		if _, ok := mm[t1]; ok {
			return nil, fmt.Errorf("team %s appears more than once in schedules for season %s", t1.Name(), s.Season.ID)
		}
		mm[t1] = mus
	}

	return mm, nil
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

package pickem

import (
	"context"

	"cloud.google.com/go/firestore"
)

/*Predictable is an interface that returns a prediction if given a GamePredicter.  Games implement Predictable. */
type Predictable interface {
	Predict(MatchupPredicter) (prob float64, spread float64, err error)
}

// RelativeLocation describes where a game is being played relative to one team's home field.
type RelativeLocation int

const (
	// Home is a team's home field.
	Home RelativeLocation = 2

	// Near is a field closer to a given team's home field than the team's opponent's home field.
	Near = 1

	// Neutral is a truely neutral location.
	Neutral = 0

	// Far is a field closer to a given team's opponent's home field than the team's home field.
	Far = -1

	// Away is a team's opponent's home field.
	Away = -2
)

func (rl RelativeLocation) String() string {
	switch rl {
	case -2:
		return "Away"
	case -1:
		return "Far"
	case 0:
		return "Neutral"
	case 1:
		return "Near"
	case 2:
		return "Home"
	}
	return "Unknown"
}

// Matchup represents a matchup between two teams.
type Matchup struct {
	Team1    *Team
	Team2    *Team
	Location RelativeLocation
}

// MatchupRef is a representation of a matchup in Firestore.
type MatchupRef struct {
	Team1    *firestore.DocumentRef `json:"team1" firestore:"team1"`
	Team2    *firestore.DocumentRef `json:"team2" firestore:"team2"`
	Location RelativeLocation       `json:"location" firestore:"location"`
}

// GetTo implements FirestoreGetter and fills MatchupRef with primatives.
func (mr *MatchupRef) GetTo(ctx context.Context, fs *firestore.Client, m *Matchup) error {
	var t1, t2 Team
	var ts *firestore.DocumentSnapshot
	var err error
	if ts, err = mr.Team1.Get(ctx); err != nil {
		return err
	}
	if err = ts.DataTo(&t1); err != nil {
		return err
	}
	if ts, err = mr.Team2.Get(ctx); err != nil {
		return err
	}
	if err = ts.DataTo(&t2); err != nil {
		return err
	}
	m.Team1 = &t1
	m.Team2 = &t2
	m.Location = mr.Location
	return nil
}

// NewMatchup makes a game between two teams.
func NewMatchup(team1, team2 *Team, locRelTeam1 RelativeLocation) *Matchup {
	return &Matchup{Team1: team1, Team2: team2, Location: locRelTeam1}
}

// RelativeLocation returns the location of the game relative to the given team.
func (g *Matchup) RelativeLocation(relToFirstTeam bool) RelativeLocation {
	if relToFirstTeam {
		return g.Location
	}
	return -g.Location
}

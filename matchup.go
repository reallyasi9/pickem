package pickem

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
	team1    *Team
	team2    *Team
	location RelativeLocation
}

// NewMatchup makes a game between two teams.
func NewMatchup(team1, team2 *Team, locRelTeam1 RelativeLocation) *Matchup {
	return &Matchup{team1: team1, team2: team2, location: locRelTeam1}
}

// Location returns the location of the game relative to the given team.
func (g *Matchup) Location(relToFirstTeam bool) RelativeLocation {
	if relToFirstTeam {
		return g.location
	}
	return -g.location
}

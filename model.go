package pickem

import (
	"fmt"

	"github.com/atgjack/prob"
)

/*MatchupPredicter describes an object that can predict the probability of win if two teams play each other.

The RelativeLocation argument is relative to the first Team, so a value of Home means the first Team is
	playing at home while the second team is playing on the road.

The returned probability is relative to the first Team argument, so a probability of .9 means the first Team
	has a 90% probability of win and the second Team has a 10% probability of win.

The returned spread is positive if the first Team is predicted to win, negative if the second Team is predicted
	to win, and zero if there is no favorite.*/
type MatchupPredicter interface {
	Predict(Team, Team, RelativeLocation) (prob float64, spread float64, err error)
}

/*GaussianSpreadModel implements GamePredicter and uses a normal distribution based on spreads to calculate
win probabilities.

The spread is determined by team rating difference and where the game is being played (to account for bias).*/
type GaussianSpreadModel struct {
	dist      prob.Normal
	homeBias  float64
	closeBias float64
	ratings   map[Team]float64
}

/*NewGaussianSpreadModel makes a model.
Note that positive homeBias and closeBias are points added to the home/close team's predicted spread.*/
func NewGaussianSpreadModel(ratings map[Team]float64, stdDev, homeBias, closeBias float64) *GaussianSpreadModel {
	return &GaussianSpreadModel{ratings: ratings, dist: prob.Normal{Mu: 0, Sigma: stdDev}, homeBias: homeBias, closeBias: closeBias}
}

// Predict returns the probability and spread for team1.
func (m *GaussianSpreadModel) Predict(t1, t2 Team, loc RelativeLocation) (float64, float64, error) {
	if t1 == NONE {
		// The second team has a bye week, so wins automatically.
		return 0., 0., nil
	}
	if t2 == NONE {
		// The first team has a bye week, so wins automatically.
		return 1., 0., nil
	}
	spread, err := m.spread(t1, t2, loc)
	if err != nil {
		return 0., 0., fmt.Errorf("Predict failed to calculate spread: %v", err)
	}
	prob := m.dist.Cdf(spread)

	return prob, spread, nil
}

func (m GaussianSpreadModel) spread(t1, t2 Team, loc RelativeLocation) (float64, error) {
	r1, ok := m.ratings[t1]
	if !ok {
		return 0., fmt.Errorf("team 1 '%s' has no rating", t1.Name())
	}
	r2, ok := m.ratings[t2]
	if !ok {
		return 0., fmt.Errorf("team 2 '%s' has no rating", t2.Name())
	}

	diff := r1 - r2
	switch loc {
	case Home:
		diff += m.homeBias
	case Near:
		diff += m.closeBias
	case Far:
		diff -= m.closeBias
	case Away:
		diff -= m.homeBias
	}
	return diff, nil
}

// A teamPair is just a way to store two teams in a lookup table and allow fast searching by teams in either order.
type teamPair struct {
	team1 Team
	team2 Team
}

// A matchupMap allows searching for matchup spreads with teams in either order.
type matchupMap map[teamPair]float64

// Get searches for a matchup in the matchupMap
func (mm matchupMap) get(t1, t2 Team) (spread float64, swap bool, ok bool) {
	m := teamPair{t1, t2}
	if spread, ok = mm[m]; ok {
		return
	}
	m = teamPair{t2, t1}
	if spread, ok = mm[m]; ok {
		swap = true
		return
	}
	// Not found.  That's a shame.
	return
}

/*LookupModel implements GamePredicter and uses a simple lookup table to calculate spreads and a gaussian model
to calculate win probabilities.*/
type LookupModel struct {
	dist      prob.Normal
	homeBias  float64
	closeBias float64
	spreads   matchupMap
}

// NewLookupModel makes a model.
func NewLookupModel(homeTeams, roadTeams []Team, spreads []float64, stdDev, homeBias, closeBias float64) *LookupModel {
	if len(homeTeams) != len(roadTeams) || len(homeTeams) != len(spreads) {
		panic(fmt.Errorf("mismatched length of home (%d), road (%d), and spread (%d) slices", len(homeTeams), len(roadTeams), len(spreads)))
	}
	mm := make(map[teamPair]float64)
	for i := 0; i < len(homeTeams); i++ {
		mm[teamPair{homeTeams[i], roadTeams[i]}] = spreads[i]
	}
	return &LookupModel{spreads: mm, dist: prob.Normal{Mu: 0, Sigma: stdDev}, homeBias: homeBias, closeBias: closeBias}
}

// Predict returns the probability and spread for team1.
func (m *LookupModel) Predict(t1, t2 Team, loc RelativeLocation) (float64, float64, error) {
	if t1 == NONE {
		// The second team has a bye week, so wins automatically.
		return 0., 0., nil
	}
	if t2 == NONE {
		// The first team has a bye week, so wins automatically.
		return 1., 0., nil
	}
	spread, swap, ok := m.spreads.get(t1, t2)
	if !ok {
		return 0., 0., fmt.Errorf("spread between teams %s and %s not found", t1.Name(), t2.Name())
	}
	if swap {
		spread = -spread
	}
	switch loc {
	case Home:
		spread += m.homeBias
	case Near:
		spread += m.closeBias
	case Far:
		spread -= m.closeBias
	case Away:
		spread -= m.homeBias
	}

	prob := m.dist.Cdf(spread)

	return prob, spread, nil
}

package pickem

import (
	"time"

	"cloud.google.com/go/firestore"
)

// A Game represents a matchup between two Teams, the conditions of the matchup (time, locale, etc.), and the outcomes.
type Game struct {
	Season         *firestore.DocumentRef `json:"season"`
	Week           int                    `json:"week"`
	Postseason     bool                   `json:"postseason"`
	StartTime      time.Time              `json:"start_time"`
	NeutralSite    bool                   `json:"neutral_site"`
	ConferenceGame bool                   `json:"conference_game"`
	Attendance     *int                   `json:"attendance"`
	Venue          *firestore.DocumentRef `json:"venue"`
	HomeTeam       *firestore.DocumentRef `json:"home_team"`
	HomePoints     *int                   `json:"home_points"`
	HomeLineScores []int                  `json:"home_line_scores"`
	AwayTeam       *firestore.DocumentRef `json:"away_team"`
	AwayPoints     *int                   `json:"away_points"`
	AwayLineScores []int                  `json:"away_line_scores"`
}

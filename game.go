package pickem

import (
	"time"

	"cloud.google.com/go/firestore"
)

// A Game represents a matchup between two Teams, the conditions of the matchup (time, locale, etc.), and the outcomes.
type Game struct {
	Season         *firestore.DocumentRef `json:"season" firestore:"season"`
	Week           int                    `json:"week" firestore:"week"`
	Postseason     bool                   `json:"postseason" firestore:"postseason"`
	StartTime      time.Time              `json:"start_time" firestore:"start_time"`
	NeutralSite    bool                   `json:"neutral_site" firestore:"neutral_site"`
	ConferenceGame bool                   `json:"conference_game" firestore:"conference_game"`
	Attendance     *int                   `json:"attendance" firestore:"attendance"`
	Venue          *firestore.DocumentRef `json:"venue" firestore:"venue"`
	HomeTeam       *firestore.DocumentRef `json:"home_team" firestore:"home_team"`
	HomePoints     *int                   `json:"home_points" firestore:"home_points"`
	HomeLineScores []int                  `json:"home_line_scores" firestore:"home_line_scores"`
	AwayTeam       *firestore.DocumentRef `json:"away_team" firestore:"away_team"`
	AwayPoints     *int                   `json:"away_points" firestore:"away_points"`
	AwayLineScores []int                  `json:"away_line_scores" firestore:"away_line_scores"`
	Timestamp      time.Time              `json:"timestamp" firestore:"timestamp,serverTimestamp"`
}

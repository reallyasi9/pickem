package pickem

import (
	"time"

	"cloud.google.com/go/firestore"
)

// Player represents a player's current status in the competition.
type Player struct {
	Name     string    `json:"name",firestore:"name"`
	JoinDate time.Time `json:"join_date",firestore:"join_date"`
}

// PlayerPreferences holds preferred options for the player.
type PlayerPreferences struct {
	FavoriteTeam       *firestore.DocumentRef `json:"favorite_team",firestore:"favorite_team"`
	StraightUpModel    *firestore.DocumentRef `json:"straight_up_model",firestore:"straight_up_model"`
	NoisySpreadModel   *firestore.DocumentRef `json:"noisy_spread_model",firestore:"noisy_spread_model"`
	SuperdogModel      *firestore.DocumentRef `json"superdog_model",firefirestore:"superdog_model"`
	PonyModel          *firestore.DocumentRef `json:"pony_model",firestore:"pony_model"`
	BeatTheStreakModel *firestore.DocumentRef `json:"beat_the_streak_model",firestore:"beat_the_streak_model"`
}

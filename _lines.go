package pickem

import (
	"encoding/csv"
	"math"
	"net/http"
	"strconv"
)

// MatchupModel is a combined game and model for 2D lookup of lines
type MatchupModel struct {
	Game  *Matchup
	Model string
}

// LineMap is a mapping of game/model combinations with a line
type LineMap map[MatchupModel]float64

// MakeLines makes a map of games to lines
func MakeLines(url string) (LineMap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	lines := make(LineMap)

	// first line contains the header information
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	models := header[2:len(header)]

	record, err := reader.Read()
	for ; record != nil; record, err = reader.Read() {
		if record == nil {
			break
		}
		if err != nil {
			return nil, err
		}

		game := NewMatchup(Team{Name4: record[0]}, Team{Name4: record[1]}, Neutral)

		for i, line := range record[2:len(record)] {
			model := models[i]
			gm := MatchupModel{Game: game, Model: model}

			val, err := strconv.ParseFloat(line, 64)
			if err != nil {
				val = math.NaN() // Not an error, just missing data
			}
			lines[gm] = val
		}
	}

	return lines, nil
}

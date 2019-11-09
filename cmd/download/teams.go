package main

import (
	"flag"

	"github.com/reallyasi9/pickem"
)

var teamsFlagSet flag.FlagSet
var teamsConferenceFlag string
var teamsDryRun bool

func init() {
	commands["teams"] = teams

	teamsFlagSet.StringVar(&teamsConferenceFlag, "conference", "", "Conference of teams to download")
	teamsFlagSet.BoolVar(&teamsDryRun, "dryrun", false, "Do not upload to firestore")
}

type cfbdTeam struct {
	ID           int      `json:"id"`
	School       string   `json:"school"`
	Mascot       string   `json:"mascot"`
	Abbreviation string   `json:"abbreviation"`
	AltName1     *string  `json:"alt_name_1"`
	AltName2     *string  `json:"alt_name_2"`
	AltName3     *string  `json:"alt_name_3"`
	Conference   *string  `json:"conference"`
	Division     *string  `json:"division"`
	Color        string   `json:"color"`
	AltColor     *string  `json:"alt_color"`
	Logos        []string `json:"logos"`
}

func (t cfbdTeam) pickem() (*pickem.Team, error) {
	return nil, nil
}

func teams(args []string) error {
	if err := teamsFlagSet.Parse(args); err != nil {
		return err
	}

	return nil
}

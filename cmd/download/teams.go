package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/reallyasi9/pickem"
)

var teamsFlagSet flag.FlagSet
var teamsConferenceFlag string

func init() {
	commands["teams"] = teams

	teamsFlagSet.StringVar(&teamsConferenceFlag, "conference", "", "conference download filter")
	teamsFlagSet.BoolVar(&dryRunFlag, "dryrun", false, "download and print actions only (do not upload to Firestore)")
	teamsFlagSet.BoolVar(&overwriteFlag, "overwrite", false, "overwrite documents in Firestore if they already exist")
}

type cfbdTeam struct {
	ID           int      `json:"id"`
	School       string   `json:"school"`
	Mascot       *string  `json:"mascot"`
	Abbreviation *string  `json:"abbreviation"`
	AltName1     *string  `json:"alt_name_1"`
	AltName2     *string  `json:"alt_name_2"`
	AltName3     *string  `json:"alt_name_3"`
	Conference   *string  `json:"conference"`
	Division     *string  `json:"division"`
	Color        *string  `json:"color"`
	AltColor     *string  `json:"alt_color"`
	Logos        []string `json:"logos"`
}

func (t cfbdTeam) pickem() (*pickem.Team, error) {
	// If there is no school, I can't do anything
	if t.School == "" {
		return nil, fmt.Errorf("team %d missing school name", t.ID)
	}
	var team pickem.Team
	team.Names = make([]string, 0)
	team.Colors = make([]pickem.RGBHex, 0)
	team.Logos = t.Logos

	team.Names = append(team.Names, t.School)
	team.SchoolName = t.School
	// "State" is always abbreviated both "St." and "St"
	if strings.Contains(team.SchoolName, "State") {
		team.Names = append(team.Names, strings.Replace(team.SchoolName, "State", "St.", 1))
		team.Names = append(team.Names, strings.Replace(team.SchoolName, "State", "St", 1))
	}
	if t.Abbreviation != nil {
		team.Names = append(team.Names, *t.Abbreviation)
		team.Abbreviation = *t.Abbreviation
	} else {
		// Need something here
		team.Abbreviation = strings.ToUpper(team.SchoolName)
	}
	if t.AltName1 != nil {
		team.Names = append(team.Names, *t.AltName1)
	}
	if t.AltName2 != nil {
		team.Names = append(team.Names, *t.AltName2)
	}
	if t.AltName3 != nil {
		team.Names = append(team.Names, *t.AltName3)
	}
	if t.Color != nil {
		team.Colors = append(team.Colors, pickem.RGBHex(*t.Color))
	}
	if t.AltColor != nil {
		team.Colors = append(team.Colors, pickem.RGBHex(*t.AltColor))
	}
	if t.Mascot != nil {
		team.TeamName = *t.Mascot
	}

	team.Conference = t.Conference
	team.Division = t.Division

	return &team, nil
}

func teams(ctx context.Context, args []string) error {
	flag.Parse()
	if err := teamsFlagSet.Parse(args); err != nil {
		return err
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return err
	}
	u.Path = "/teams"

	q := u.Query()
	q.Set("conference", teamsConferenceFlag)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var teams []cfbdTeam
	err = json.Unmarshal(body, &teams)
	if err != nil {
		return err
	}

	toWrite := fs.Batch()
	collection := fs.Collection("xteams")

	var i int
	var t cfbdTeam
	for i, t = range teams {
		team, err := t.pickem()
		if err != nil {
			return err
		}
		ref := collection.Doc(strconv.Itoa(t.ID))
		if dryRunFlag {
			fmt.Printf("%s <- %v\n", ref.ID, team)
			continue
		}

		if overwriteFlag {
			toWrite = toWrite.Set(ref, &team)
		} else {
			toWrite = toWrite.Create(ref, &team)
		}

		if i%500 == 499 {
			_, err := toWrite.Commit(ctx)
			if err != nil {
				return err
			}
			toWrite = fs.Batch()
		}
	}

	if !dryRunFlag && i%500 != 499 {
		_, err = toWrite.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

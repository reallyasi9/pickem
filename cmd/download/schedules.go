package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/reallyasi9/pickem"
)

var schedulesFlagSet flag.FlagSet
var schedulesYearFlag int
var schedulesWeekFlag int
var schedulesTeamFlag string

func init() {
	commands["schedules"] = schedules

	schedulesFlagSet.IntVar(&schedulesYearFlag, "year", time.Now().Year(), "year to download")
	schedulesFlagSet.IntVar(&schedulesWeekFlag, "week", 0, "week download filter (starting with 1, any number < 1 will download all weeks in the season)")
	schedulesFlagSet.StringVar(&schedulesTeamFlag, "team", "", "team download filter")
}

func schedules(ctx context.Context, args []string) error {
	if err := schedulesFlagSet.Parse(args); err != nil {
		return err
	}

	matches, err := pickem.Matchups(ctx, fs, schedulesTeamFlag, schedulesYearFlag, schedulesWeekFlag)
	if err != nil {
		return err
	}

	for _, m := range matches {
		fmt.Println(m)
	}

	return nil
}

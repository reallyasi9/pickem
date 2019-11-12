package pickem

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// A Team in Pick 'Em terms is an entity that participates in Matchups.  Teams have many representations:
// Names are a names that are used in Excel-style slates or by NCAA prediction models, and are determined by the whims of the authors;
// a SchoolName and TeamName are an attempt at a consistent naming scheme determined by the official school and team names published by the schools themselves.
// Other components of a team are for aesthetic purposes only.
type Team struct {
	ID           int                    `json:"id" firestore:"id"`
	Names        []string               `json:"names" firestore:"names"`
	Abbreviation string                 `json:"abbreviation" firestore:"abbreviation"`
	SchoolName   string                 `json:"school_name" firestore:"school_name"`
	TeamName     string                 `json:"team_name" firestore:"team_name"`
	Colors       []RGBHex               `json:"colors" firestore:"colors"`
	Logos        []string               `json:"logos" firestore:"logos"`
	Conference   *string                `json:"conference" firestore:"conference"`
	Division     *string                `json:"division" firestore:"division"`
	HomeVenue    *firestore.DocumentRef `firestore:"home_venue"`
}

// Name implements NameStringer interface.
func (t *Team) Name() string {
	return t.SchoolName + " " + t.TeamName
}

// ShortName implements NameStringer interface.
func (t *Team) ShortName() string {
	return t.Abbreviation
}

// LookupTeam looks up a team in Firestore by, in order, ID==SchoolName, Abbreviation, then Names.  If multiple teams match, an error is returned.
func LookupTeam(ctx context.Context, fs *firestore.Client, name string) (*Team, error) {
	collection := fs.Collection("xteams")
	var team Team

	doc, err := collection.Doc(name).Get(ctx)
	if err != nil {
		if err := doc.DataTo(&team); err != nil {
			return nil, err
		}
		return &team, nil
	}

	itr := collection.Where("Abbreviation", "==", name).Documents(ctx)
	defer itr.Stop()
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err := doc.DataTo(&team); err != nil {
			return nil, err
		}
		if _, err := itr.Next(); err != iterator.Done {
			return nil, fmt.Errorf("ambiguous team abbreviation '%s'", name)
		}
		return &team, nil
	}
	itr.Stop()

	itr = collection.Where("Names", "array-contains", name).Documents(ctx)
	defer itr.Stop()
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err := doc.DataTo(&team); err != nil {
			return nil, err
		}
		if _, err := itr.Next(); err != iterator.Done {
			return nil, fmt.Errorf("ambiguous team name '%s'", name)
		}
		return &team, nil
	}

	return nil, fmt.Errorf("name '%s' not found in teams", name)
}

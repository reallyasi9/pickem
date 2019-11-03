package pickem

// A Team in Pick 'Em terms is an entity that participates in Matchups.  Teams have many representations:
// Names are a names that are used in Excel-style slates or by NCAA prediction models, and are determined by the whims of the authors;
// a SchoolName and TeamName are an attempt at a consistent naming scheme determined by the official school and team names published by the schools themselves.
// Other components of a team are for aesthetic purposes only.
type Team struct {
	Names      []string `json:"names",firestore:"names"`
	SchoolName string   `json:"school_name",firestore:"school_name"`
	TeamName   string   `json:"team_name",firestore:"team_name"`
	MascotName string   `json:"mascot_name",firestore:"mascot_name"`
	Colors     []RGBHex `json:"colors",firestore:"colors"`
}

// Name implements NameStringer interface.
func (t *Team) Name() string {
	return t.SchoolName + " " + t.TeamName
}

// ShortName implements NameStringer interface.
func (t *Team) ShortName() string {
	shortest := t.Names[0]
	for i := 1; i < len(t.Names); i++ {
		if len(t.Names[i]) < len(shortest) {
			shortest = t.Names[i]
		}
	}
	return shortest
}

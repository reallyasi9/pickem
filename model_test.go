package pickem

import (
	"math"
	"reflect"
	"testing"

	"github.com/atgjack/prob"
)

func isClose(a, b, epsilon float64) bool {
	i, r := math.Modf(a - b)
	return i == 0 && math.Abs(r) < math.Abs(epsilon)
}

func isDarnClose(a, b float64) bool {
	n := math.Nextafter(a, b)
	if a > b {
		return n <= b
	}
	return n >= b
}

func TestGaussianSpreadModel_Predict(t *testing.T) {
	rm := make(map[Team]float64)
	rm[Team{"A"}] = 0.
	rm[Team{"B"}] = 1.

	winBy1sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(1)
	winBy2sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(2)

	winBy1sigma1 := prob.Normal{Mu: 0, Sigma: 1}.Cdf(1)

	type fields struct {
		dist      prob.Normal
		homeBias  float64
		closeBias float64
		ratings   map[Team]float64
	}
	type args struct {
		t1  Team
		t2  Team
		loc RelativeLocation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{name: "same team neutral",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"A"}, Neutral},
			want:   .5, want1: 0, wantErr: false},
		{name: "same team near",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"A"}, Near},
			want:   winBy1sigma12, want1: 1., wantErr: false},
		{name: "same team home",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"A"}, Home},
			want:   winBy2sigma12, want1: 2., wantErr: false},
		{name: "same team far",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"A"}, Far},
			want:   1 - winBy1sigma12, want1: -1., wantErr: false},
		{name: "same team away",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"A"}, Away},
			want:   1 - winBy2sigma12, want1: -2., wantErr: false},
		{name: "rating by -1 neutral",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"B"}, Neutral},
			want:   1 - winBy1sigma12, want1: -1., wantErr: false},
		{name: "rating by 1 neutral",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"B"}, Team{"A"}, Neutral},
			want:   winBy1sigma12, want1: 1., wantErr: false},
		{name: "rating by -1 neutral standard",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"B"}, Neutral},
			want:   1 - winBy1sigma1, want1: -1., wantErr: false},
		{name: "rating by 1 neutral standard",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"B"}, Team{"A"}, Neutral},
			want:   winBy1sigma1, want1: 1., wantErr: false},
		{name: "missing team 1",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"Q"}, Team{"B"}, Neutral},
			want:   0, want1: 0, wantErr: true},
		{name: "missing team 2",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{"Z"}, Neutral},
			want:   0, want1: 0, wantErr: true},
		{name: "bye team 1",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{}, Team{"B"}, Neutral},
			want:   0, want1: 0, wantErr: false},
		{name: "bye team 2",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{"A"}, Team{}, Neutral},
			want:   1, want1: 0, wantErr: false},
		{name: "bye both teams",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, ratings: rm},
			args:   args{Team{}, Team{}, Neutral},
			want:   math.NaN(), want1: math.NaN(), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GaussianSpreadModel{
				dist:      tt.fields.dist,
				homeBias:  tt.fields.homeBias,
				closeBias: tt.fields.closeBias,
				ratings:   tt.fields.ratings,
			}
			got, got1, err := m.Predict(tt.args.t1, tt.args.t2, tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GaussianSpreadModel.Predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !isDarnClose(got, tt.want) {
				t.Errorf("GaussianSpreadModel.Predict() got = %v, want %v", got, tt.want)
			}
			if !isDarnClose(got1, tt.want1) {
				t.Errorf("GaussianSpreadModel.Predict() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLookupModel_Predict(t *testing.T) {
	mm := make(matchupMap)
	mm[teamPair{Team{"A"}, Team{"B"}}] = -1
	mm[teamPair{Team{"A"}, Team{"C"}}] = 2
	mm[teamPair{Team{"B"}, Team{"C"}}] = 0

	winBy1sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(1)
	winBy2sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(2)

	winBy1sigma1 := prob.Normal{Mu: 0, Sigma: 1}.Cdf(1)

	type fields struct {
		dist      prob.Normal
		homeBias  float64
		closeBias float64
		spreads   matchupMap
	}
	type args struct {
		t1  Team
		t2  Team
		loc RelativeLocation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{name: "0 neutral",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"B"}, Team{"C"}, Neutral},
			want:   .5, want1: 0, wantErr: false},
		{name: "0 near",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"B"}, Team{"C"}, Near},
			want:   winBy1sigma12, want1: 1, wantErr: false},
		{name: "0 home",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"B"}, Team{"C"}, Home},
			want:   winBy2sigma12, want1: 2, wantErr: false},
		{name: "0 far",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"B"}, Team{"C"}, Far},
			want:   1 - winBy1sigma12, want1: -1, wantErr: false},
		{name: "0 away",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"B"}, Team{"C"}, Away},
			want:   1 - winBy2sigma12, want1: -2, wantErr: false},
		{name: "1 neutral",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"A"}, Team{"B"}, Neutral},
			want:   1 - winBy1sigma12, want1: -1, wantErr: false},
		{name: "2 neutral",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"A"}, Team{"C"}, Neutral},
			want:   winBy2sigma12, want1: 2, wantErr: false},
		{name: "1 standard",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"A"}, Team{"B"}, Neutral},
			want:   1 - winBy1sigma1, want1: -1, wantErr: false},
		{name: "1 neutral swap",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"B"}, Team{"A"}, Neutral},
			want:   winBy1sigma12, want1: 1, wantErr: false},
		{name: "0 near swap",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"C"}, Team{"B"}, Near},
			want:   1 - winBy1sigma12, want1: -1, wantErr: false},
		{name: "missing game",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"A"}, Team{"Z"}, Neutral},
			want:   0, want1: 0, wantErr: true},
		{name: "bye team 1",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{}, Team{"B"}, Neutral},
			want:   0, want1: 0, wantErr: false},
		{name: "bye team 2",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{"A"}, Team{}, Neutral},
			want:   1, want1: 0, wantErr: false},
		{name: "bye both teams",
			fields: fields{dist: prob.Normal{Mu: 0, Sigma: 1}, homeBias: 2, closeBias: 1, spreads: mm},
			args:   args{Team{}, Team{}, Neutral},
			want:   math.NaN(), want1: math.NaN(), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LookupModel{
				dist:      tt.fields.dist,
				homeBias:  tt.fields.homeBias,
				closeBias: tt.fields.closeBias,
				spreads:   tt.fields.spreads,
			}
			got, got1, err := m.Predict(tt.args.t1, tt.args.t2, tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupModel.Predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !isDarnClose(got, tt.want) {
				t.Errorf("LookupModel.Predict() got = %v, want %v", got, tt.want)
			}
			if !isDarnClose(got1, tt.want1) {
				t.Errorf("LookupModel.Predict() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewLookupModel(t *testing.T) {

	mm := make(matchupMap)
	mm[teamPair{Team{"A"}, Team{"B"}}] = -1
	mm[teamPair{Team{"A"}, Team{"C"}}] = 2
	mm[teamPair{Team{"B"}, Team{"C"}}] = 0
	want := &LookupModel{spreads: mm, dist: prob.Normal{Mu: 0, Sigma: 12}, homeBias: 2, closeBias: 1}

	type args struct {
		homeTeams []Team
		roadTeams []Team
		spreads   []float64
		stdDev    float64
		homeBias  float64
		closeBias float64
	}
	tests := []struct {
		name string
		args args
		want *LookupModel
	}{
		{"basic", args{
			homeTeams: []Team{Team{"A"}, Team{"A"}, Team{"B"}},
			roadTeams: []Team{Team{"B"}, Team{"C"}, Team{"C"}},
			spreads:   []float64{-1, 2, 0},
			stdDev:    12, homeBias: 2, closeBias: 1},
			want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLookupModel(tt.args.homeTeams, tt.args.roadTeams, tt.args.spreads, tt.args.stdDev, tt.args.homeBias, tt.args.closeBias); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLookupModel() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test bad arguments
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewLookupModel([]Team{Team{"A"}}, []Team{Team{"A"}, Team{"B"}}, []float64{1}, 12, 2, 1)
}

func TestNewGaussianSpreadModel(t *testing.T) {
	rm := make(map[Team]float64)
	rm[Team{"A"}] = 0.
	rm[Team{"B"}] = 1.

	want := &GaussianSpreadModel{prob.Normal{Mu: 0, Sigma: 12}, 2, 1, rm}

	type args struct {
		ratings   map[Team]float64
		stdDev    float64
		homeBias  float64
		closeBias float64
	}
	tests := []struct {
		name string
		args args
		want *GaussianSpreadModel
	}{
		{"basic", args{rm, 12, 2, 1}, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGaussianSpreadModel(tt.args.ratings, tt.args.stdDev, tt.args.homeBias, tt.args.closeBias); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGaussianSpreadModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

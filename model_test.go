package pickem

import (
	"math"
	"testing"

	"github.com/atgjack/prob"
)

func isClose(a, b, epsilon float64) bool {
	i, r := math.Modf(a - b)
	return i == 0 && math.Abs(r) < math.Abs(epsilon)
}

func TestGaussianSpreadModelPredict(t *testing.T) {

	rm := make(map[Team]float64)
	rm[Team{"A"}] = 0.
	rm[Team{"B"}] = 1.
	rm[BYE] = -100.

	gsm := NewGaussianSpreadModel(rm, 12., 2., 1.)

	tests := []struct {
		team1  Team
		team2  Team
		loc    RelativeLocation
		prob   float64
		spread float64
		fail   bool
	}{
		{Team{"A"}, Team{"B"}, Neutral, .466793, -1., false},
		{Team{"A"}, Team{"B"}, Home, .533207, 1., false},
		{Team{"A"}, Team{"B"}, Away, .401294, -3., false},
		{Team{"A"}, Team{"B"}, Near, .5, 0., false},
		{Team{"A"}, Team{"B"}, Far, .433816, -2., false},

		{Team{"B"}, Team{"A"}, Far, .5, 0., false},

		{Team{}, Team{"B"}, Neutral, 0., 0., false},
		{Team{"A"}, Team{}, Neutral, 1., 0., false},

		{Team{"A"}, Team{"C"}, Neutral, 0., 0., true},
		{Team{"Q"}, Team{"B"}, Neutral, 0., 0., true},
	}

	for _, test := range tests {
		prob, spread, err := gsm.Predict(test.team1, test.team2, test.loc)
		if !test.fail && err != nil {
			t.Errorf("test failed: %v", err)
			continue
		}
		if test.fail && err == nil {
			t.Error("expected failure, got none")
			continue
		}
		if !isClose(prob, test.prob, 1e-6) {
			t.Errorf("expected probability %f, got %f", test.prob, prob)
		}
		if !isClose(spread, test.spread, 1e-6) {
			t.Errorf("expected spread %f, got %f", test.spread, spread)
		}
	}
}

func TestGaussianSpreadModel_Predict(t *testing.T) {
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
		// TODO: Add test cases.
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
			if got != tt.want {
				t.Errorf("GaussianSpreadModel.Predict() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GaussianSpreadModel.Predict() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLookupModel_Predict(t *testing.T) {
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
		// TODO: Add test cases.
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
			if got != tt.want {
				t.Errorf("LookupModel.Predict() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LookupModel.Predict() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

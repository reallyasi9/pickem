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
	rm[BYE] = -100.

	winBy1sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(1)
	winBy2sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(2)
	// winBy3sigma12 := prob.Normal{Mu: 0, Sigma: 12}.Cdf(3)

	winBy1sigma1 := prob.Normal{Mu: 0, Sigma: 1}.Cdf(1)
	// winBy2sigma1 := prob.Normal{Mu: 0, Sigma: 1}.Cdf(2)
	// winBy3sigma1 := prob.Normal{Mu: 0, Sigma: 1}.Cdf(3)

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

package calculator_test

import (
	"github.com/jattento/passenger-tracker/internal/calculator"
	"math/rand"
	"testing"
)

func TestCalculator_FirstAndLast(t *testing.T) {
	tests := []struct {
		name             string
		flights          [][]calculator.AirportShortName
		wantFirstAirport string
		wantLastAirport  string
		wantErr          error
	}{
		{
			name:             "simple flight route",
			flights:          [][]calculator.AirportShortName{{"EZE", "BRA"}, {"BRA", "MIA"}},
			wantFirstAirport: "EZE",
			wantLastAirport:  "MIA",
			wantErr:          nil,
		},
		{
			name: "complex flights route with start and end airport repeating itself in the middle",
			flights: [][]calculator.AirportShortName{{"EZE", "BRA"}, {"BRA", "NYC"}, {"NYC", "BER"}, {"BER", "EZE"},
				{"EZE", "MIA"}, {"MIA", "EZE"}, {"EZE", "COL"}, {"COL", "NIC"}, {"NIC", "MIA"}},
			wantFirstAirport: "EZE",
			wantLastAirport:  "MIA",
			wantErr:          nil,
		},
		{
			name:             "single same airport route (doesnt make much sense but still technically possible)",
			flights:          [][]calculator.AirportShortName{{"EZE", "EZE"}},
			wantFirstAirport: "EZE",
			wantLastAirport:  "EZE",
			wantErr:          nil,
		},
		{
			name: "ERR - start and final airport are the same",
			flights: [][]calculator.AirportShortName{{"EZE", "BRA"}, {"BRA", "NYC"}, {"NYC", "BER"}, {"BER", "EZE"},
				{"EZE", "MIA"}, {"MIA", "EZE"}, {"EZE", "COL"}, {"COL", "NIC"}, {"NIC", "MIA"}, {"MIA", "EZE"}},
			wantErr: calculator.ErrNondeterministic,
		},
		{
			name:    "ERR - no data",
			flights: [][]calculator.AirportShortName{},
			wantErr: calculator.ErrMissingData,
		},
		{
			name:    "ERR - wrong formatted data",
			flights: [][]calculator.AirportShortName{{"EZE", "BRA"}, {"MIA"}},
			wantErr: calculator.ErrMissingData,
		},
		{
			name:    "ERR - flights not connected",
			flights: [][]calculator.AirportShortName{{"EZE", "BRA"}, {"NYC", "MIA"}},
			wantErr: calculator.ErrFlightsNotConnected,
		},
		{
			name:    "ERR - flights not connected 2",
			flights: [][]calculator.AirportShortName{{"EZE", "BRA"}, {"BRA", "MIA"}, {"BAL", "MIA"}},
			wantErr: calculator.ErrFlightsNotConnected,
		},
	}
	for _, tt := range tests {
		// This ensures all the test cases behave the same every time we run the tests
		deterministicRandSource := rand.New(rand.NewSource(123151232312))

		t.Run(tt.name, func(t *testing.T) {
			deterministicRandSource.Shuffle(len(tt.flights), func(i, j int) {
				tt.flights[i], tt.flights[j] = tt.flights[j], tt.flights[i]
			})

			gotFirstAirport, gotLastAirport, err := (&calculator.Calculator{}).FirstAndLast(tt.flights...)
			if (err != nil) != (tt.wantErr != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("FirstAndLast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFirstAirport != tt.wantFirstAirport {
				t.Errorf("FirstAndLast() gotFirstAirport = %v, want %v", gotFirstAirport, tt.wantFirstAirport)
			}
			if gotLastAirport != tt.wantLastAirport {
				t.Errorf("FirstAndLast() gotLastAirport = %v, want %v", gotLastAirport, tt.wantLastAirport)
			}
		})
	}
}

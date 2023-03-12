package calculator

import "errors"

type AirportShortName = string

// Calculator helper that performs efficient calculations.
type Calculator struct {
	// We are going to attach the functionality of this package to this struct since its probably we are going to add
	// dependencies...
}

var (
	// Possible package errors...

	ErrFlightsNotConnected = errors.New("flights not connected")
	ErrMissingData         = errors.New("missing flights data")
	ErrNondeterministic    = errors.New("start airport and final airport are the same")
)

// FirstAndLast returns the start and end airport of a connected flight scheme.
// The input doesn't need to be sorted.
//
//	Example: [][]route.AirportShortName{{"BRA", "MIA"},{"EZE", "BRA"}} returns "EZE","MIA",nil
func (c *Calculator) FirstAndLast(flights ...[]AirportShortName) (firstAirport, lastAirport AirportShortName, err error) {
	// Check if there is anything to calculate.
	if len(flights) == 0 {
		return "", "", ErrMissingData
	}

	count := make(map[AirportShortName]int)

	// This maps acts like a helper to keep track of how many times each airport got visited.
	unbalancedAmountOfVisitedAirports := make(map[AirportShortName]struct{})

	var first, last AirportShortName

	// Count the number of appearances of each airport as origin and destination.
	for _, flight := range flights {
		// Check data is formatted correctly.
		if len(flight) != 2 {
			return "", "", ErrMissingData
		}

		count[flight[0]]++
		count[flight[1]]--

		for _, airportName := range []AirportShortName{flight[0], flight[1]} {
			// This statement is going to keep track which airports have an unbalanced amount of in and out flights.
			// At the end, we should only have two of these: the starting and ending airports
			if count[airportName] != 0 {
				unbalancedAmountOfVisitedAirports[airportName] = struct{}{}
			} else {
				delete(unbalancedAmountOfVisitedAirports, airportName)
			}

			if count[airportName] == 1 {
				// If we exited the airport one more than we entered it then it's the first one.
				first = airportName
			}
			if count[airportName] == -1 {
				// If we entered the airport one more than we exited it then it's the first one.
				last = airportName
			}
		}
	}

	// Check if there is a single one airport flight
	// (doesn't make much sense but its possible, and it's not responsibility of this code to decide if it's allowed)
	if len(count) == 1 {
		for airportName := range count {
			return airportName, airportName, nil
		}
	}

	// Check that all flights are connected.
	if len(unbalancedAmountOfVisitedAirports) != 2 {
		if len(unbalancedAmountOfVisitedAirports) == 0 {
			// Check if it's impossible to determine which are the initial and final airport because they are the same.
			return "", "", ErrNondeterministic
		}

		return "", "", ErrFlightsNotConnected
	}

	return first, last, nil
}

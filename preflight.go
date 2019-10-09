package preflight

import (
	"context"
	"strings"
)

// statuses
const (
	_ = iota
	StatusGreen
	StatusYellow
	StatusRed
	StatusUnknown
)

// Check defines the interface all preflight checks must expose.
type Check interface {
	// Name should return a unique name for this check.
	Name() string
	// Check should execute the check and return a list of Results.
	Check(ctx context.Context) []Result
}

// Result represents the result of a check.
type Result struct {
	// Name is the name or short desription of the check.
	Name string `json:"name"`
	// Message is a status or error message.
	Message string `json:"message"`
	// Status represents the check's state: green/yellow/red.
	Status int `json:"status"`
}

// less is a comparison function that can be passed to sort.Slice. It sorts the
// data by Status (descending), then alphabetically by Name (ascending).
func less(data []Result) func(i, j int) bool {
	return func(i, j int) bool {
		if data[i].Status == data[j].Status {
			return data[i].Name < data[j].Name
		}
		return data[i].Status > data[j].Status
	}
}

// ConvertStatusString converts a string ("red"/"yellow"/"green") into the status code.
func ConvertStatusString(color string) int {
	switch strings.ToLower(color) {
	case "green":
		return StatusGreen
	case "yellow":
		return StatusYellow
	case "red":
		return StatusRed
	default:
		return StatusUnknown
	}
}

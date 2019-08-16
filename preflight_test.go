package preflight

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLess(t *testing.T) {
	results := []Result{
		{Status: StatusGreen, Name: "foo"},
		{Status: StatusGreen, Name: "abc"},
		{Status: StatusRed, Name: "def"},
		{Status: StatusRed, Name: "456"},
		{Status: StatusYellow, Name: "kershaw"},
	}

	sort.Slice(results, less(results))

	require.Len(t, results, 5)
	expectedOrder := []string{
		"456",
		"def",
		"kershaw",
		"abc",
		"foo",
	}

	for i, expected := range expectedOrder {
		require.Equal(t, expected, results[i].Name)
	}
}

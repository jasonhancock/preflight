package preflight

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChecker(t *testing.T) {
	gimmieCheck := func(status int, name, message string, sleep time.Duration) *CheckMock {
		return &CheckMock{
			CheckFunc: func(ctx context.Context) []Result {
				time.Sleep(sleep)
				return []Result{
					{
						Status:  status,
						Name:    name,
						Message: message,
					},
				}
			},
			NameFunc: func() string {
				return name
			},
		}
	}

	t.Run("normal", func(t *testing.T) {
		checker := NewChecker(
			WithCheck(gimmieCheck(StatusGreen, "green", "it's green!!!", 0)),
			WithCheck(gimmieCheck(StatusRed, "red2", "red2", 0)),
			WithCheck(gimmieCheck(StatusRed, "red1", "red1", 0)),
		)

		results := checker.Check(context.Background())

		require.Len(t, results, 3)
		require.Equal(t, "red1", results[0].Name)
		require.Equal(t, "red2", results[1].Name)
		require.Equal(t, "green", results[2].Name)
	})

	t.Run("check timed out", func(t *testing.T) {
		timeout := 100 * time.Millisecond
		checker := NewChecker(
			WithTimeout(timeout),
			WithCheck(gimmieCheck(StatusGreen, "green", "it's green!!!", 0)),
			WithCheck(gimmieCheck(StatusGreen, "slow poke", "taking my sweet time", 2*timeout)),
		)

		results := checker.Check(context.Background())
		require.Len(t, results, 2)
		require.Equal(t, StatusRed, results[0].Status)
		require.Equal(t, "Check timed out.", results[0].Message)
		require.Equal(t, "slow poke", results[0].Name)

		require.Equal(t, StatusGreen, results[1].Status)
		require.Equal(t, "green", results[1].Name)
	})
}

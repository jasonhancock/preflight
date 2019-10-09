package preflight

import (
	"context"
	"sort"
	"time"
)

// Option is a function for customizing the checker.
type Option func(*options)

type options struct {
	timeout time.Duration
	checks  []Check
}

// WithCheck adds a Check to the Checker.
func WithCheck(c Check) Option {
	return func(o *options) {
		o.checks = append(o.checks, c)
	}
}

// WithTimeout sets the Checker's timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// Checker executes multiple checks in parallel and returns the list of results.
type Checker struct {
	checks  []Check
	timeout time.Duration
}

// NewChecker constructs a new Checker.
func NewChecker(opts ...Option) *Checker {
	opt := &options{
		timeout: 10 * time.Second,
	}

	for _, o := range opts {
		o(opt)
	}

	return &Checker{
		timeout: opt.timeout,
		checks:  opt.checks,
	}
}

// Check executes all of the checks and returns a sorted list of results.
func (c *Checker) Check(ctx context.Context) []Result {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, c.timeout)
	defer cancel()

	results := make([]Result, 0, len(c.checks))

	type namedResult struct {
		name   string
		result []Result
	}

	timedout := make(map[string]struct{}, len(c.checks))
	res := make(chan namedResult)
	for _, check := range c.checks {
		timedout[check.Name()] = struct{}{}

		go func(chk Check) {
			r := namedResult{
				name:   chk.Name(),
				result: chk.Check(ctx),
			}

			select {
			case <-ctx.Done():
				return
			case res <- r:
			}
		}(check)
	}

	for i := 0; i < len(c.checks); i++ {
		select {
		case <-ctx.Done():
			for name := range timedout {
				results = append(results, Result{
					Status:  StatusRed,
					Message: "Check timed out.",
					Name:    name,
				})
			}
			break
		case r := <-res:
			results = append(results, r.result...)
			delete(timedout, r.name)
		}
	}

	sort.Slice(results, less(results))

	return results
}

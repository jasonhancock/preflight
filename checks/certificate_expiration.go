package checks

import (
	"context"
	"crypto/x509"
	"time"

	"github.com/jasonhancock/preflight"
)

// CertExpirationCheck is a preflight check that will warn when an x509 cert is going to expire
type CertExpirationCheck struct {
	name       string
	expiration time.Time
	yellow     time.Time
	red        time.Time
}

// NewCertExpirationCheck initializes a CertExpirationCheck. name is used in the
// result. threshYellow is when to start alerting a yellow status. threshRed is
// when to start alerting a red status.
func NewCertExpirationCheck(cert *x509.Certificate, name string, threshYellow, threshRed time.Duration) *CertExpirationCheck {
	return &CertExpirationCheck{
		name:       name,
		expiration: cert.NotAfter,
		yellow:     cert.NotAfter.Add(threshYellow),
		red:        cert.NotAfter.Add(threshRed),
	}
}

// Name returns the name of this check.
func (c *CertExpirationCheck) Name() string {
	return c.name + " certificate expiration check"
}

// Check runs the expiration check.
func (c *CertExpirationCheck) Check(ctx context.Context) preflight.Result {
	status := preflight.StatusGreen
	if time.Now().After(c.yellow) {
		status = preflight.StatusYellow
	}
	if time.Now().After(c.red) {
		status = preflight.StatusRed
	}

	return preflight.Result{
		Name:    c.name + " - Cert expiration",
		Message: "cert expires at " + c.expiration.String(),
		Status:  status,
	}
}

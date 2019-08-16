package checks

import (
	"context"
	"fmt"

	"github.com/jasonhancock/preflight"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// AWSSESVerificationStatusCheck is a preflight.Check that looks at whether or
// not a given email address has been verified in the SES console.
type AWSSESVerificationStatusCheck struct {
	svc     sesiface.SESAPI
	address string
}

//NewAWSSESVerificationStatusCheck initializes a new AWSSESVerificationStatusCheck.
func NewAWSSESVerificationStatusCheck(svc sesiface.SESAPI, emailAddress string) *AWSSESVerificationStatusCheck {
	return &AWSSESVerificationStatusCheck{
		svc:     svc,
		address: emailAddress,
	}
}

var _ preflight.Check = (*AWSSESVerificationStatusCheck)(nil)

// Name returns the unique name for this check.
func (c *AWSSESVerificationStatusCheck) Name() string {
	return fmt.Sprintf("AWS SES Verification Status - %s", c.address)
}

// Check executes the check.
func (c *AWSSESVerificationStatusCheck) Check(ctx context.Context) preflight.Result {
	in := &ses.GetIdentityVerificationAttributesInput{
		Identities: []*string{&c.address},
	}

	result := preflight.Result{
		Name: c.Name(),
	}

	out, err := c.svc.GetIdentityVerificationAttributesWithContext(ctx, in)
	if err != nil {
		result.Message = err.Error()
		result.Status = preflight.StatusRed
		return result
	}

	attrs, ok := out.VerificationAttributes[c.address]
	if !ok {
		result.Message = "Address does not exist in AWS SES console for this region."
		result.Status = preflight.StatusRed
	} else if *attrs.VerificationStatus != "Success" {
		result.Message = "Address has not been verified. Current state: " + *attrs.VerificationStatus
		result.Status = preflight.StatusRed
	} else {
		result.Status = preflight.StatusGreen
		result.Message = "Address has been verified."
	}

	return result
}

// AWSSESAccountSendingEnabledCheck is a preflight.Check that looks at whether or
// not an account is still in sandbox mode.
type AWSSESAccountSendingEnabledCheck struct {
	svc sesiface.SESAPI
}

//NewAWSSESAccountSendingEnabledCheck initializes a new AWSSESAccountSendingStatusCheck.
func NewAWSSESAccountSendingEnabledCheck(svc sesiface.SESAPI) *AWSSESAccountSendingEnabledCheck {
	return &AWSSESAccountSendingEnabledCheck{
		svc: svc,
	}
}

var _ preflight.Check = (*AWSSESAccountSendingEnabledCheck)(nil)

// Name returns the unique name for this check.
func (c *AWSSESAccountSendingEnabledCheck) Name() string {
	return "AWS SES Account Sending Status"
}

// Check executes the check.
func (c *AWSSESAccountSendingEnabledCheck) Check(ctx context.Context) preflight.Result {
	in := &ses.GetAccountSendingEnabledInput{}

	result := preflight.Result{
		Name:    c.Name(),
		Message: "Sending has been enabled.",
		Status:  preflight.StatusGreen,
	}

	out, err := c.svc.GetAccountSendingEnabledWithContext(ctx, in)
	if err != nil {
		result.Message = err.Error()
		result.Status = preflight.StatusRed
		return result
	}

	if !*out.Enabled {
		result.Message = "Sending not enabled. Please request to move out of the sandbox."
		result.Status = preflight.StatusRed
	}

	return result
}

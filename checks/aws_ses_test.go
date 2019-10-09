package checks

import (
	"context"
	"testing"

	"github.com/jasonhancock/preflight"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/stretchr/testify/require"
)

func TestAWSSESVerificationStatusCheck(t *testing.T) {
	const addr = "bob@example.com"
	var tests = []struct {
		description string
		out         *ses.GetIdentityVerificationAttributesOutput
		expected    preflight.Result
	}{
		{
			"normal",
			&ses.GetIdentityVerificationAttributesOutput{
				VerificationAttributes: map[string]*ses.IdentityVerificationAttributes{
					addr: &ses.IdentityVerificationAttributes{
						VerificationStatus: aws.String("Success"),
					},
				},
			},
			preflight.Result{
				Name:    "AWS SES Verification Status - bob@example.com",
				Message: "Address has been verified.",
				Status:  preflight.StatusGreen,
			},
		},
		{
			"not verified",
			&ses.GetIdentityVerificationAttributesOutput{
				VerificationAttributes: map[string]*ses.IdentityVerificationAttributes{
					addr: &ses.IdentityVerificationAttributes{
						VerificationStatus: aws.String("Pending"),
					},
				},
			},
			preflight.Result{
				Name:    "AWS SES Verification Status - bob@example.com",
				Message: "Address has not been verified. Current state: Pending",
				Status:  preflight.StatusRed,
			},
		},
		{
			"not found",
			&ses.GetIdentityVerificationAttributesOutput{
				VerificationAttributes: map[string]*ses.IdentityVerificationAttributes{},
			},
			preflight.Result{
				Name:    "AWS SES Verification Status - bob@example.com",
				Message: "Address does not exist in AWS SES console for this region.",
				Status:  preflight.StatusRed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			svc := &mockSES{verificationOut: tt.out}
			check := NewAWSSESVerificationStatusCheck(svc, addr)
			result := check.Check(context.Background())

			require.Equal(t, tt.expected.Name, result[0].Name)
			require.Equal(t, tt.expected.Status, result[0].Status)
			require.Equal(t, tt.expected.Message, result[0].Message)
		})
	}
}

func TestAWSSESAccountSendingEnabledCheck(t *testing.T) {
	var tests = []struct {
		description string
		value       bool
		expected    preflight.Result
	}{
		{
			"enabled",
			true,
			preflight.Result{
				Name:    "AWS SES Account Sending Status",
				Message: "Sending has been enabled.",
				Status:  preflight.StatusGreen,
			},
		},
		{
			"not enabled",
			false,
			preflight.Result{
				Name:    "AWS SES Account Sending Status",
				Message: "Sending not enabled. Please request to move out of the sandbox.",
				Status:  preflight.StatusRed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			svc := &mockSES{sendingOut: &ses.GetAccountSendingEnabledOutput{Enabled: aws.Bool(tt.value)}}
			check := NewAWSSESAccountSendingEnabledCheck(svc)
			result := check.Check(context.Background())

			require.Equal(t, tt.expected.Name, result[0].Name)
			require.Equal(t, tt.expected.Status, result[0].Status)
			require.Equal(t, tt.expected.Message, result[0].Message)
		})
	}
}

type mockSES struct {
	verificationOut *ses.GetIdentityVerificationAttributesOutput
	sendingOut      *ses.GetAccountSendingEnabledOutput
	sesiface.SESAPI
}

func (m *mockSES) GetIdentityVerificationAttributesWithContext(aws.Context, *ses.GetIdentityVerificationAttributesInput, ...request.Option) (*ses.GetIdentityVerificationAttributesOutput, error) {
	return m.verificationOut, nil
}

func (m *mockSES) GetAccountSendingEnabledWithContext(aws.Context, *ses.GetAccountSendingEnabledInput, ...request.Option) (*ses.GetAccountSendingEnabledOutput, error) {
	return m.sendingOut, nil
}

package bootstrap

import (
	"strings"
	"testing"
)

func TestParseEnvironment(t *testing.T) {
	for _, environment := range []string{
		EnvironmentDev,
		EnvironmentStaging,
		EnvironmentProduction,
	} {
		t.Run(environment, func(t *testing.T) {
			got, err := parseEnvironment(environment)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != environment {
				t.Fatalf("expected %q, got %q", environment, got)
			}
		})
	}
}

func TestParseEnvironmentRejectsUnknownValue(t *testing.T) {
	_, err := parseEnvironment("prod")
	if err == nil {
		t.Fatal("expected unknown environment to be rejected")
	}
}

func TestValidateDebugTokenExposure(t *testing.T) {
	if err := validateDebugTokenExposure(EnvironmentStaging, true); err != nil {
		t.Fatalf("expected staging debug token exposure to be allowed, got %v", err)
	}

	err := validateDebugTokenExposure(EnvironmentProduction, true)
	if err == nil {
		t.Fatal("expected production debug token exposure to be rejected")
	}
	if !strings.Contains(err.Error(), "must be false in production") {
		t.Fatalf("unexpected error: %v", err)
	}
}

package application

import (
	"net/http"
	"testing"

	"github.com/seu-usuario/bank-api/internal/security/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func TestRegisterErrors_TransactionPasswordMappings(t *testing.T) {
	RegisterErrors()

	tests := []struct {
		name       string
		err        error
		wantCode   string
		wantStatus int
	}{
		{
			name:       "already set",
			err:        domain.ErrTransactionPasswordAlreadySet,
			wantCode:   sharederrors.ErrCodeTransactionPasswordAlreadySet,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "not set",
			err:        domain.ErrTransactionPasswordNotSet,
			wantCode:   sharederrors.ErrCodeTransactionPasswordNotSet,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "invalid",
			err:        domain.ErrTransactionPasswordInvalid,
			wantCode:   sharederrors.ErrCodeTransactionPasswordInvalid,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "locked",
			err:        domain.ErrTransactionPasswordLocked,
			wantCode:   sharederrors.ErrCodeTransactionPasswordLocked,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "step-up endpoint not allowed",
			err:        domain.ErrStepUpEndpointNotAllowed,
			wantCode:   sharederrors.ErrCodeStepUpEndpointNotAllowed,
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sharederrors.MapError(tt.err)

			if got.Code != tt.wantCode {
				t.Fatalf("expected code %q, got %q", tt.wantCode, got.Code)
			}
			if got.Status != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, got.Status)
			}
		})
	}
}

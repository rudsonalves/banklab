package delivery

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/installation/application"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
	sharedheaders "github.com/seu-usuario/bank-api/internal/shared/http/headers"
)

type registerInstallationUseCase interface {
	Execute(ctx context.Context, input application.RegisterInstallationInput) (*application.RegisterInstallationOutput, error)
}

type listInstallationsUseCase interface {
	Execute(ctx context.Context) (*application.ListInstallationsOutput, error)
}

type revokeInstallationUseCase interface {
	Execute(ctx context.Context, input application.RevokeInstallationInput) (*application.RevokeInstallationOutput, error)
}

type Handler struct {
	registerInstallation registerInstallationUseCase
	listInstallations    listInstallationsUseCase
	revokeInstallation   revokeInstallationUseCase
}

type registerInstallationData struct {
	AccessToken            string    `json:"access_token"`
	RefreshToken           string    `json:"refresh_token"`
	InstallationResourceID uuid.UUID `json:"installation_resource_id"`
	InstallationStatus     string    `json:"installation_status"`
}

type installationData struct {
	ResourceID  uuid.UUID  `json:"resource_id"`
	Status      string     `json:"status"`
	FirstSeenAt time.Time  `json:"first_seen_at"`
	LastSeenAt  time.Time  `json:"last_seen_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type listInstallationsData struct {
	Installations []installationData `json:"installations"`
}

type revokeInstallationData struct {
	ResourceID uuid.UUID  `json:"resource_id"`
	Status     string     `json:"status"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

func New(
	registerInstallation registerInstallationUseCase,
	listInstallations listInstallationsUseCase,
	revokeInstallation revokeInstallationUseCase,
) *Handler {
	return &Handler{
		registerInstallation: registerInstallation,
		listInstallations:    listInstallations,
		revokeInstallation:   revokeInstallation,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if h.registerInstallation == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	installationID, err := parseCanonicalUUID(r.Header.Get(sharedheaders.InstallationID))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidInstallationID))
		return
	}

	output, err := h.registerInstallation.Execute(r.Context(), application.RegisterInstallationInput{
		PresentedInstallationID: installationID,
		StepUpToken:             r.Header.Get(sharedheaders.StepUpToken),
	})
	if err != nil {
		log.Printf("event=register_installation error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}
	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, registerInstallationData{
		AccessToken:            output.AccessToken,
		RefreshToken:           output.RefreshToken,
		InstallationResourceID: output.InstallationResourceID,
		InstallationStatus:     output.InstallationStatus,
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if h.listInstallations == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	output, err := h.listInstallations.Execute(r.Context())
	if err != nil {
		log.Printf("event=list_installations error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}
	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	items := make([]installationData, 0, len(output.Installations))
	for _, installation := range output.Installations {
		items = append(items, installationData{
			ResourceID:  installation.ResourceID,
			Status:      installation.Status,
			FirstSeenAt: installation.FirstSeenAt,
			LastSeenAt:  installation.LastSeenAt,
			RevokedAt:   installation.RevokedAt,
			CreatedAt:   installation.CreatedAt,
			UpdatedAt:   installation.UpdatedAt,
		})
	}

	sharedhttp.WriteJSON(w, http.StatusOK, listInstallationsData{Installations: items})
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	if h.revokeInstallation == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	resourceID, err := parseCanonicalUUID(r.PathValue("installation_resource_id"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	output, err := h.revokeInstallation.Execute(r.Context(), application.RevokeInstallationInput{
		ResourceID: resourceID,
	})
	if err != nil {
		log.Printf("event=revoke_installation error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}
	if output == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, revokeInstallationData{
		ResourceID: output.ResourceID,
		Status:     output.Status,
		RevokedAt:  output.RevokedAt,
	})
}

func parseCanonicalUUID(raw string) (uuid.UUID, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return uuid.Nil, sharederrors.ErrInvalidRequest
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, err
	}
	if parsed.Version() != 4 || parsed.String() != value {
		return uuid.Nil, sharederrors.ErrInvalidRequest
	}

	return parsed, nil
}

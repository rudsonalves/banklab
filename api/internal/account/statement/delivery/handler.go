package delivery

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	statementapp "github.com/seu-usuario/bank-api/internal/account/statement/application"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type statementUseCase interface {
	Execute(ctx context.Context, input statementapp.GetStatementInput) (*statementapp.Statement, error)
}

type Handler struct {
	statement statementUseCase
}

func New(statement statementUseCase) *Handler {
	return &Handler{statement: statement}
}

// Statement handles the HTTP request for retrieving an account statement.
func (h *Handler) Statement(w http.ResponseWriter, r *http.Request) {
	if h.statement == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := requireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	limit, err := parseOptionalInt(r.URL.Query().Get("limit"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	cursor, err := parseOptionalTime(r.URL.Query().Get("cursor"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	cursorID, err := parseOptionalUUID(r.URL.Query().Get("cursor_id"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	if (cursor == nil) != (cursorID == nil) {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	from, err := parseOptionalTime(r.URL.Query().Get("from"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	to, err := parseOptionalTime(r.URL.Query().Get("to"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	result, err := h.statement.Execute(r.Context(), statementapp.GetStatementInput{
		User:      user,
		AccountID: accountID,
		Limit:     limit,
		Cursor:    cursor,
		CursorID:  cursorID,
		From:      from,
		To:        to,
	})
	if err != nil {
		log.Printf("event=get_statement error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	items := make([]StatementItemData, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, StatementItemData{
			TransactionID: item.TransactionID,
			Type:          item.Type,
			Amount:        item.Amount,
			BalanceAfter:  item.BalanceAfter,
			ReferenceID:   item.ReferenceID,
			Description:   item.Description,
			CreatedAt:     item.CreatedAt,
		})
	}

	var nextCursor *StatementCursorData
	if result.NextCursor != nil {
		nextCursor = &StatementCursorData{
			CreatedAt: result.NextCursor.CreatedAt,
			ID:        result.NextCursor.ID,
		}
	}

	sharedhttp.WriteJSON(w, http.StatusOK, StatementData{
		AccountID:  result.AccountID,
		Items:      items,
		NextCursor: nextCursor,
	})
}

func requireUser(ctx context.Context) (*authdomain.AuthenticatedUser, error) {
	user, ok := sharedauthctx.GetAuthenticatedUser(ctx)
	if !ok || user == nil {
		return nil, authdomain.ErrUnauthorized
	}

	return user, nil
}

// parseOptionalInt parses an optional integer query parameter. If the
// input string is empty, it returns 0 and no error.
func parseOptionalInt(raw string) (int, error) {
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// parseOptionalTime parses an optional time query parameter in RFC3339 format.
func parseOptionalTime(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

// parseOptionalUUID parses an optional UUID query parameter.
func parseOptionalUUID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

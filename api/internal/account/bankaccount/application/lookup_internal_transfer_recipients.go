package application

import (
	"context"
	"strings"
	"unicode"

	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type LookupInternalTransferRecipientsInput struct {
	User          *authdomain.AuthenticatedUser
	Branch        string
	AccountNumber string
	Document      string
}

type lookupInternalTransferRecipientsRepository interface {
	FindTransferRecipientsByBranchAndNumber(ctx context.Context, branch, number string) ([]domain.TransferRecipient, error)
	FindTransferRecipientsByDocument(ctx context.Context, document string) ([]domain.TransferRecipient, error)
}

type LookupInternalTransferRecipients struct {
	repo lookupInternalTransferRecipientsRepository
}

func NewLookupInternalTransferRecipients(repo lookupInternalTransferRecipientsRepository) *LookupInternalTransferRecipients {
	return &LookupInternalTransferRecipients{repo: repo}
}

func (uc *LookupInternalTransferRecipients) Execute(
	ctx context.Context,
	input LookupInternalTransferRecipientsInput,
) ([]domain.TransferRecipient, error) {
	if !CanListOwnAccounts(input.User) {
		return nil, domain.ErrForbidden
	}

	branch := normalizeBranch(input.Branch)
	accountNumber := normalizeAccountNumber(input.AccountNumber)
	document := normalizeDocument(input.Document)

	hasAccountLookup := branch != "" || accountNumber != ""
	hasDocumentLookup := document != ""

	if hasAccountLookup && hasDocumentLookup {
		return nil, domain.ErrInvalidData
	}

	if hasAccountLookup {
		if branch == "" || accountNumber == "" {
			return nil, domain.ErrInvalidData
		}
		return uc.repo.FindTransferRecipientsByBranchAndNumber(ctx, branch, accountNumber)
	}

	if hasDocumentLookup {
		if len(document) != 11 {
			return nil, domain.ErrInvalidData
		}
		return uc.repo.FindTransferRecipientsByDocument(ctx, document)
	}

	return nil, domain.ErrInvalidData
}

func normalizeBranch(value string) string {
	return strings.TrimSpace(value)
}

func normalizeAccountNumber(value string) string {
	return onlyDigits(value)
}

func normalizeDocument(value string) string {
	return domain.NormalizeDocument(value)
}

func onlyDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

package delivery

import (
	"net/http"

	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	sharedhttp.WriteSuccess(w, status, data)
}

func writeError(w http.ResponseWriter, status int, err *sharederrors.AppError) {
	sharedhttp.WriteError(w, status, err)
}

package delivery

import (
	"os"
	"testing"

	"github.com/seu-usuario/bank-api/internal/bootstrap"
)

func TestMain(m *testing.M) {
	bootstrap.RegisterErrors()
	os.Exit(m.Run())
}

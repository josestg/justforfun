package health

import (
	"context"

	"github.com/josestg/justforfun/internal/domain/sys"

	dHealth "github.com/josestg/justforfun/internal/domain/health"
)

// UseCase implements the health check use case interface.
type UseCase struct {
	// todo: injects some dependencies here...
}

// implementation checks.
var _ dHealth.UseCase = &UseCase{}

// NewUseCase creates a new health check use case.
func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) HealthReport(_ context.Context) (*dHealth.Report, error) {

	report := dHealth.Report{
		Info:    sys.NewInfo("example-id", "develop"),
		Support: sys.Support{},
	}

	return &report, nil
}

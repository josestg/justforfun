package health

import (
	"context"

	"github.com/josestg/justforfun/internal/domain/sys"
)

// UseCase is contract that must be implemented by the health check use case.
type UseCase interface {
	// HealthReport returns the health report.
	HealthReport(ctx context.Context) (*Report, error)
}

// Report represents the health report.
type Report struct {
	Info    sys.Info    `json:"info"`
	Support sys.Support `json:"support"`
}

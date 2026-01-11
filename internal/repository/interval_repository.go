package repository

import "github.com/mariowenwz/flo-energy-nem12-etl/internal/domain"

// IntervalRepository defines persistence behavior for interval records.
type IntervalRepository interface {

    // SaveBatch persists a batch of interval records atomically.
    // Implementation may choose transaction / bulk insert / buffering.
    SaveBatch(records []domain.IntervalRecord) error
}

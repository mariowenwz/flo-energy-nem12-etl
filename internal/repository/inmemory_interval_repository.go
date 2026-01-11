package repository

import "github.com/mariowenwz/flo-energy-nem12-etl/internal/domain"

// IntervalRepository abstracts persistence of meter readings.
// The current implementation is in-memory for testing purposes.
type InMemoryIntervalRepository struct {
    Records []domain.IntervalRecord
}

func NewInMemoryIntervalRepository() *InMemoryIntervalRepository {
    return &InMemoryIntervalRepository{
        Records: make([]domain.IntervalRecord, 0),
    }
}

func (r *InMemoryIntervalRepository) SaveBatch(
    records []domain.IntervalRecord,
) error {
    r.Records = append(r.Records, records...)
    return nil
}

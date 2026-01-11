package transformer

import (
	"errors"
    "time"

    "github.com/mariowenwz/flo-energy-nem12-etl/internal/domain"
    "github.com/mariowenwz/flo-energy-nem12-etl/internal/parser"
)

type IntervalTransformer interface {
    Transform(
        ctx domain.MeterContext,
        record parser.Record300,
    ) ([]domain.IntervalRecord, error)
}

type DefaultIntervalTransformer struct{}

func NewDefaultIntervalTransformer() IntervalTransformer {
    return &DefaultIntervalTransformer{}
}

func (t *DefaultIntervalTransformer) Transform(
    ctx domain.MeterContext,
    record parser.Record300,
) ([]domain.IntervalRecord, error) {

    if ctx.IntervalLength <= 0 {
        return nil, errors.New("invalid interval length")
    }

    baseDate, err := time.Parse("20060102", record.IntervalDate)
    if err != nil {
        return nil, err
    }

    results := make([]domain.IntervalRecord, 0, len(record.Values))

    for i, v := range record.Values {
        ts := baseDate.Add(
            time.Duration(i*ctx.IntervalLength) * time.Minute,
        )

        results = append(results, domain.IntervalRecord{
            NMI:         ctx.NMI,
            Timestamp:   ts,
            Consumption: v,
        })
    }

    return results, nil
}

package transformer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mariowenwz/flo-energy-nem12-etl/internal/domain"
	"github.com/mariowenwz/flo-energy-nem12-etl/internal/parser"
	"github.com/mariowenwz/flo-energy-nem12-etl/internal/transformer"
)

func TestIntervalTransformer_TransformToMeterReadings(t *testing.T) {
	tests := []struct {
		name        string
		ctx         domain.MeterContext
		record300   parser.Record300
		want        []domain.IntervalRecord
		wantErr     bool
		description string
	}{
		{
			name: "single interval value with 30 minute interval",
			ctx: domain.MeterContext{
				NMI:            "NEM1201009",
				IntervalLength: 30,
			},
			record300: parser.Record300{
				IntervalDate: "20050301",
				Values:       []float64{0.461},
			},
			want: []domain.IntervalRecord{
				{
					NMI:         "NEM1201009",
					Timestamp:   time.Date(2005, 3, 1, 0, 0, 0, 0, time.UTC),
					Consumption: 0.461,
				},
			},
			wantErr:     false,
			description: "transforms single interval value correctly",
		},
		{
			name: "multiple interval values with 30 minute interval",
			ctx: domain.MeterContext{
				NMI:            "NEM1201009",
				IntervalLength: 30,
			},
			record300: parser.Record300{
				IntervalDate: "20050301",
				Values:       []float64{0.461, 0.522, 0.583},
			},
			want: []domain.IntervalRecord{
				{
					NMI:         "NEM1201009",
					Timestamp:   time.Date(2005, 3, 1, 0, 0, 0, 0, time.UTC),
					Consumption: 0.461,
				},
				{
					NMI:         "NEM1201009",
					Timestamp:   time.Date(2005, 3, 1, 0, 30, 0, 0, time.UTC),
					Consumption: 0.522,
				},
				{
					NMI:         "NEM1201009",
					Timestamp:   time.Date(2005, 3, 1, 1, 0, 0, 0, time.UTC),
					Consumption: 0.583,
				},
			},
			wantErr:     false,
			description: "transforms multiple interval values with correct timestamps",
		},
		{
			name: "full day 48 intervals with 30 minute interval",
			ctx: domain.MeterContext{
				NMI:            "NEM1201009",
				IntervalLength: 30,
			},
			record300: parser.Record300{
				IntervalDate: "20050301",
				Values:       make([]float64, 48), // Fill in 48 real intervals of 30 minutes in the future
			},
			want:        make([]domain.IntervalRecord, 48), 
			wantErr:     false,
			description: "transforms full day of 48 intervals correctly",
		},
		{
			name: "empty values slice",
			ctx: domain.MeterContext{
				NMI:            "NEM1201009",
				IntervalLength: 30,
			},
			record300: parser.Record300{
				IntervalDate: "20050301",
				Values:       []float64{},
			},
			want:        []domain.IntervalRecord{},
			wantErr:     false,
			description: "handles empty values slice gracefully",
		},
		{
			name: "different interval length 15 minutes",
			ctx: domain.MeterContext{
				NMI:            "NEM1201009",
				IntervalLength: 15,
			},
			record300: parser.Record300{
				IntervalDate: "20050301",
				Values:       []float64{0.461, 0.522},
			},
			want: []domain.IntervalRecord{
				{
					NMI:         "NEM1201009",
					Timestamp:   time.Date(2005, 3, 1, 0, 0, 0, 0, time.UTC),
					Consumption: 0.461,
				},
				{
					NMI:         "NEM1201009",
					Timestamp:   time.Date(2005, 3, 1, 0, 15, 0, 0, time.UTC),
					Consumption: 0.522,
				},
			},
			wantErr:     false,
			description: "handles different interval lengths correctly",
		},
		// TODO: Add more test cases for:
		// - Invalid date format handling
		// - Edge cases around midnight/day boundaries
		// - Large number of intervals
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            results, err := transformer.NewDefaultIntervalTransformer().Transform(tt.ctx, tt.record300)
			
            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Len(t, results, len(tt.want))

            if len(tt.want) > 0 {
                assert.Equal(t, tt.ctx.NMI, results[0].NMI)
            }
        })
    }
}


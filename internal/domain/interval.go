package domain

import "time"

type IntervalRecord struct {
    NMI         string
    Timestamp   time.Time
    Consumption float64
}

package parser

import (
    "bufio"
    "io"
    "strconv"
    "strings"
    "time"

    "github.com/mariowenwz/flo-energy-nem12-etl/internal/domain"
)


// Record300 represents a raw 300 record parsed from NEM12 input.
// It is an input DTO, not a domain model.
type Record300 struct {
    IntervalDate string
    Values       []float64
}

type NEM12Parser interface {
    Parse(r io.Reader) error
    MeterContext() domain.MeterContext
    Records300() <-chan Record300
}

type nem12Parser struct {
    ctx     domain.MeterContext
    records chan Record300
}

func NewNEM12Parser() NEM12Parser {
    return &nem12Parser{
        records: make(chan Record300, 1000), // Used Buffered Channel to decouple IO-bound parsing from CPU/IO-bound transformation and persistence.
    }
}

func (p *nem12Parser) Parse(r io.Reader) error {
    scanner := bufio.NewScanner(r)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        fields := strings.Split(line, ",")
        if len(fields) == 0 {
            continue
        }

        switch fields[0] {
        case "200":
            p.parse200(fields)
        case "300":
            p.parse300(fields)
        }
    }

    close(p.records)
    return scanner.Err()
}

func (p *nem12Parser) parse200(fields []string) {
    if len(fields) < 9 {
        return
    }

    intervalLen, err := strconv.Atoi(fields[8])
    if err != nil {
        return
    }

    p.ctx = domain.MeterContext{
        NMI:            fields[1],
        IntervalLength: intervalLen,
    }
}

func (p *nem12Parser) parse300(fields []string) {
    if len(fields) < 3 {
        return
    }

    date, err := time.Parse("20060102", fields[1])
    if err != nil {
        return
    }

    values := make([]float64, 0)
    for _, v := range fields[2:] {
        if v == "" {
            continue
        }
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            continue
        }
        values = append(values, f)
    }

    p.records <- Record300{
        IntervalDate: date.Format("20060102"),
        Values:       values,
    }
}


func (p *nem12Parser) MeterContext() domain.MeterContext {
    return p.ctx
}

func (p *nem12Parser) Records300() <-chan Record300 {
    return p.records
}

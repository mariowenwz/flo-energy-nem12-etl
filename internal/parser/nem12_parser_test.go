package parser_test

import (
    "strings"
    "testing"

    //"github.com/mariowenwz/flo-energy-nem12-etl/internal/domain"
    "github.com/mariowenwz/flo-energy-nem12-etl/internal/parser"
)

func TestNEM12Parser_Parse_Single200And300Record(t *testing.T) {
    input := `
100,NEM12,200506081149,UNITEDDP,NEMMCO
200,NEM1201009,E1E2,1,E1,N1,01009,kWh,30,20050610
300,20050301,0,0,0,0,0,0,0,0,0,0,0,0,0.461
`

    r := strings.NewReader(input)

    p := parser.NewNEM12Parser()

    err := p.Parse(r)
    if err != nil {
        t.Fatalf("unexpected parse error: %v", err)
    }

    ctx := p.MeterContext()
    if ctx.NMI != "NEM1201009" {
        t.Errorf("expected NMI %q, got %q", "NEM1201009", ctx.NMI)
    }

    if ctx.IntervalLength != 30 {
        t.Errorf("expected interval length 30, got %d", ctx.IntervalLength)
    }

    records := p.Records300()

    rec, ok := <-records
    if !ok {
        t.Fatal("expected at least one 300 record")
    }

    if rec.IntervalDate != "20050301" {
        t.Errorf("unexpected interval date: %v", rec.IntervalDate)
    }

    if len(rec.Values) == 0 {
        t.Error("expected consumption values, got empty slice")
    }
}

func TestNEM12Parser_Parse_StreamingLargeFile(t *testing.T) {
    // This test verifies that the parser processes records in a streaming manner
    // without loading the entire file into memory.
    t.Skip("implementation pending")
}


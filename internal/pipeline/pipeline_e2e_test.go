package pipeline_test

import (
    "strings"
    "testing"

    "github.com/mariowenwz/flo-energy-nem12-etl/internal/parser"
    "github.com/mariowenwz/flo-energy-nem12-etl/internal/pipeline"
    "github.com/mariowenwz/flo-energy-nem12-etl/internal/repository"
)

func TestPipeline_EndToEnd_SingleNMI_SingleDay(t *testing.T) {
    // given
    input := strings.NewReader(`
200,NEM1201009,E1E2,1,E1,N1,01009,kWh,30,20050610
300,20050301,0,0,0,0,0,0,0,0,0,0,0,0,0.461
`)

    repo := repository.NewInMemoryIntervalRepository()
    p := pipeline.New(
        parser.NewNEM12Parser(),
        repo,
    )

    // when
    err := p.Run(input)

    // then
    if err != nil {
        t.Fatalf("pipeline failed: %v", err)
    }

    if len(repo.Records) != 13 {
        t.Fatalf("expected 13 interval records, got %d", len(repo.Records))
    }
}

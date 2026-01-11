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
200,NEM1202009,NMI1234567,30
300,20240101,1,1,1,1
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

    if len(repo.Records) != 4 {
        t.Fatalf("expected 4 interval records, got %d", len(repo.Records))
    }
}

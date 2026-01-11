package pipeline

import (
    "io"

    "github.com/mariowenwz/flo-energy-nem12-etl/internal/parser"
    "github.com/mariowenwz/flo-energy-nem12-etl/internal/repository"
    "github.com/mariowenwz/flo-energy-nem12-etl/internal/transformer"
)

type Pipeline struct {
    parser      parser.NEM12Parser
    transformer *transformer.DefaultIntervalTransformer
    repo        repository.IntervalRepository
}

func New(
    p parser.NEM12Parser,
    r repository.IntervalRepository,
) *Pipeline {
    return &Pipeline{
        parser:      p,
        transformer: transformer.NewDefaultIntervalTransformer(),
        repo:        r,
    }
}

func (p *Pipeline) Run(input io.Reader) error {
    err := p.parser.Parse(input)
    if err != nil {
        return err
    }

    for record300 := range p.parser.Records300() {
        intervalRecords, err := p.transformer.Transform(p.parser.MeterContext(), record300)
        if err != nil {
            return err
        }

		err = p.repo.SaveBatch(intervalRecords)
		if err != nil {
			return err
		}
    }

    return nil
}



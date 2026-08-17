package traces

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) Fail(err error) {
	s.span.RecordError(err)
	s.span.SetStatus(codes.Error, err.Error())
}

func (s *otelSpan) End() {
	s.span.End()
}

type resolutionSpan struct {
	otelSpan
}

func (s *resolutionSpan) RecordSource(source string) {
	s.span.SetAttributes(attribute.String("ens.source", source))
}

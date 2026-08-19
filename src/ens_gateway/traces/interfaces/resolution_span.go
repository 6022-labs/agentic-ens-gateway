package interfaces

// ResolutionSpan is the root span of a resolution; the answering tier is only known at the end.
type ResolutionSpan interface {
	Span
	RecordSource(source string)
}

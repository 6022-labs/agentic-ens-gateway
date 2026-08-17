package interfaces

// Span is the tracing handle a domain service holds; End must run even on the failure path.
type Span interface {
	Fail(err error)
	End()
}

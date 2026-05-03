package rest

// RestServiceOption represents a runtime option accepted by NewRestService.
// Unknown option types panic during option application.
type RestServiceOption interface{}

type upstreamOption[C, U any] struct {
	mapper RestUpstreamMapper[C, U]
}

type downstreamOption[R any] struct {
	mapper RestDownstreamMapper[R]
}

// WithUpstreamMapper sets a custom upstream mapper.
// Panics if mapper is nil.
func WithUpstreamMapper[C, U any](mapper RestUpstreamMapper[C, U]) RestServiceOption {
	if mapper == nil {
		panic("rest: WithUpstreamMapper: mapper must not be nil")
	}
	return upstreamOption[C, U]{mapper: mapper}
}

// WithDownstreamMapper sets a custom downstream mapper.
// Panics if mapper is nil.
func WithDownstreamMapper[R any](mapper RestDownstreamMapper[R]) RestServiceOption {
	if mapper == nil {
		panic("rest: WithDownstreamMapper: mapper must not be nil")
	}
	return downstreamOption[R]{mapper: mapper}
}

type restServiceOptions[C, R, U any] struct {
	upstream   RestUpstreamMapper[C, U]
	downstream RestDownstreamMapper[R]
}

func newRestServiceOptions[C, R, U any](opts []RestServiceOption) restServiceOptions[C, R, U] {
	o := restServiceOptions[C, R, U]{
		upstream:   &DefaultRestUpstreamMapper[C, U]{},
		downstream: &DefaultRestDownstreamMapper[R]{},
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case upstreamOption[C, U]:
			o.upstream = v.mapper
		case downstreamOption[R]:
			o.downstream = v.mapper
		default:
			panic("rest: unknown or mismatched RestServiceOption type")
		}
	}
	return o
}

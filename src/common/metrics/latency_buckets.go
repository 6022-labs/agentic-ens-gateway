package metrics

// LatencySecondsBuckets are second-scaled histogram bounds (the SDK defaults are ms-scaled, useless for seconds).
var LatencySecondsBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
}

// LongLatencySecondsBuckets extend the tail for seconds-to-minutes ops (LLM calls, inference loops).
var LongLatencySecondsBuckets = []float64{
	0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20, 30, 60, 120,
}

// GenAiDurationSecondsBuckets are the advisory bounds the OTel GenAI metrics conventions attach to their duration
// histograms; the MCP conventions use a different set.
var GenAiDurationSecondsBuckets = []float64{
	0.01, 0.02, 0.04, 0.08, 0.16, 0.32, 0.64, 1.28, 2.56, 5.12, 10.24, 20.48, 40.96, 81.92, 163.84,
}

// DbLatencySecondsBuckets are the DB semconv advisory set. Its 1ms floor is five times finer than the
// HTTP set's; a sub-millisecond cache GET still lands in the first bucket.
var DbLatencySecondsBuckets = []float64{
	0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10,
}

// RpcLatencyMillisecondsBuckets are millisecond-scaled: rpc.client.duration is defined in ms.
var RpcLatencyMillisecondsBuckets = []float64{
	1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 30000, 120000,
}

// ResponseBodySizeBytesBuckets are byte-scaled; the SDK default tops out at 10kB, well under a
// typical agent card or completion payload.
var ResponseBodySizeBytesBuckets = []float64{
	256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304,
}

// LoopIterationsBuckets bound small integer counts; the SDK defaults start at 5 and hide what a bounded loop does.
var LoopIterationsBuckets = []float64{1, 2, 3, 4, 5, 8, 12, 20}

// RetryAttemptsBuckets bound small retry counts; the SDK defaults start at 5 and hide a 3-attempt loop.
var RetryAttemptsBuckets = []float64{1, 2, 3, 4, 5, 8}

module github.com/nickbryan/slogutil/benchmarks

go 1.23.3

replace github.com/nickbryan/slogutil => ../

require (
	github.com/nickbryan/slogutil v0.0.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
)

require github.com/google/go-cmp v0.6.0 // indirect

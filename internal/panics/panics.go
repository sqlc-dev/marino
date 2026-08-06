// Package panics exists because the parser package's generated token
// constants include `any` and `recover` (the MySQL keywords), shadowing
// the Go builtins at package scope. Handle is usable as a deferred
// function from that package: recover only works when called directly by
// a function the panicking goroutine deferred.
package panics

// Handle calls h with the recovered value if the goroutine is panicking.
// Use as: defer panics.Handle(func(e interface{}) { ... }).
func Handle(h func(e interface{})) {
	if e := recover(); e != nil {
		h(e)
	}
}

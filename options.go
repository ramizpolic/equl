package equl

type options struct {
	filter        func(string) bool
	includeFields []string
	ignoreFields  []string
	onlyEqual     bool
}

// WithFieldFilter can be used to specify a custom key filter to use for
// when doing Diff or Equal. If the function returns true for the given key,
// then that key will be taken into account when performing these operations.
//
// If WithFieldFilter is specified, WithFields and WithoutFields is ignored.
func WithFieldFilter(filter func(string) bool) func(*options) {
	return func(o *options) {
		o.filter = filter
	}
}

// WithFields specifies which fields to take into account when doing
// Diff or Equal operations. Only these keys will be added.
// If empty, then all keys will be taken into account.
// When combined with WithoutFields, it is possible to specify
// sub-key filtration.
//
// For example, building a request from following values
// WithFields(".Parent"), WithoutFields(".Parent.sub-key")
// means that all sub-keys in ".Parent" will be taken
// into account, except ".Parent.sub-key" keys.
func WithFields(fields ...string) func(*options) {
	return func(o *options) {
		o.includeFields = fields
	}
}

// WithoutFields specifies which fields should be ignored when doing
// Diff or Equal operations. These keys will always be ignored.
// When combined with WithFields, it is possible to specify
// sub-key filtration. For example:
//
// For example, building a request from following values
// WithFields(".Parent"), WithoutFields(".Parent.sub-key")
// means that all sub-keys in ".Parent" will be taken
// into account, except ".Parent.sub-key" keys.
func WithoutFields(fields ...string) func(*options) {
	return func(o *options) {
		o.ignoreFields = fields
	}
}

func withOnlyEqual() func(*options) {
	return func(o *options) {
		o.onlyEqual = true
	}
}

package filter

// PolicyVersion is part of consumers' result keys. Change it whenever the
// default file universe changes, even if the annotation format stays the same.
const PolicyVersion = "file-scope-v1"

// DefaultExcludePathSegments returns an owned copy of the system exclusions.
// These are directory components, never prefix/glob patterns.
func DefaultExcludePathSegments() []string {
	return []string{"vendor", "node_modules"}
}

// WithDefaultPathSegments preserves caller rules and adds the fixed policy.
func WithDefaultPathSegments(opts *Options) *Options {
	result := Options{}
	if opts != nil {
		result = *opts
	}
	result.ExcludePathSegments = DefaultExcludePathSegments()
	if opts != nil {
		for _, segment := range opts.ExcludePathSegments {
			found := false
			for _, existing := range result.ExcludePathSegments {
				if existing == segment {
					found = true
					break
				}
			}
			if !found && segment != "" {
				result.ExcludePathSegments = append(result.ExcludePathSegments, segment)
			}
		}
	}
	return &result
}

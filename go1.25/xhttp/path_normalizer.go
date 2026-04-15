package xhttp

import "net/http"

// TODO: option to enable caching and check whether it's enabled (if e.g. RoundTripWrap receives
// path normalized with caching disabled, then it enables it for itself only = makes a copy with caching enabled).

// HTTPPathNormalizer returns normalized path with low cardinality.
// For example it may be able to identify that request path '/users/sDaf34Fb9'
// is a variabtion of normalized path '/users/:id'.
// It returns nil if no normalized path found.
type PathNormalizer func(*http.Request) *string

var NoopPathNormalizer = func(*http.Request) *string { return nil }

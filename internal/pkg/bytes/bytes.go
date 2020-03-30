// Package bytes wraps a byte slice for GraphQL un/marshaling.
package bytes

import (
	"io"

	"github.com/99designs/gqlgen/graphql"
)

// Bytes is a byte slice of variable length.
type Bytes []byte

// MarshalGQL marshals a byte slice to a GraphQL string.
func (b Bytes) MarshalGQL(w io.Writer) {
	graphql.MarshalString(string(b)).MarshalGQL(w)
}

// UnmarshalGQL unmarshals a GraphQL string into a byte slice.
func (b *Bytes) UnmarshalGQL(v interface{}) error {
	s, err := graphql.UnmarshalString(v)
	*b = []byte(s)

	return err
}

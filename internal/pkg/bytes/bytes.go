package bytes

import (
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type Bytes []byte

func (b Bytes) MarshalGQL(w io.Writer) {
	graphql.MarshalString(string(b)).MarshalGQL(w)
}

func (b *Bytes) UnmarshalGQL(v interface{}) error {
	s, err := graphql.UnmarshalString(v)
	*b = []byte(s)
	return err
}

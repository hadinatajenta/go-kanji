package identity

import (
	"fmt"

	hashids "github.com/speps/go-hashids/v2"
)

// UserReferenceEncoder encodes and decodes user identifiers into opaque references.
type UserReferenceEncoder struct {
	hash *hashids.HashID
}

// NewUserReferenceEncoder constructs a new encoder using the provided salt.
func NewUserReferenceEncoder(salt string) (*UserReferenceEncoder, error) {
	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = 12

	h, err := hashids.NewWithData(data)
	if err != nil {
		return nil, fmt.Errorf("create hashids encoder: %w", err)
	}

	return &UserReferenceEncoder{hash: h}, nil
}

// Encode converts the given numeric identifier into an opaque string reference.
func (e *UserReferenceEncoder) Encode(id int64) (string, error) {
	ref, err := e.hash.EncodeInt64([]int64{id})
	if err != nil {
		return "", fmt.Errorf("encode user id: %w", err)
	}

	return ref, nil
}

// Decode restores the numeric identifier from a string reference.
func (e *UserReferenceEncoder) Decode(reference string) (int64, error) {
	values, err := e.hash.DecodeInt64WithError(reference)
	if err != nil || len(values) == 0 {
		return 0, fmt.Errorf("decode user reference: %w", err)
	}

	return values[0], nil
}

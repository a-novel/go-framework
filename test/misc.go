package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	ErrNil = fmt.Errorf("nil error")
)

func RequireError(t *testing.T, expect, err error) {
	if expect == nil {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		require.ErrorIs(t, err, expect)
	}
}

func Concat[T any](a ...[]T) []T {
	var result []T
	for _, arr := range a {
		result = append(result, arr...)
	}
	return result
}

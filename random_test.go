package goframework_test

import (
	goframework "github.com/a-novel/go-framework"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestRandom(t *testing.T) {
	r20 := regexp.MustCompile(`[a-zA-Z]{20}`)
	r100 := regexp.MustCompile(`[a-zA-Z]{100}`)

	// Result should be of correct length.
	res := goframework.Random(100)
	require.Len(t, res, 100)
	require.Regexp(t, r100, res)

	res = goframework.Random(20)
	require.Len(t, res, 20)
	require.Regexp(t, r20, res)

	// Result should be different each call.
	res2 := goframework.Random(20)
	require.Len(t, res2, 20)
	require.NotEqual(t, res, res2)
}

func TestRandomAlphabet(t *testing.T) {
	r1 := regexp.MustCompile(`(foo|bar|qux){3}`)
	r2 := regexp.MustCompile(`(foo|bar|qux)((-foo|-bar|-qux){2})`)
	alphabet := []string{"foo", "bar", "qux"}

	// Result should be of correct length.
	res := goframework.RandomAlphabet(alphabet, 3, "")
	require.Len(t, res, 9)
	require.Regexp(t, r1, res)

	// Join should work.
	res3 := goframework.RandomAlphabet(alphabet, 3, "-")
	require.Len(t, res3, 11)
	require.Regexp(t, r2, res3)
}

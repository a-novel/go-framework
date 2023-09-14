package goframework_test

import (
	"fmt"
	"github.com/a-novel/go-framework"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	fakeErr = fmt.Errorf("fake error")
)

func TestWaitForPing(t *testing.T) {
	t.Log("WARNING: this test uses timeouts and may take a certain time to complete")

	t.Run("Success", func(t *testing.T) {
		require.NoError(t, goframework.WaitForPingAuto(func() error {
			return nil
		}))
	})

	t.Run("Success/WithRetries", func(t *testing.T) {
		var count int

		require.NoError(t, goframework.WaitForPingAuto(func() error {
			if count > 3 {
				return nil
			}

			count++
			return fakeErr
		}))
	})

	t.Run("Error/Timeout", func(t *testing.T) {
		require.ErrorIs(t, goframework.WaitForPingAuto(func() error {
			return fakeErr
		}), fakeErr)
	})
}

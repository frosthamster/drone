//go:build e2e

package leetcode_test

import (
	"context"
	"testing"

	"github.com/boar-d-white-foundation/drone/leetcode"
	"github.com/stretchr/testify/require"
)

func TestE2EGetDailyLink(t *testing.T) {
	ctx := context.Background()
	link, err := leetcode.GetDailyLink(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, link)
}
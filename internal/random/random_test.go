package random_test

import (
	"testing"

	"github.com/asankov/shortener/internal/random"
	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	random := random.New()

	id := random.ID(5)
	require.Len(t, id, 5)

	pwd := random.Password(10)
	require.Len(t, pwd, 10)
}

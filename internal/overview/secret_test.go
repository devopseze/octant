package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestSecretData_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewSecretData("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestSecretData(t *testing.T) {
	v := NewSecretData("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()
	c := cache.NewMemoryCache()

	secret := loadFromFile(t, "secret-1.yaml")

	got, err := v.Content(ctx, secret, c)
	require.NoError(t, err)

	dataSection := content.NewSection()
	dataSection.AddText("ca.crt", "1025 bytes")
	dataSection.AddText("namespace", "8 bytes")
	dataSection.AddText("token", "token")

	dataSummary := content.NewSummary("Data", []content.Section{dataSection})

	expected := []content.Content{
		&dataSummary,
	}

	assert.Equal(t, got, expected)
}

package memory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	storage := NewStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key, err := storage.Set(ctx, "https://site.com")
	assert.NotEmpty(t, key)
	assert.NoError(t, err)

	uri, err := storage.Get(ctx, key)
	assert.Equal(t, uri, "https://site.com")
	assert.NoError(t, err)

	uri, err = storage.Get(ctx, "not-exists")
	assert.Empty(t, uri)
	assert.Error(t, err)
}

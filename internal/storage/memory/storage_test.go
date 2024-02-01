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

	key, err := storage.Set(ctx, "https://site.com", 0)
	assert.NotEmpty(t, key)
	assert.NoError(t, err)

	uri, err := storage.Get(ctx, key, 0)
	assert.Equal(t, uri, "https://site.com")
	assert.NoError(t, err)

	uri, err = storage.Get(ctx, "not-exists", 0)
	assert.Empty(t, uri)
	assert.Error(t, err)

	data, err := storage.GetByUser(ctx, 0)
	assert.Nil(t, data)
	assert.NoError(t, err)
}

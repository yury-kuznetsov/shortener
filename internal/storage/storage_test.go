package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	storage, err := NewStorage("")
	assert.NoError(t, err)

	key, err := storage.Set("https://site.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	uri := storage.Get(key)
	assert.Equal(t, uri, "https://site.com")
}

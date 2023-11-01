package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	storage := NewStorage("")

	key := storage.Set("https://site.com")
	assert.NotEmpty(t, key)

	uri := storage.Get(key)
	assert.Equal(t, uri, "https://site.com")
}

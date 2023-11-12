package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	storage, err := NewStorage("")
	assert.NoError(t, err)

	key, err := storage.Set("https://site.com")
	assert.NotEmpty(t, key)
	assert.NoError(t, err)

	uri, err := storage.Get(key)
	assert.Equal(t, uri, "https://site.com")
	assert.NoError(t, err)

	uri, err = storage.Get("not-exists")
	assert.Empty(t, uri)
	assert.Error(t, err)
}

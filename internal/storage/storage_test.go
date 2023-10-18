package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	key := ArrStorage.Set("https://site.com")
	assert.NotEmpty(t, key)

	uri := ArrStorage.Get(key)
	assert.Equal(t, uri, "https://site.com")
}

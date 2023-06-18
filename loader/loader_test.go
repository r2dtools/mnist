package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	train, test, err := LoadData("")

	assert.Nil(t, err)
	assert.Equal(t, 60000, train.Count())
	assert.Equal(t, 10000, test.Count())
}

package scm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionNumber(t *testing.T) {
	v := versionNumber("1.7")

	assert.True(t, v.meets(versionNumber("1")), "version 1.7 should meet 1")
	assert.True(t, v.meets(versionNumber("1.6")), "version 1.7 should meet 1.6")
	assert.True(t, v.meets(versionNumber("1.6.2")), "version 1.7 should meet 1.6.2")
	assert.True(t, v.meets(versionNumber("1.7")), "version 1.7 should meet 1.7")
	assert.True(t, v.meets(versionNumber("1.7.0")), "version 1.7 should meet 1.7.0")
	assert.False(t, v.meets(versionNumber("2")), "version 1.7 should not meet 2")
	assert.False(t, v.meets(versionNumber("2.0.7")), "version 1.7 should not meet 2.0.7")
}

package tool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompareVer(t *testing.T) {
	// https://cs.opensource.google/go/x/mod/+/refs/tags/v0.5.1:semver/semver_test.go
	actual := compareVer("v2.6.0", "v2.7.0")
	assert.True(t, actual < 0)
}

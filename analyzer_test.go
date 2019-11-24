package prodinspect

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFromFileSystem(t *testing.T) {
	assert := assert.New(t)

	testdata := analysistest.TestData()
	rs := analysistest.Run(t, testdata, Analyzer, "a")

	assert.IsType((*Inspector)(nil), rs[0].Result)
}

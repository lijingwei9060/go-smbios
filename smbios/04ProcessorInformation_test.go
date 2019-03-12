package smbios_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	b := []byte{0xAB, 0xAC, 0x1C, 0x1B}
	s := fmt.Sprintf("%X %X %X %X", b[0], b[1], b[2], b[3])
	t.Logf(s)
	assert.NotNil(t, s)
}

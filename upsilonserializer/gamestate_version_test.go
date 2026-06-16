package upsilonserializer_test

import (
	"testing"

	"github.com/ecumeurs/upsilonserializer"
)

// TestCurrentSerializerVersion_IsPositive ensures the constant is non-zero so that
// a zero-value (absent) field can always be distinguished from a real version.
func TestCurrentSerializerVersion_IsPositive(t *testing.T) {
	if upsilonserializer.CurrentSerializerVersion <= 0 {
		t.Fatalf("CurrentSerializerVersion must be > 0, got %d", upsilonserializer.CurrentSerializerVersion)
	}
}

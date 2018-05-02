package keys

import (
	"testing"

	"github.com/lioneagle/goutil/src/test"
)

func TestKeyContextUnmarshalError(t *testing.T) {
	var context KeyContext
	err := context.UnmarshalJSON([]byte(``))

	test.EXPECT_NE(t, err, nil, "")
}

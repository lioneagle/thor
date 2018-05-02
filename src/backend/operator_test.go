package backend

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lioneagle/goutil/src/test"
)

func TestOpUnmarshalJSONError(t *testing.T) {
	var op Op
	err := op.UnmarshalJSON(nil)

	test.EXPECT_NE(t, err, nil, "")
}

func TestOpUnmarshalJSONOk(t *testing.T) {
	testdata := []struct {
		src string
	}{
		{""},
		{"not_equal"},
		{"regex_match"},
		{"not_regex_match"},
		{"regex_contains"},
		{"not_regex_contains"},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			var op Op
			d, err := json.Marshal(v.src)
			test.ASSERT_EQ(t, err, nil, "err =", err)

			err = op.UnmarshalJSON(d)
			test.EXPECT_EQ(t, err, nil, "err =", err)
		})
	}
}

package keys

import (
	"fmt"
	"testing"

	"github.com/lioneagle/goutil/src/test"
)

func TestKeyString(t *testing.T) {
	testdata := []struct {
		key Key
		str string
	}{
		{UP, "up"},
		{LEFT, "left"},
		{RIGHT, "right"},
		{DOWN, "down"},
		{ENTER, "enter"},
		{Key('\t'), "tab"},
		{ESCAPE, "escape"},
		{Key(' '), "space"},
		{F1, "f1"},
		{F2, "f2"},
		{F3, "f3"},
		{F4, "f4"},
		{F5, "f5"},
		{F6, "f6"},
		{F7, "f7"},
		{F8, "f8"},
		{F9, "f9"},
		{F10, "f10"},
		{F11, "f11"},
		{F12, "f12"},
		{BACKSPACE, "backspace"},
		{DELETE, "delete"},
		{INSERT, "insert"},
		{PAGE_UP, "pageup"},
		{PAGE_DOWN, "pagedown"},
		{HOME, "home"},
		{END, "end"},
		{BREAK, "break"},
		{Key('/'), "forward_slash"},
		{Key('`'), "backquote"},
		{Key('"'), "\\\""},
		{Key('+'), "plus"},
		{Key('-'), "minus"},
		{Key('='), "equals"},
		{ANY, "<character>"},
		{Key('i'), "i"},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			test.EXPECT_EQ(t, v.key.String(), v.str, "")
		})
	}
}

func TestStr2Key(t *testing.T) {
	testdata := []struct {
		key Key
		str string
	}{
		{UP, "up"},
		{LEFT, "left"},
		{RIGHT, "right"},
		{DOWN, "down"},
		{ENTER, "enter"},
		{Key('\t'), "tab"},
		{ESCAPE, "escape"},
		{Key(' '), "space"},
		{F1, "f1"},
		{F2, "f2"},
		{F3, "f3"},
		{F4, "f4"},
		{F5, "f5"},
		{F6, "f6"},
		{F7, "f7"},
		{F8, "f8"},
		{F9, "f9"},
		{F10, "f10"},
		{F11, "f11"},
		{F12, "f12"},
		{BACKSPACE, "backspace"},
		{DELETE, "delete"},
		{INSERT, "insert"},
		{PAGE_UP, "pageup"},
		{PAGE_DOWN, "pagedown"},
		{HOME, "home"},
		{END, "end"},
		{BREAK, "break"},
		{Key('/'), "forward_slash"},
		{Key('`'), "backquote"},
		{Key('"'), "\\\""},
		{Key('+'), "plus"},
		{Key('-'), "minus"},
		{Key('='), "equals"},
		{ANY, "<character>"},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			key, ok := g_str2Key[v.str]
			test.ASSERT_TRUE(t, ok, "str = \"%s\"", v.str)
			test.EXPECT_EQ(t, key, v.key, "")
		})
	}
}

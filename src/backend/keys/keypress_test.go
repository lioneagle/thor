package keys

import (
	"fmt"
	"testing"

	"github.com/lioneagle/goutil/src/test"
)

func TestKeyPressIndex(t *testing.T) {
	testdata := []struct {
		key   Key
		shift bool
		super bool
		alt   bool
		ctrl  bool
		index int
	}{
		{'a', false, false, false, false, 'a'},
		{'a', true, false, false, false, 'a' + SHIFT},
		{'a', true, true, false, false, 'a' + SHIFT + SUPER},
		{'a', true, true, true, false, 'a' + SHIFT + SUPER + ALT},
		{'a', true, true, true, true, 'a' + SHIFT + SUPER + ALT + CTRL},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			keypress := KeyPress{Key: v.key, Shift: v.shift, Super: v.super, Alt: v.alt, Ctrl: v.ctrl}
			test.EXPECT_EQ(t, keypress.Index(), v.index, "")
		})
	}
}

func TestKeyPressIsCharacter(t *testing.T) {
	testdata := []struct {
		key    Key
		shift  bool
		super  bool
		alt    bool
		ctrl   bool
		isChar bool
	}{
		{'a', false, false, false, false, true},
		{'a', true, false, false, false, true},
		{'a', false, true, false, false, false},
		{'a', false, false, false, true, false},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			keypress := KeyPress{Key: v.key, Shift: v.shift, Super: v.super, Alt: v.alt, Ctrl: v.ctrl}
			test.EXPECT_EQ(t, keypress.IsCharacter(), v.isChar, "")
		})
	}
}

func TestKeyPressFix(t *testing.T) {
	keypress := KeyPress{"A", 'A', false, false, false, false}
	keypress.fix()
	test.EXPECT_EQ(t, keypress.Key, Key('a'), "")
	test.EXPECT_TRUE(t, keypress.Shift, "")
}

func TestKeyPressUnmarshalJSON(t *testing.T) {
	var keypress KeyPress
	d := `"super+ctrl+alt+shift+f1+λλλ"`
	err := keypress.UnmarshalJSON([]byte(d))
	test.EXPECT_EQ(t, err, nil, "")
}

func TestKeyPressString(t *testing.T) {
	testdata := []struct {
		key   Key
		shift bool
		super bool
		alt   bool
		ctrl  bool
		str   string
	}{
		{'a', false, false, false, false, "a"},
		{'a', true, false, false, false, "shift+a"},
		{'a', true, true, false, false, "super+shift+a"},
		{'a', true, true, true, false, "super+alt+shift+a"},
		{'a', true, true, true, true, "super+ctrl+alt+shift+a"},
	}

	for i, v := range testdata {
		v := v

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			keypress := KeyPress{Key: v.key, Shift: v.shift, Super: v.super, Alt: v.alt, Ctrl: v.ctrl}
			test.EXPECT_EQ(t, keypress.String(), v.str, "")
		})
	}
}

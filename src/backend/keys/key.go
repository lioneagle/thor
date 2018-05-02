package keys

import (
	"unicode"
)

type Key rune

const (
	LEFT Key = 0x2190 + iota
	UP
	RIGHT
	DOWN

	ENTER        = '\n'
	ESCAPE       = 0x001b
	BACKSPACE    = 0x0008
	DELETE       = 0x007f
	KEYPAD_ENTER = '\n'

	F1 Key = 0x2701 + iota
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	INSERT
	PAGE_UP
	PAGE_DOWN
	HOME
	END
	BREAK
	ANY Key = unicode.MaxRune
)

const (
	SHIFT = (1 << (29 - iota))
	CTRL
	ALT
	SUPER
)

var g_str2Key = map[string]Key{
	"up":            UP,
	"left":          LEFT,
	"right":         RIGHT,
	"down":          DOWN,
	"enter":         ENTER,
	"tab":           '\t',
	"escape":        ESCAPE,
	"space":         ' ',
	"f1":            F1,
	"f2":            F2,
	"f3":            F3,
	"f4":            F4,
	"f5":            F5,
	"f6":            F6,
	"f7":            F7,
	"f8":            F8,
	"f9":            F9,
	"f10":           F10,
	"f11":           F11,
	"f12":           F12,
	"backspace":     BACKSPACE,
	"delete":        DELETE,
	"keypad_enter":  KEYPAD_ENTER,
	"insert":        INSERT,
	"pageup":        PAGE_UP,
	"pagedown":      PAGE_DOWN,
	"home":          HOME,
	"end":           END,
	"break":         BREAK,
	"forward_slash": '/',
	"backquote":     '`',
	"\\\"":          '"',
	"plus":          '+',
	"minus":         '-',
	"equals":        '=',
	"<character>":   ANY,
}

var g_key2str = map[Key]string{
	UP:        "up",
	LEFT:      "left",
	RIGHT:     "right",
	DOWN:      "down",
	ENTER:     "enter",
	'\t':      "tab",
	ESCAPE:    "escape",
	' ':       "space",
	F1:        "f1",
	F2:        "f2",
	F3:        "f3",
	F4:        "f4",
	F5:        "f5",
	F6:        "f6",
	F7:        "f7",
	F8:        "f8",
	F9:        "f9",
	F10:       "f10",
	F11:       "f11",
	F12:       "f12",
	BACKSPACE: "backspace",
	DELETE:    "delete",
	INSERT:    "insert",
	PAGE_UP:   "pageup",
	PAGE_DOWN: "pagedown",
	HOME:      "home",
	END:       "end",
	BREAK:     "break",
	'/':       "forward_slash",
	'`':       "backquote",
	'"':       "\\\"",
	'+':       "plus",
	'-':       "minus",
	'=':       "equals",
	ANY:       "<character>",
}

func (this Key) String() string {
	if v, ok := g_key2str[this]; ok {
		return v
	}
	return string(this)
}

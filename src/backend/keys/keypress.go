package keys

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/lioneagle/goutil/src/logger"
)

// KeyPress describes a key press event.
// Note that Key does not distinguish between capital and non-capital letters;
// use the Text property for this purpose.
type KeyPress struct {
	Text  string
	Key   Key
	Shift bool
	Super bool
	Alt   bool
	Ctrl  bool
}

// Returns an index used for sorting key presses.
// TODO(.): This is in no way a unique index with quite a lot of collisions and potentially resulting
// in bad lookups.
func (this *KeyPress) Index() int {
	ret := int(this.Key)
	ret = int(this.Key)
	if this.Shift {
		ret += SHIFT
	}
	if this.Alt {
		ret += ALT
	}
	if this.Ctrl {
		ret += CTRL
	}
	if this.Super {
		ret += SUPER
	}
	return ret
}

func (this *KeyPress) IsCharacter() bool {
	return unicode.IsPrint(rune(this.Key)) && !this.Super && !this.Ctrl
}

// Modifies the KeyPress so that it's Key is a unicode lower case
// rune and if it was in uppercase before this modification, the
// "Shift" modifier is also enabled.
func (this *KeyPress) fix() {
	lower := Key(unicode.ToLower(rune(this.Key)))
	if lower != this.Key {
		this.Shift = true
		this.Key = lower
	}
}

func (this *KeyPress) UnmarshalJSON(d []byte) error {
	combo := strings.Split(string(d[1:len(d)-1]), "+")
	for _, c := range combo {
		lower := strings.ToLower(c)
		switch lower {
		case "super":
			this.Super = true
		case "ctrl":
			this.Ctrl = true
		case "alt":
			this.Alt = true
		case "shift":
			this.Shift = true
		default:
			if v, ok := g_str2Key[lower]; ok {
				this.Key = v
			} else {
				r := []Key(c)
				if len(r) != 1 {
					logger.Warning("Unknown key value with %d bytes: %s", len(c), c)
					return nil
				}
				this.Key = Key(c[0])
				this.fix()
			}
		}
	}
	return nil
}

func (this *KeyPress) String() (ret string) {
	if this.Super {
		ret += "super+"
	}
	if this.Ctrl {
		ret += "ctrl+"
	}
	if this.Alt {
		ret += "alt+"
	}
	if this.Shift {
		ret += "shift+"
	}
	ret += fmt.Sprintf("%s", this.Key)
	return ret
}

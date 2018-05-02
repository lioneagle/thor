package keys

import (
	"encoding/json"
	"fmt"
	"sort"

	"backend"

	"github.com/lioneagle/goutil/src/buffer"
	"github.com/lioneagle/goutil/src/profile"
)

// A single KeyBinding for which after pressing the given
// sequence of Keys, and the Context matches,
// the Command will be invoked with the provided Args.
type KeyBinding struct {
	Keys     []KeyPress
	Command  int
	Args     map[string]interface{}
	Context  []KeyContext
	priority int
}

// An utility struct that(same as HasSettings ) is typically embedded in
// other type structs to make that type implement the KeyBindingsInterface
type HasKeyBindings struct {
	keybindings KeyBindings
}

func (this *HasKeyBindings) KeyBindings() *KeyBindings {
	return &this.keybindings
}

type KeyBindingsInterface interface {
	KeyBindings() *KeyBindings
}

type KeyBindings struct {
	Bindings []*KeyBinding
	seqIndex int
	parent   KeyBindingsInterface
}

func (this *KeyBindings) Len() int {
	return len(this.Bindings)
}

/*func (this *KeyBindings) Swap(i, j int) {
	this.Bindings[i], this.Bindings[j] = this.Bindings[j], this.Bindings[i]
}

func (this *KeyBindings) Less(i, j int) bool {
	return this.Bindings[i].Keys[this.seqIndex].Index() < this.Bindings[j].Keys[this.seqIndex].Index()
}*/

// Drops all KeyBindings that are a sequence of key presses less or equal
// to the given number.
func (this *KeyBindings) DropLessEqualKeys(count int) {
	for {
		for i := 0; i < len(this.Bindings); {
			if len(this.Bindings[i].Keys) <= count {
				this.Bindings[i] = this.Bindings[len(this.Bindings)-1]
				this.Bindings = this.Bindings[:len(this.Bindings)-1]
			} else {
				i++
			}
		}
		sort.Slice(this.Bindings, func(i, j int) bool {
			return this.Bindings[i].Keys[this.seqIndex].Index() < this.Bindings[j].Keys[this.seqIndex].Index()
		})
		if this.parent == nil {
			break
		}
		this = this.parent.KeyBindings()
	}
}

func (this *KeyBindings) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, &this.Bindings); err != nil {
		return err
	}
	for i := range this.Bindings {
		this.Bindings[i].priority = i
	}
	this.DropLessEqualKeys(0)
	return nil
}

func (this *KeyBindings) SetParent(p KeyBindingsInterface) {
	this.parent = p
	// All parents and childs seqIndex must be equal
	p.KeyBindings().seqIndex = this.seqIndex
}

func (this *KeyBindings) Parent() KeyBindingsInterface {
	return this.parent
}

func (this *KeyBindings) filter(ki int, ret *KeyBindings) {
	for {
		idx := sort.Search(this.Len(), func(i int) bool {
			return this.Bindings[i].Keys[this.seqIndex].Index() >= ki
		})
		for i := idx; i < len(this.Bindings) && this.Bindings[i].Keys[this.seqIndex].Index() == ki; i++ {
			ret.Bindings = append(ret.Bindings, this.Bindings[i])
		}
		if this.parent == nil {
			break
		}
		this = this.parent.KeyBindings()
		if ret.parent == nil {
			ret.SetParent(new(HasKeyBindings))
		}
		ret = ret.parent.KeyBindings()
	}
}

// Filters the KeyBindings, returning a new KeyBindings object containing
// a subset of matches for the given key press.
func (this *KeyBindings) Filter(keypress KeyPress) (ret KeyBindings) {
	p := profile.Enter("key.filter")
	defer p.Exit()

	keypress.fix()
	this.DropLessEqualKeys(this.seqIndex)
	ret.seqIndex = this.seqIndex + 1
	ki := keypress.Index()

	this.filter(ki, &ret)

	if keypress.IsCharacter() {
		this.filter(int(ANY), &ret)
	}
	return
}

// Tries to resolve all the current KeyBindings in k to a single
// action. If any action is appropriate as determined by context,
// the return value will be the specific KeyBinding that is possible
// to execute now, otherwise it is nil.
func (this *KeyBindings) Action(qc func(key string, operator backend.Op, operand interface{}, match_all bool) bool) (keybinding *KeyBinding) {
	p := profile.Enter("key.action")
	defer p.Exit()

	for {
		for i := range this.Bindings {
			if len(this.Bindings[i].Keys) > this.seqIndex {
				// This key binding is of a key sequence longer than what is currently
				// probed for. For example, the binding is for the sequence ['a','b','c'], but
				// the user has only pressed ['a','b'] so far.
				continue
			}
			for _, c := range this.Bindings[i].Context {
				if !qc(c.Key, c.Operator, c.Operand, c.MatchAll) {
					goto skip
				}
			}
			if keybinding == nil || keybinding.priority < this.Bindings[i].priority {
				keybinding = this.Bindings[i]
			}
		skip:
		}
		if keybinding != nil || this.parent == nil {
			break
		}
		this = this.parent.KeyBindings()
	}
	return
}

func (this *KeyBindings) SeqIndex() int {
	return this.seqIndex
}

func (this KeyBindings) String() string {
	buf := buffer.NewByteBuffer(nil)
	for _, b := range this.Bindings {
		buf.WriteString(fmt.Sprintf("%+v\n", b))
	}
	return buf.String()
}

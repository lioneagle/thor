package keys

import (
	"encoding/json"

	"backend"
)

type rawKeyContext struct {
	Key      string      //< The context's name.
	Operator backend.Op  //< The operation to perform.
	Operand  interface{} //< The operand on which this operation should be performed.
	MatchAll bool        `json:"match_all"` //< Whether all selections should match the context or if it's enough for just one to match.
}

type KeyContext struct {
	rawKeyContext
}

func (this *KeyContext) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, &this.rawKeyContext); err != nil {
		return err
	}
	if this.Operand == nil {
		this.Operand = true
	}
	return nil
}

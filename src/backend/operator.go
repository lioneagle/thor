package backend

import (
	"encoding/json"
)

type Op int

const (
	OpEqual            Op = iota //< Compare for equality.
	OpNotEqual                   //< Compare for difference.
	OpRegexMatch                 //< Compare for a regular expression match.
	OpNotRegexMatch              //< Compare for a regular expression difference.
	OpRegexContains              //< Compare whether the given regular expression matches some substring of the operand.
	OpNotRegexContains           //< Compare whether the given regular expression does not match some substring of the operand.
)

func (o *Op) UnmarshalJSON(d []byte) error {
	var tmp string
	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}
	switch tmp {
	default:
		*o = OpEqual
	case "not_equal":
		*o = OpNotEqual
	case "regex_match":
		*o = OpRegexMatch
	case "not_regex_match":
		*o = OpNotRegexMatch
	case "regex_contains":
		*o = OpRegexContains
	case "not_regex_contains":
		*o = OpNotRegexContains
	}

	return nil
}

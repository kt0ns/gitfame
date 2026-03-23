package flags

import (
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type Value interface {
	pflag.Value
}

var _ Value = &StringValue{}
var _ Value = &BoolValue{}
var _ Value = &StringSliceValue{}

type StringValue struct{ val string }

func (s *StringValue) Set(val string) error {
	s.val = val
	return nil
}

func (s *StringValue) Type() string {
	return "string"
}

func (s *StringValue) String() string {
	return s.val
}

func (s *StringValue) Value() string {
	return s.val
}

func newStringValue(def string) *StringValue {
	return &StringValue{val: def}
}

type BoolValue struct {
	val bool
}

func (b *BoolValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	b.val = v
	return nil
}
func (b *BoolValue) Type() string {
	return "bool"
}

func (b *BoolValue) String() string {
	return strconv.FormatBool(b.val)
}

func (b *BoolValue) IsBoolFlag() bool {
	return true
}

func (b *BoolValue) Value() bool {
	return b.val
}

func newBoolValue(def bool) *BoolValue {
	return &BoolValue{val: def}
}

type StringSliceValue struct {
	val     []string
	changed bool
}

func (s *StringSliceValue) Set(val string) error {
	if !s.changed {
		s.val = []string{}
		s.changed = true
	}
	parts := strings.Split(val, ",")
	for _, p := range parts {
		s.val = append(s.val, strings.TrimSpace(p))
	}
	return nil
}
func (s *StringSliceValue) Type() string {
	return "stringSlice"
}

func (s *StringSliceValue) String() string {
	return strings.Join(s.val, ",")
}

func (s *StringSliceValue) Value() []string {
	return s.val
}

func newStringSliceValue(def []string) *StringSliceValue {
	copied := make([]string, len(def))
	copy(copied, def)
	return &StringSliceValue{val: copied}
}

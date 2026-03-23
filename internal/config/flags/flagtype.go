package flags

import (
	stderrors "errors"
	"fmt"

	"gitlab.com/slon/shad-go/gitfame/internal/errors"
)

type Flag struct {
	Name  string
	Use   string
	Value Value
}

var (
	ErrNotString      = stderrors.New("cannot get string value")
	ErrNotBool        = stderrors.New("cannot get bool value")
	ErrNotStringSlice = stderrors.New("cannot get stringSlice value")
)

func (f *Flag) GetString() (string, error) {
	switch val := f.Value.(type) {
	case *StringValue:
		return val.Value(), nil
	default:
		return "", fmt.Errorf(errors.MsgFlagTypeMismatch, ErrNotString, f.Name, f.Value.Type())
	}
}

func (f *Flag) GetBool() (bool, error) {
	switch val := f.Value.(type) {
	case *BoolValue:
		return val.Value(), nil
	default:
		return false, fmt.Errorf(errors.MsgFlagTypeMismatch, ErrNotBool, f.Name, f.Value.Type())
	}
}

func (f *Flag) GetStringSlice() ([]string, error) {
	switch val := f.Value.(type) {
	case *StringSliceValue:
		return val.Value(), nil
	default:
		return nil, fmt.Errorf(errors.MsgFlagTypeMismatch, ErrNotStringSlice, f.Name, f.Value.Type())
	}
}

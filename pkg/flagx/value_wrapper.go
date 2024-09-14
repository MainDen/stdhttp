package flagx

import (
	"fmt"
	"strings"
)

type Wrapper interface {
	Wrap(value Value) Value
}

func Wrap(value Value, wrappers ...Wrapper) Value {
	for _, wrapper := range wrappers {
		value = wrapper.Wrap(value)
	}
	return value
}

func WithFunc(fn func(value Value) Value) *wrapperFunc {
	return &wrapperFunc{fn: fn}
}

type wrapperFunc struct {
	fn func(value Value) Value
}

func (wrapper *wrapperFunc) Wrap(value Value) Value {
	return wrapper.fn(value)
}

type wrapperInlined struct {
	inlined bool
}

func WithInlined(inlined bool) *wrapperInlined {
	return &wrapperInlined{inlined: inlined}
}

type valueInlined struct {
	value   Value
	inlined bool
}

func WrapInlined(value Value, inlined bool) *valueInlined {
	return &valueInlined{value: value, inlined: inlined}
}

func (value *valueInlined) Parse(args ...string) (int, error) {
	return value.value.Parse(args...)
}

func (value *valueInlined) Format() []string {
	return FormatSlice(value.value)
}

func (value *valueInlined) IsInlined() bool {
	return value.inlined
}

type wrapperDefaults struct {
	defaults []string
}

func WithDefaults(defaults ...string) *wrapperDefaults {
	return &wrapperDefaults{defaults: defaults}
}

func (wrapper *wrapperDefaults) Wrap(value Value) Value {
	return WrapDefaults(value, wrapper.defaults...)
}

type valueDefaults struct {
	value    Value
	defaults []string
}

func WrapDefaults(value Value, defaults ...string) *valueDefaults {
	return &valueDefaults{value: value, defaults: defaults}
}

func (value *valueDefaults) Parse(args ...string) (int, error) {
	return value.value.Parse(args...)
}

func (value *valueDefaults) Format() []string {
	return value.defaults
}

func (value *valueDefaults) IsInlined() bool {
	return IsInlined(value.value)
}

type valueJoin struct {
	values []Value
}

func Join(values ...Value) *valueJoin {
	return &valueJoin{values: values}
}

func (value *valueJoin) Parse(args ...string) (int, error) {
	var n int
	for _, value := range value.values {
		m, err := value.Parse(args...)
		n += m
		if err != nil {
			return n, err
		}
		args = args[m:]
	}
	return n, nil
}

func (value *valueJoin) Format() []string {
	var result []string
	for _, value := range value.values {
		result = append(result, FormatSlice(value)...)
	}
	return result
}

type valueOptional struct {
	value Value
}

func Optional(value Value) *valueOptional {
	return &valueOptional{value: value}
}

func (value *valueOptional) Parse(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}
	return value.value.Parse(args...)
}

func (value *valueOptional) Format() []string {
	return FormatSlice(value.value)
}

func (value *valueOptional) IsInlined() bool {
	return IsInlined(value.value)
}

type wrapperArgs struct {
	args []string
}

func WithArgs(args ...string) *wrapperArgs {
	return &wrapperArgs{args: args}
}

func (wrapper *wrapperArgs) Wrap(value Value) Value {
	return Args(value, wrapper.args...)
}

type valueArgs struct {
	value Value
	args  []string
}

func Args(value Value, args ...string) *valueArgs {
	return &valueArgs{value: value, args: args}
}

func (value *valueArgs) Parse(args ...string) (int, error) {
	n, err := value.value.Parse(value.args...)
	if err != nil {
		return 0, err
	}
	if n < len(value.args) {
		if len(args)-n == 0 {
			panic(fmt.Errorf("extra argument: %v", formatQuote(value.args[n])))
		}
		panic(fmt.Errorf("extra arguments: %v", strings.Join(formatQuoteList(value.args[n:]), " ")))
	}
	return 0, nil
}

type wrapperEnum struct {
	values []string
}

func WithEnum(values ...string) *wrapperEnum {
	return &wrapperEnum{values: values}
}

func (wrapper *wrapperEnum) Wrap(value Value) Value {
	return Enum(value, wrapper.values...)
}

type valueEnum struct {
	value  Value
	enum   map[string]struct{}
	values []string
}

func Enum(value Value, values ...string) *valueEnum {
	enum := make(map[string]struct{}, len(values))
	for _, value := range values {
		enum[value] = struct{}{}
	}
	return &valueEnum{value: value, enum: enum, values: values}
}

func (value *valueEnum) Parse(args ...string) (int, error) {
	n, err := value.value.Parse(args...)
	if err != nil {
		return n, err
	}
	for i := 0; i < n; i++ {
		if _, ok := value.enum[args[i]]; !ok {
			return i, fmt.Errorf("invalid value: %v (Allowed one of: %v)", formatQuote(args[i]), strings.Join(formatQuoteList(value.values), ", "))
		}
	}
	return n, nil
}

func (value *valueEnum) Format() []string {
	return FormatSlice(value.value)
}

func (value *valueEnum) IsInlined() bool {
	return IsInlined(value.value)
}

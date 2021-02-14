/*
   Copyright 2021 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package expect

import (
	"errors"
	"fmt"
	"go/types"
	"strconv"
	"strings"
)

var ErrInvalidCondition = errors.New("Invalid condition")

const ConfigBoolKeyName = "key_bool"
const ConfigStrKeyName = "key_str"
const ConfigIntKeyName = "key_int"
const ConfigUintKeyName = "key_uint"
const ConfigDoubleKeyName = "key_double"

type Condition struct {
	ctype  types.BasicKind
	ccase  int
	cvalue interface{}
}
type TypeCondition struct {
	Keys             Keys
	Condition        Condition
	TypeConditionStr string
}

const (
	CaseInvalid     = iota
	CaseGt          //  >
	CaseGe          //  >=
	CaseLt          //  <
	CaseLe          //  <=
	CaseEq          //  ==
	CaseNe          //  !=
	CaseContains    // for string.
	CaseNotContains // for string.
)

// Str2IntCase converts string case to int case.
//  e.g. ">=" ->  CaseGe
func Str2IntCase(s string) int {
	ret := CaseInvalid
	switch s {
	case ">":
		ret = CaseGt
	case ">=":
		ret = CaseGe
	case "<":
		ret = CaseLt
	case "<=":
		ret = CaseLe
	case "==":
		ret = CaseEq
	case "!=":
		ret = CaseNe
	case "contains":
		ret = CaseContains
	case "not_contains":
		ret = CaseNotContains
	}
	return ret
}

// IntCase2Str converts int case to string case.
//  e.g. CaseGe -> ">="
func IntCase2Str(i int) string {
	ret := "Invalid"
	switch i {
	case CaseGt:
		ret = ">"
	case CaseGe:
		ret = ">="
	case CaseLt:
		ret = "<"
	case CaseLe:
		ret = "<="
	case CaseEq:
		ret = "=="
	case CaseNe:
		ret = "!="
	case CaseContains:
		ret = "contains"
	case CaseNotContains:
		ret = "not_contains"
	}
	return ret
}

func (c Condition) matchString(s string) bool {
	switch c.ccase {
	case CaseEq:
		return c.cvalue.(string) == s
	case CaseNe:
		return c.cvalue.(string) != s
	case CaseContains:
		return strings.Contains(s, c.cvalue.(string))
	case CaseNotContains:
		return !strings.Contains(s, c.cvalue.(string))
	}
	return false
}

func (c Condition) matchBool(b bool) bool {
	switch c.ccase {
	case CaseEq:
		return c.cvalue.(bool) == b
	case CaseNe:
		return c.cvalue.(bool) != b
	}
	return false
}

func (c Condition) matchInt(i int) bool {
	switch c.ccase {
	case CaseGt:
		return i > c.cvalue.(int)
	case CaseGe:
		return i >= c.cvalue.(int)
	case CaseLt:
		return i < c.cvalue.(int)
	case CaseLe:
		return i <= c.cvalue.(int)
	case CaseEq:
		return c.cvalue.(int) == i
	case CaseNe:
		return c.cvalue.(int) != i
	}
	return false
}

func (c Condition) matchUint(i uint) bool {
	switch c.ccase {
	case CaseGt:
		return i > c.cvalue.(uint)
	case CaseGe:
		return i >= c.cvalue.(uint)
	case CaseLt:
		return i < c.cvalue.(uint)
	case CaseLe:
		return i <= c.cvalue.(uint)
	case CaseEq:
		return c.cvalue.(uint) == i
	case CaseNe:
		return c.cvalue.(uint) != i
	}
	return false
}

func (c Condition) matchDouble(d float64) bool {
	switch c.ccase {
	case CaseGt:
		return d > c.cvalue.(float64)
	case CaseGe:
		return d >= c.cvalue.(float64)
	case CaseLt:
		return d < c.cvalue.(float64)
	case CaseLe:
		return d <= c.cvalue.(float64)
	case CaseEq:
		return c.cvalue.(float64) == d
	case CaseNe:
		return c.cvalue.(float64) != d
	}
	return false
}

// IsMatch check if v matches condition.
func (c Condition) IsMatch(v interface{}) (bool, error) {
	if v == nil {
		return false, errors.New("value is nil")
	}

	switch c.ctype {
	case types.Bool:
		b, ok := v.(bool)
		if ok {
			return c.matchBool(b), nil
		}
	case types.String:
		s, ok := v.(string)
		if ok {
			return c.matchString(s), nil
		}
		// may send byte slice
		ba, ok := v.([]byte)
		if ok {
			return c.matchString(string(ba)), nil
		}
	case types.Int:
		switch v.(type) {
		case int64:
			i := v.(int64)
			return c.matchInt(int(i)), nil
		case int32:
			i := v.(int32)
			return c.matchInt(int(i)), nil
		case int16:
			i := v.(int16)
			return c.matchInt(int(i)), nil
		case int8:
			i := v.(int8)
			return c.matchInt(int(i)), nil
		case int:
			i := v.(int)
			return c.matchInt(i), nil
		case uint64:
			i := v.(uint64)
			return c.matchInt(int(i)), nil
		case uint32:
			i := v.(uint32)
			return c.matchInt(int(i)), nil
		case uint16:
			i := v.(uint16)
			return c.matchInt(int(i)), nil
		case uint8:
			i := v.(uint8)
			return c.matchInt(int(i)), nil
		case uint:
			i := v.(uint)
			return c.matchInt(int(i)), nil
		}

	case types.Uint:
		switch v.(type) {
		case uint64:
			i := v.(uint64)
			return c.matchUint(uint(i)), nil
		case uint32:
			i := v.(uint32)
			return c.matchUint(uint(i)), nil
		case uint16:
			i := v.(uint16)
			return c.matchUint(uint(i)), nil
		case uint8:
			i := v.(uint8)
			return c.matchUint(uint(i)), nil
		case uint:
			i := v.(uint)
			return c.matchUint(i), nil
		case int64:
			i := v.(int64)
			return c.matchUint(uint(i)), nil
		case int32:
			i := v.(int32)
			return c.matchUint(uint(i)), nil
		case int16:
			i := v.(int16)
			return c.matchUint(uint(i)), nil
		case int8:
			i := v.(int8)
			return c.matchUint(uint(i)), nil
		case int:
			i := v.(int)
			return c.matchUint(uint(i)), nil
		}

	case types.Float64:
		switch v.(type) {
		case float64:
			d := v.(float64)
			return c.matchDouble(d), nil
		case float32:
			d := v.(float32)
			return c.matchDouble(float64(d)), nil
		}
	}
	return false, fmt.Errorf("can not cast: type=%d v=%s\n", c.ctype, v)
}

func (c Condition) String() string {
	ret := IntCase2Str(c.ccase) + " "
	switch c.ctype {
	case types.Uint:
		u, ok := c.cvalue.(uint)
		if ok {
			ret += strconv.FormatUint(uint64(u), 10)
		}
	case types.Int:
		i, ok := c.cvalue.(int)
		if ok {
			ret += strconv.FormatInt(int64(i), 10)
		}
	case types.Float64:
		f, ok := c.cvalue.(float64)
		if ok {
			ret += strconv.FormatFloat(f, 'f', -1, 64)
		}
	case types.Bool:
		b, ok := c.cvalue.(bool)
		if ok {
			if b {
				ret += "true"
			} else {
				ret += "false"
			}
		}
	case types.String:
		ret += c.cvalue.(string)
	}
	return ret
}

// NewBoolCondition returns Condition c of boolean.
//  c must be CaseEq or CaseNe.
func NewBoolCondition(c int, b bool) (*Condition, error) {
	if c != CaseEq && c != CaseNe {
		return nil, ErrInvalidCondition
	}
	ret := &Condition{ctype: types.Bool, ccase: c, cvalue: b}

	return ret, nil
}

// NewStringCondition returns Condition c of string.
//  c must be CaseEq, CaseNe CaseContains or CaseNotContains.
func NewStringCondition(c int, s string) (*Condition, error) {
	if c != CaseEq && c != CaseNe && c != CaseContains && c != CaseNotContains {
		return nil, ErrInvalidCondition
	}
	ret := &Condition{ctype: types.String, ccase: c, cvalue: s}

	return ret, nil
}

// NewIntCondition returns Condition c of int.
//  c must be CaseEq, CaseNe CaseGt, CaseGe, CaseLt or CaseLe.
func NewIntCondition(c int, i int) (*Condition, error) {
	if c != CaseEq && c != CaseNe && c != CaseGt && c != CaseGe && c != CaseLt && c != CaseLe {
		return nil, ErrInvalidCondition
	}
	ret := &Condition{ctype: types.Int, ccase: c, cvalue: i}

	return ret, nil
}

// NewUintCondition returns Condition c of uint.
//  c must be CaseEq, CaseNe CaseGt, CaseGe, CaseLt or CaseLe.
func NewUintCondition(c int, i uint) (*Condition, error) {
	if c != CaseEq && c != CaseNe && c != CaseGt && c != CaseGe && c != CaseLt && c != CaseLe {
		return nil, ErrInvalidCondition
	}
	ret := &Condition{ctype: types.Uint, ccase: c, cvalue: i}

	return ret, nil
}

// NewDoubleCondition returns Condition c of double
//  c must be CaseEq, CaseNe CaseGt, CaseGe, CaseLt or CaseLe.
func NewDoubleCondition(c int, d float64) (*Condition, error) {
	if c != CaseEq && c != CaseNe && c != CaseGt && c != CaseGe && c != CaseLt && c != CaseLe {
		return nil, ErrInvalidCondition
	}
	ret := &Condition{ctype: types.Float64, ccase: c, cvalue: d}
	return ret, nil
}

// Compare compares c and ic.
func (c Condition) Compare(ic Condition) bool {
	if c.ccase != ic.ccase || c.ctype != ic.ctype || c.cvalue == nil || ic.cvalue == nil {
		return false
	}
	switch c.ctype {
	case types.Uint:
		return c.cvalue.(uint) == ic.cvalue.(uint)
	case types.Int:
		return c.cvalue.(int) == ic.cvalue.(int)
	case types.Float64:
		return c.cvalue.(float64) == ic.cvalue.(float64)
	case types.String:
		return c.cvalue.(string) == ic.cvalue.(string)
	case types.Bool:
		return c.cvalue.(bool) == ic.cvalue.(bool)
	}
	return false
}

func (cnf *Config) SetTypeCondition(c *ConfigLine, t types.BasicKind) error {
	if c == nil {
		return errors.New("ConfigLine is nil")
	}
	k, err := convertKeys(c.ClKey)
	if err != nil {
		return fmt.Errorf("SetExists:%w", err)
	}
	tc := &TypeCondition{Keys: *k}
	cnd := &Condition{}
	switch t {
	case types.Uint:
		jn, ok := c.ClValue.(float64)
		if !ok {
			return fmt.Errorf("json number convert error. type=%T", c.ClValue)
		}
		i := uint(jn)
		cnd, err = NewUintCondition(Str2IntCase(c.ClCondition), i)
		if err != nil {
			return fmt.Errorf("NewUintCondition err:%s", err)
		}

	case types.Int:
		jn, ok := c.ClValue.(float64)
		if !ok {
			return errors.New("json number convert error")
		}
		i := int(jn)
		cnd, err = NewIntCondition(Str2IntCase(c.ClCondition), i)
		if err != nil {
			return fmt.Errorf("NewIntCondition err:%s", err)
		}

	case types.Float64:
		jn, ok := c.ClValue.(float64)
		if !ok {
			return errors.New("json number convert error")
		}
		cnd, err = NewDoubleCondition(Str2IntCase(c.ClCondition), jn)
		if err != nil {
			return fmt.Errorf("NewDoubleCondition err:%s", err)
		}
	case types.String:
		s, ok := c.ClValue.(string)
		if !ok {
			return errors.New("json string convert error")
		}
		cnd, err = NewStringCondition(Str2IntCase(c.ClCondition), s)
		if err != nil {
			return fmt.Errorf("NewStringCondition err:%s", err)
		}
	case types.Bool:
		b, ok := c.ClValue.(bool)
		if !ok {
			return fmt.Errorf("json bool convert error type=%T", c.ClValue)
		}
		cnd, err = NewBoolCondition(Str2IntCase(c.ClCondition), b)
		if err != nil {
			return fmt.Errorf("NewBoolCondition err:%s", err)
		}

	default:
		return errors.New("Invalid type")
	}
	tc.Condition = *cnd
	tc.TypeConditionStr = tc.String()
	cnf.TypeConditions = append(cnf.TypeConditions, *tc)
	return nil
}

func (tc TypeCondition) String() string {
	return fmt.Sprintf("%s %s", tc.Keys.String(), tc.Condition.String())

}

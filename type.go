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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
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
	Keys      Keys
	Condition Condition
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
	case types.Int:
		i, ok := v.(int)
		if ok {
			return c.matchInt(i), nil
		}
	case types.Uint:
		i, ok := v.(uint)
		if ok {
			return c.matchUint(i), nil
		}
	case types.Float64:
		d, ok := v.(float64)
		if ok {
			return c.matchDouble(d), nil
		}
	}
	return false, errors.New("can not cast")
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
	}
	return false
}

func (cnf *Config) SetTypeCondition(c *ConfigLine, t types.BasicKind) error {
	if c == nil {
		return errors.New("ConfigLine is nil")
	}
	k, err := convertKeys(c.Key)
	if err != nil {
		return fmt.Errorf("SetExists:%w", err)
	}
	tc := &TypeCondition{Keys: *k}
	cnd := &Condition{ctype: t, ccase: Str2IntCase(c.Condition)}
	switch t {
	case types.Uint:
		jn, ok := c.Value.(json.Number)
		if !ok {
			return errors.New("json number convert error")
		}
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf("json.Number.Int64() err:%s", err)
		}
		cnd.cvalue = uint(i)
	case types.Int:
		jn, ok := c.Value.(json.Number)
		if !ok {
			return errors.New("json number convert error")
		}
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf("json.Number.Int64() err:%s", err)
		}
		cnd.cvalue = int(i)
	case types.Float64:
		jn, ok := c.Value.(json.Number)
		if !ok {
			return errors.New("json number convert error")
		}
		i, err := jn.Float64()
		if err != nil {
			return fmt.Errorf("json.Number.Int64() err:%s", err)
		}
		cnd.cvalue = i
	case types.String:
		s, ok := c.Value.(string)
		if !ok {
			return errors.New("json string convert error")
		}
		cnd.cvalue = s
	default:
		return errors.New("Invalid type")
	}
	tc.Condition = *cnd
	cnf.TypeConditions = append(cnf.TypeConditions, *tc)
	return nil
}

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
	"errors"
	"go/types"
	"strings"
)

var ErrInvalidCondition = errors.New("Invalid condition")

const ConfigBoolKeyName = "key_bool"
const ConfigStrKeyName = "key_str"

type Condition struct {
	ctype  types.BasicKind
	ccase  int
	cvalue interface{}
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

type TypeCondition struct {
	Keys      Keys
	Condition Condition
}

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

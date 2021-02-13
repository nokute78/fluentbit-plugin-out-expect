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
	"testing"
)

func TestStr2IntCase(t *testing.T) {
	type testcase struct {
		name   string
		input  string
		expect int
	}

	cases := []testcase{
		{"gt", ">", CaseGt},
		{"ge", ">=", CaseGe},
		{"lt", "<", CaseLt},
		{"le", "<=", CaseLe},
		{"eq", "==", CaseEq},
		{"ne", "!=", CaseNe},
		{"invalid", "<=>", CaseInvalid},
	}

	for i, v := range cases {
		ret := Str2IntCase(v.input)
		if ret != v.expect {
			t.Errorf("%d:%s mismatch:\n given :%d\n expect:%d", i, v.name, ret, v.expect)
		}
	}
}

func TestNewBoolCondition(t *testing.T) {
	type testcase struct {
		name      string
		inputCase int
		inputVal  bool
	}

	okCases := []testcase{
		{"eq true", CaseEq, true},
		{"eq false", CaseEq, false},
		{"ne true", CaseNe, true},
		{"ne false", CaseNe, false},
	}

	for i, v := range okCases {
		_, err := NewBoolCondition(v.inputCase, v.inputVal)
		if err != nil {
			t.Errorf("%d:%s error=%s", i, v.name, err)
		}
	}
	ngCases := []testcase{
		{"ge", CaseGe, true},
		{"gt", CaseGt, false},
		{"le", CaseLe, true},
		{"lt", CaseLt, false},
	}
	for i, v := range ngCases {
		_, err := NewBoolCondition(v.inputCase, v.inputVal)
		if err == nil {
			t.Errorf("%d:%s should be error", i, v.name)
		}
	}

}

func TestNewStringCondition(t *testing.T) {
	type testcase struct {
		name      string
		inputCase int
		inputVal  string
	}

	okCases := []testcase{
		{"eq", CaseEq, "hoge"},
		{"ne", CaseNe, "hoge"},
		{"cont", CaseContains, "hoge"},
		{"not cont", CaseNotContains, "hoge"},
	}

	for i, v := range okCases {
		_, err := NewStringCondition(v.inputCase, v.inputVal)
		if err != nil {
			t.Errorf("%d:%s error=%s", i, v.name, err)
		}
	}
	ngCases := []testcase{
		{"ge", CaseGe, "hoge"},
		{"gt", CaseGt, "hoge"},
		{"le", CaseLe, "hoge"},
		{"lt", CaseLt, "hoge"},
	}
	for i, v := range ngCases {
		_, err := NewStringCondition(v.inputCase, v.inputVal)
		if err == nil {
			t.Errorf("%d:%s should be error", i, v.name)
		}
	}

}

func testMatch(t *testing.T, c *Condition, v interface{}, e bool) {
	t.Helper()

	b, err := c.IsMatch(v)
	if err != nil {
		t.Errorf("IsMatch err:%s", err)
	} else if b != e {
		t.Errorf("ret mismatch\n given :%t\n expect:%t", b, e)
	}

}

func TestMatchBool(t *testing.T) {
	// true case
	c, err := NewBoolCondition(CaseEq, true)
	if err != nil {
		t.Fatalf("NewBoolCondition err:%s", err)
	}
	testMatch(t, c, true, true)
	testMatch(t, c, false, false)

	// false case
	c, err = NewBoolCondition(CaseEq, false)
	if err != nil {
		t.Fatalf("NewBoolCondition err:%s", err)
	}
	testMatch(t, c, true, false)
	testMatch(t, c, false, true)
}

func TestMatchString(t *testing.T) {
	// true case
	c, err := NewStringCondition(CaseEq, "one two three")
	if err != nil {
		t.Fatalf("NewStringCondition err:%s", err)
	}
	testMatch(t, c, "one two three", true)
	testMatch(t, c, "one", false)

	c, err = NewStringCondition(CaseNe, "one two three")
	if err != nil {
		t.Fatalf("NewStringCondition err:%s", err)
	}
	testMatch(t, c, "one two three", false)
	testMatch(t, c, "one", true)

	c, err = NewStringCondition(CaseContains, "one")
	if err != nil {
		t.Fatalf("NewStringCondition err:%s", err)
	}
	testMatch(t, c, "one two three", true)
	testMatch(t, c, "four", false)

	c, err = NewStringCondition(CaseNotContains, "one")
	if err != nil {
		t.Fatalf("NewStringCondition err:%s", err)
	}
	testMatch(t, c, "one two three", false)
	testMatch(t, c, "four", true)

}

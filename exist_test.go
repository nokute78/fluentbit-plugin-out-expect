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

func TestGetValueFromMap(t *testing.T) {
	type testcase struct {
		name string
		k    Keys
		m    map[interface{}]interface{}
		ok   bool
	}
	cases := []testcase{
		testcase{"normal case", Keys{Keys: []string{"key"}}, map[interface{}]interface{}{"key": "hoge"}, true},
		testcase{"missing", Keys{Keys: []string{"key"}}, map[interface{}]interface{}{"a": "hoge"}, false},
		testcase{"nest", Keys{Keys: []string{"key", "nest"}}, map[interface{}]interface{}{"key": map[interface{}]interface{}{"nest": "hoge"}}, true},
		testcase{"too nest input", Keys{Keys: []string{"key", "nest", "errornest"}}, map[interface{}]interface{}{"key": map[interface{}]interface{}{"nest": "hoge"}}, false},
	}

	for i, v := range cases {
		_, ok := v.k.GetValueFromMap(v.m)
		if ok != v.ok {
			t.Errorf("case(%d):%s mismatch\n given :%t\n expect:%t", i, v.name, ok, v.ok)
		}
	}

}

func TestSetExistSingleKey(t *testing.T) {
	cnf := &Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":"keytest"}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}

	err = cnf.SetExist(cnfl, true)
	if err != nil {
		t.Errorf("set key error:%s", err)
	}

	if len(cnf.Exists[0].Keys) != 1 {
		t.Errorf("length is not 1")
	} else if cnf.Exists[0].Keys[0] != "keytest" {
		t.Errorf("includes mismatch:\n given :%s\n expect:%s", cnf.Exists[0], "keytest")
	}

	cnfl, err = NewConfigLineFromJson(`{"key":""}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, true)
	if err == nil {
		t.Errorf("blank string should be error")
	}
}

func TestSetExistKeyTwice(t *testing.T) {
	cnf := &Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":"keytest"}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, true)
	if err != nil {
		t.Errorf("set key error:%s", err)
	}

	err = cnf.SetExist(cnfl, true)
	if err == nil {
		t.Errorf("set key twice should be errror")
	}
}

func TestSetExistKeys(t *testing.T) {
	cnf := &Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":["key1","key2","key3"]}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, true)
	if err != nil {
		t.Errorf("set keys error:%s", err)
	}

	expect := []string{"key1", "key2", "key3"}
	if len(cnf.Exists[0].Keys) != len(expect) {
		t.Fatalf("length mismatch:\n given :%d\n expect:%d", len(cnf.Exists[0].Keys), len(expect))
	}
	for i, v := range cnf.Exists[0].Keys {
		if v != expect[i] {
			t.Errorf("mismatch(%d):\n given :%s\n expect:%s", i, v, expect[i])
		}
	}
}

func TestSetExistKeysTwice(t *testing.T) {
	cnf := &Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":["key1","key2","key3"]}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, true)
	if err != nil {
		t.Errorf("set keys error:%s", err)
	}
	err = cnf.SetExist(cnfl, true)
	if err == nil {
		t.Errorf("set keys twice should be errror")
	}
}

func TestSetNotExistSingleKey(t *testing.T) {
	cnf := &Config{}
	expect := "keytest"
	cnfl, err := NewConfigLineFromJson(`{"key":"keytest"}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, false)
	if err != nil {
		t.Errorf("set key error:%s", err)
	}

	if len(cnf.NotExists[0].Keys) != 1 {
		t.Errorf("length is not 1")
	} else if cnf.NotExists[0].Keys[0] != expect {
		t.Errorf("exclude mismatch:\n given :%s\n expect:%s", cnf.NotExists[0].Keys[0], expect)
	}

	cnfl, err = NewConfigLineFromJson(`{"key":""}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, false)
	if err == nil {
		t.Errorf("blank string should be error")
	}
}

func TestSetNotExistKeyTwice(t *testing.T) {
	cnf := &Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":"keytest"}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, false)
	if err != nil {
		t.Errorf("set key error:%s", err)
	}

	err = cnf.SetExist(cnfl, false)
	if err == nil {
		t.Errorf("set key twice should be errror")
	}
}

func TestSetNotExistKeys(t *testing.T) {
	cnf := &Config{}
	expect := []string{"key1", "key2", "key3"}
	cnfl, err := NewConfigLineFromJson(`{"key":["key1","key2","key3"]}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, false)
	if err != nil {
		t.Errorf("set keys error:%s", err)
	}

	if len(cnf.NotExists[0].Keys) != len(expect) {
		t.Fatalf("length mismatch:\n given :%d\n expect:%d", len(cnf.NotExists[0].Keys), len(expect))
	}
	for i, v := range cnf.NotExists[0].Keys {
		if v != expect[i] {
			t.Errorf("mismatch(%d):\n given :%s\n expect:%s", i, v, expect[i])
		}
	}
}

func TestSetNotExistKeysTwice(t *testing.T) {
	cnf := &Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":["key1","key2","key3"]}`)
	if err != nil {
		t.Errorf("NewConfigLineFromJson err:%s", err)
	}
	err = cnf.SetExist(cnfl, false)
	if err != nil {
		t.Errorf("set keys error:%s", err)
	}
	err = cnf.SetExist(cnfl, false)
	if err == nil {
		t.Errorf("set keys twice should be errror")
	}
}

func TestKeyCompare(t *testing.T) {
	type testcase struct {
		name   string
		src    Keys
		dst    Keys
		expect bool
	}

	cases := []testcase{
		testcase{"normal", Keys{Keys: []string{"key"}}, Keys{Keys: []string{"key"}}, true},
		testcase{"src > dst", Keys{Keys: []string{"key", "key2"}}, Keys{Keys: []string{"key"}}, false},
		testcase{"src < dst", Keys{Keys: []string{"key"}}, Keys{Keys: []string{"key", "key2"}}, false},
		testcase{"blanks", Keys{Keys: []string{}}, Keys{Keys: []string{}}, false},
		testcase{"blank", Keys{Keys: []string{""}}, Keys{Keys: []string{}}, false},
	}

	for i, v := range cases {
		ret := v.src.Compare(v.dst)
		if ret != v.expect {
			t.Errorf("%d:%s mismatch\n given :%t\n expect:%t", i, v.name, ret, v.expect)
		}
	}

}

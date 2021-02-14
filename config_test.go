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

func TestConflict(t *testing.T) {
	cnf := Config{}
	cnfl, err := NewConfigLineFromJson(`{"key":["key","key2"]}`)
	if err != nil {
		t.Fatalf("NewConfigLine error:%s", err)
	}

	err = cnf.SetExist(cnfl, true)
	if err != nil {
		t.Errorf("set exist(true) error:%s %+v", err, cnfl)
	}

	cnfl2, err := NewConfigLineFromJson(`{"key":"key"}`)
	if err != nil {
		t.Fatalf("NewConfigLine error:%s", err)
	}
	err = cnf.SetExist(cnfl2, true)
	if err != nil {
		t.Errorf("set key(false) error:%s", err)
	}

	err = cnf.Validate()
	if err != nil {
		t.Errorf("Validate error: %s", err)
	}

	err = cnf.SetExist(cnfl, false)
	if err != nil {
		t.Errorf("set key(false) error:%s", err)
	}

	err = cnf.Validate()
	if err == nil {
		t.Errorf("It should be conflict")
	}

}

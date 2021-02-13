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
	"strings"
)

type Keys struct {
	Keys        []string
	FlattenKeys string
}

// String implements fmt.Stringer.
func (k Keys) String() string {
	if len(k.Keys) == 0 {
		return ""
	} else if len(k.Keys) == 1 {
		return `"` + k.Keys[0] + `"`
	}
	ret := `"` + k.Keys[0] + `"`
	for i := 1; i < len(k.Keys); i++ {
		ret = ret + "->" + `"` + k.Keys[i] + `"`
	}

	return ret
}

// GetValueFromMap returns the value from map m.
//   If the value is not found, it returns nil, false.
func (k Keys) GetValueFromMap(m map[interface{}]interface{}) (interface{}, bool) {
	if len(k.Keys) == 0 || m == nil {
		return nil, false
	}
	var ret interface{}
	ret = m

	for _, key := range k.Keys {
		im, ok := ret.(map[interface{}]interface{})
		if !ok {
			return nil, false
		}
		ret, ok = im[key]
		if !ok {
			return nil, false
		}
	}
	return ret, true
}

// Compare compares k and ks.
func (k Keys) Compare(ks Keys) bool {
	if len(k.Keys) != len(ks.Keys) || len(k.Keys) == 0 {
		return false
	}
	for i, key := range k.Keys {
		if ks.Keys[i] != key {
			return false
		}
	}
	return true
}

func containsKeys(keyss []Keys, ks Keys) bool {
	if len(ks.Keys) == 0 {
		return false
	}

	for _, keys := range keyss {
		if keys.Compare(ks) {
			return true
		}
	}
	return false
}

// IsExistKeys check if ExistKeys of cnf has ks keys or not.
func (cnf *Config) IsExistKeys(ks Keys) bool {
	return containsKeys(cnf.Exists, ks)
}

// IsNotExistKeys check if NotExistKeys of cnf has ks keys or not.
func (cnf *Config) IsNotExistKeys(ks Keys) bool {
	return containsKeys(cnf.NotExists, ks)
}

// SetKey parses input param string and set member variable.
//   If isExist is true, param is treated as ExistKey.
//   If isExist is false, param is treated as NotExistKey.
func (cnf *Config) SetKey(param string, isExist bool) error {
	if len(param) == 0 {
		return errors.New("string is blank")
	}
	strs := []string{}
	if strings.Contains(param, " ") {
		tmp := strings.Split(param, " ")
		for _, v := range tmp {
			if len(v) > 0 {
				strs = append(strs, v)
			}
		}
	} else {
		strs = append(strs, param)
	}

	if len(strs) == 0 {
		return errors.New("config not found")
	}

	k := Keys{Keys: strs}
	if isExist {
		if cnf.IsExistKeys(k) {
			return errors.New("already exist")
		}
		k.FlattenKeys = k.String()
		cnf.Exists = append(cnf.Exists, k)
	} else {
		if cnf.IsNotExistKeys(k) {
			return errors.New("already exist")
		}
		k.FlattenKeys = k.String()
		cnf.NotExists = append(cnf.NotExists, k)
	}

	return nil
}

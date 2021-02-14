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
	"fmt"
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

// SetKey set member variable via c.
//   If isExist is true, c is treated as ExistKey.
//   If isExist is false, c is treated as NotExistKey.
func (cnf *Config) SetExist(c *ConfigLine, isExist bool) error {
	if c == nil {
		return errors.New("ConfigLine is nil")
	}
	k, err := convertKeys(c.ClKey)
	if err != nil {
		return fmt.Errorf("SetExists:%w", err)
	}

	if isExist {
		if cnf.HasExistKeys(k) {
			return errors.New("already exist")
		}
		cnf.Exists = append(cnf.Exists, *k)
	} else {
		if cnf.HasNotExistKeys(k) {
			return errors.New("already exist")
		}
		cnf.NotExists = append(cnf.NotExists, *k)
	}
	return nil
}

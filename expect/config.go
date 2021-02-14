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
	"encoding/json"
	"errors"
)

const ParamNumMax = 16
const ConfigExistKeyName = "key_exists"
const ConfigNotExistKeyName = "key_not_exists"

// Config represents context of this plugin.
type Config struct {
	Exists         []Keys
	NotExists      []Keys
	TypeConditions []TypeCondition
}

// Validate check if configuration value is ok or not.
func (c Config) Validate() error {
	for _, key := range c.Exists {
		for _, vkey := range c.NotExists {
			if key.Compare(vkey) {
				return errors.New("conflict key")
			}
		}
	}
	return nil
}

// ConfigLine represents each line of config file.
type ConfigLine struct {
	ClKey       interface{} `json:"key"` // string or []string
	ClValue     interface{} `json:"value,omitempty"`
	ClCondition string      `json:"condition,omitempty"`
}

// NewConfigLineFromJson returns ConfigLine pointer via Json s.
func NewConfigLineFromJson(s string) (*ConfigLine, error) {
	ret := &ConfigLine{}
	err := json.Unmarshal([]byte(s), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// v (string or []string) -> *Keys
func convertKeys(v interface{}) (*Keys, error) {
	var ss []string
	s, ok := v.(string)
	if ok {
		if s == "" {
			return nil, errors.New("blank string")
		}
		ss = append(ss, s)
	} else {
		ia, ok := v.([]interface{})
		if !ok {
			return nil, errors.New("cannot convert key array")
		}
		ss = make([]string, len(ia))
		for i, vv := range ia {
			ss[i], ok = vv.(string)
			if !ok {
				return nil, errors.New("cannot convert key string")
			} else if ss[i] == "" {
				return nil, errors.New("blank string")
			}
		}
	}
	ret := &Keys{Keys: ss}
	ret.FlattenKeys = ret.String()

	return ret, nil
}

func containsKeys(keyss []Keys, ks *Keys) bool {
	if ks == nil || len(ks.Keys) == 0 {
		return false
	}

	for _, keys := range keyss {
		if keys.Compare(*ks) {
			return true
		}
	}
	return false
}

// HasExistKeys check if ExistKeys of cnf has ks keys or not.
func (cnf *Config) HasExistKeys(ks *Keys) bool {
	return containsKeys(cnf.Exists, ks)
}

// HasNotExistKeys check if NotExistKeys of cnf has ks keys or not.
func (cnf *Config) HasNotExistKeys(ks *Keys) bool {
	return containsKeys(cnf.NotExists, ks)
}

// HasTypeCondition check if TypeConditions of cnf has t or not.
func (cnf *Config) HasTypeCondition(t *TypeCondition) bool {
	if t == nil || len(t.Keys.Keys) == 0 {
		return false
	}
	for _, tc := range cnf.TypeConditions {
		if tc.Keys.Compare(t.Keys) && tc.Condition.Compare(t.Condition) {
			return true
		}
	}
	return false
}

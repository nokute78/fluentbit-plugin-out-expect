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
)

const ParamNumMax = 16
const ConfigExistKeyName = "key_exists"
const ConfigNotExistKeyName = "key_not_exists"

type Config struct {
	Exists    []Keys
	NotExists []Keys
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

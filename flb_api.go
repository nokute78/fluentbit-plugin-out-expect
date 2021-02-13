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
	"C"
	"errors"
	"log"
	"strconv"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)

func getParameter(p unsafe.Pointer, key string, i int) (string, error) {
	s := key + strconv.Itoa(i)
	param := output.FLBPluginConfigKey(p, s)
	if len(param) == 0 {
		return "", errors.New("Not found")
	}
	return param, nil
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "gexpect", "Check if a key/value is expected the key/value")
}

//export FLBPluginInit
// (fluentbit will call this)
// plugin (context) pointer to fluentbit context (state/ c code)
func FLBPluginInit(p unsafe.Pointer) int {
	cnf := Config{}

	// Exist
	for i := 0; i < ParamNumMax; i++ {
		param, err := getParameter(p, ConfigExistKeyName, i)
		if err == nil {
			cnf.SetKey(param, true)
		}
		param, err = getParameter(p, ConfigNotExistKeyName, i)
		if err == nil {
			cnf.SetKey(param, false)
		}
	}

	output.FLBPluginSetContext(p, cnf)
	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	log.Println("[expect] Flush called for unknown instance")
	return output.FLB_OK
}

func reportsErrors(reports []string, tag *C.char) {
	log.Println(strconv.Itoa(len(reports)) + " error(s) detected! tag:" + C.GoString(tag))
	for _, v := range reports {
		log.Println(" " + v)
	}
	log.Println("")
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx unsafe.Pointer, data unsafe.Pointer, length C.int, tag *C.char) int {
	cnf, ok := output.FLBPluginGetContext(ctx).(Config)
	if !ok {
		log.Println("[expect] Context Conversion error")
		return output.FLB_ERROR
	}

	dec := output.NewDecoder(data, int(length))

	for {
		reports := []string{}
		ret, _, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		for _, keys := range cnf.Exists {
			_, ok := keys.GetValueFromMap(record)
			if !ok {
				reports = append(reports, "Exist key not found:"+keys.FlattenKeys)
			}
		}
		for _, keys := range cnf.NotExists {
			_, ok := keys.GetValueFromMap(record)
			if ok {
				reports = append(reports, "Not Exist key found:"+keys.FlattenKeys)
			}
		}

		if len(reports) > 0 {
			reportsErrors(reports, tag)
		}

	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

// dummy
func main() {
}

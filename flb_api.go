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
	"go/types"
	"log"
	"strconv"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/nokute78/fluentbit-plugin-out-expect/expect"
)

const Version = "0.0.1"

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
	cnf := expect.Config{}
	log.Printf("[expect] Ver: %s\n", Version)

	// Exist
	for i := 0; i < expect.ParamNumMax; i++ {
		param, err := getParameter(p, expect.ConfigExistKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			cnf.SetExist(p, true)
		}
		param, err = getParameter(p, expect.ConfigNotExistKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			cnf.SetExist(p, false)
		}
		param, err = getParameter(p, expect.ConfigBoolKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			err = cnf.SetTypeCondition(p, types.Bool)
			if err != nil {
				log.Printf("bool config error=%s\n", err)
			}
		}
		param, err = getParameter(p, expect.ConfigStrKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			err = cnf.SetTypeCondition(p, types.String)
			if err != nil {
				log.Printf("string config error=%s\n", err)
			}
		}
		param, err = getParameter(p, expect.ConfigIntKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			err = cnf.SetTypeCondition(p, types.Int)
			if err != nil {
				log.Printf("int config error=%s\n", err)
			}
		}
		param, err = getParameter(p, expect.ConfigUintKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			err = cnf.SetTypeCondition(p, types.Uint)
			if err != nil {
				log.Printf("uint config error=%s\n", err)
			}
		}
		param, err = getParameter(p, expect.ConfigDoubleKeyName, i)
		if err == nil {
			p, err := expect.NewConfigLineFromJson(param)
			if err != nil {
				continue
			}
			err = cnf.SetTypeCondition(p, types.Float64)
			if err != nil {
				log.Printf("double config error=%s\n", err)
			}
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
	cnf, ok := output.FLBPluginGetContext(ctx).(expect.Config)
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
		for _, tc := range cnf.TypeConditions {
			v, ok := tc.Keys.GetValueFromMap(record)
			if !ok {
				reports = append(reports, "Key not found:"+tc.Keys.FlattenKeys)
				continue
			}
			b, err := tc.Condition.IsMatch(v)
			if err != nil {
				reports = append(reports, "IsMatch error:"+tc.Keys.FlattenKeys)
			} else if !b {
				reports = append(reports, "Error. expect: value "+i2str(v)+" of "+tc.TypeConditionStr)
			}
		}

		if len(reports) > 0 {
			reportsErrors(reports, tag)
		}

	}

	return output.FLB_OK
}

func i2str(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case bool:
		b := v.(bool)
		return strconv.FormatBool(b)
	case int:
		i := v.(int)
		return strconv.FormatInt(int64(i), 10)
	case int8:
		i := v.(int8)
		return strconv.FormatInt(int64(i), 10)
	case int16:
		i := v.(int16)
		return strconv.FormatInt(int64(i), 10)
	case int32:
		i := v.(int32)
		return strconv.FormatInt(int64(i), 10)
	case int64:
		i := v.(int64)
		return strconv.FormatInt(i, 10)
	case uint:
		i := v.(uint)
		return strconv.FormatUint(uint64(i), 10)
	case uint8:
		i := v.(uint8)
		return strconv.FormatUint(uint64(i), 10)
	case uint16:
		i := v.(uint16)
		return strconv.FormatUint(uint64(i), 10)
	case uint32:
		i := v.(uint32)
		return strconv.FormatUint(uint64(i), 10)
	case uint64:
		i := v.(uint64)
		return strconv.FormatUint(i, 10)
	case float32:
		f := v.(float32)
		return strconv.FormatFloat(float64(f), 'f', -1, 32)
	case float64:
		f := v.(float64)
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return ""
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

// dummy
func main() {
}

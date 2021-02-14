# fluent-bit-plugin-out-expect

Ouput plugin for [Fluent-Bit](https://fluentbit.io/) to verify key/value of record.

## Feature

* Conditional checking
* Nested field support

## Configuration Parameters

Each configuration name should be *key_nameN*. *N* is 0-15.

### Key Exists
*key_existsN* *Json Object*
or
*key_not_existsN* *Json Object*

Json object:
|Key|Value Type|Description|
|---|----------|-----------|
|`"key"`|string or string array|The key name to check if it exists or not. If it is array, it is recognized as nested keys.|

Example:
|use case| example configuration|
|--------|----------------------|
|Key "alert" should be exist |`key_exist0 {"key":"alert"}` |
|Key "alert" should not be exist |`key_not_exist0 {"key":"alert"}` |

### Boolean
*key_boolN* *Json Object*

Json object:
|Key|Value Type|Description|
|---|----------|-----------|
|`"key"`      |string or string array|The key name to check if it exists or not. If it is array, it is recognized as nested keys.|
|`"value"`    |boolean|Checking value.|
|`"condition"`|string|Checking Condition. `"=="`/`"!="|

Example:
|use case| example configuration|
|--------|----------------------|
|Value of key "not_nil" should be true |`key_bool0 {"key":"not_nil","condition","==", "value":true}` |
|Value of key "not_nil" should be false|`key_bool0 {"key":"not_nil","condition","contains", "value":false}` |

### String
*key_strN* *Json Object*

Json object:
|Key|Value Type|Description|
|---|----------|-----------|
|`"key"`      |string or string array|The key name to check if it exists or not. If it is array, it is recognized as nested keys.|
|`"value"`    |string|Checking value.|
|`"condition"`|string|Checking Condition. `"=="`/`"!="`/`"contains"`/`"not_contains"`|

Example:
|use case| example configuration|
|--------|----------------------|
|Value of key "name" should be match "Taro"|`key_str0 {"key":"name","condition","==", "value":"taro"}` |
|Value of key "name" should be contain "Taro"|`key_str0 {"key":"name","condition","contains", "value":"taro"}` |

### Int
*key_intN* *Json Object*

Json object:
|Key|Value Type|Description|
|---|----------|-----------|
|`"key"`      |string or string array|The key name to check if it exists or not. If it is array, it is recognized as nested keys.|
|`"value"`    |int|Checking value.|
|`"condition"`|string|Checking Condition. `"=="`/`"!="`/`">"`/`">="`/`"<"`/`"<="`|

Example:
|use case| example configuration|
|--------|----------------------|
|Value of key "log_level" should be match 3|`key_int0 {"key":"log_level","condition","==", "value":3}` |
|Value of key "log_level" should be greater than 3|`key_int0 {"key":"log_level","condition",">", "value":3}` |

### Uint
*key_uintN* *Json Object*

Json object:
|Key|Value Type|Description|
|---|----------|-----------|
|`"key"`      |string or string array|The key name to check if it exists or not. If it is array, it is recognized as nested keys.|
|`"value"`    |uint|Checking value.|
|`"condition"`|string|Checking Condition. `"=="`/`"!="`/`">"`/`">="`/`"<"`/`"<="`|

Example:
|use case| example configuration|
|--------|----------------------|
|Value of key "log_level" should be match 3|`key_uint0 {"key":"log_level","condition","==", "value":3}` |
|Value of key "log_level" should be greater than 3|`key_uint0 {"key":"log_level","condition",">", "value":3}` |

### Double
*key_doubleN* *Json Object*

Json object:
|Key|Value Type|Description|
|---|----------|-----------|
|`"key"`      |string or string array|The key name to check if it exists or not. If it is array, it is recognized as nested keys.|
|`"value"`    |double|Checking value.|
|`"condition"`|string|Checking Condition. `"=="`/`"!="`/`">"`/`">="`/`"<"`/`"<="`|

Example:
|use case| example configuration|
|--------|----------------------|
|Value of key "degree" should be match 27.3|`key_double0 {"key":"degree","condition","==", "value":27.3}` |
|Value of key "degree" should be greater than 27.3|`key_double0 {"key":"degree","condition",">", "value":27.3}` |


## Build

```
make
```

## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)
# fluent-bit-plugin-out-expect

Ouput plugin for [Fluent-Bit](https://fluentbit.io/) to verify key/value of record.

## Feature

* Nested field support

## Configuration Parameters

|Check if|Format|Description|
|-------|------|-----------|
|A key exists       | *key_existsN* *key1* [*key2*] ...     | |
|A key doesn't exist| *key_not_existsN* *key1* [*key2*] ... | |

## Examples

## Build

```
make
```

## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)
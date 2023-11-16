# UDP Server API

This API provides two methods: `get_result` and `help`.

## `get_result {mode: str}`

This method returns the result of performance test in specified mode. The `mode` parameter specifies the type of encoding to use, which must be one of the following values:

- `avro`: Use Apache Avro encoding.
- `json`: Use JSON encoding.
- `msgpack`: Use MessagePack encoding.
- `protobuf`: Use Protocol Buffers encoding.
- `python_native`: Use Python's native data structures.
- `xml`: Use XML encoding.

If the `mode` parameter is not one of the above values, the method will return an error.

### Response

The response will contain the result of performance test in specified mode.
Response has following format:
```
{mode} - {byte size of encoded test object} - {time of encoding test object 1000 times}ms - {time of decoding encoded test object 1000 times}ms\n
```

### Example
```bash
echo 'get_result avro' | nc -u localhost 2000 -w 3
# avro - 221 - 341ms - 191ms
```

## `help`
_aliases: `?`, `info`_

This method returns a description of the API.

### Response

The response will contain a short description of the available API methods.

### Example
```bash
echo 'help' | nc -u localhost 2000 -w 3
# get_result {mode}

# Available modes:
# avro
# json
# msgpack
# protobuf
# python_native
# xml
```

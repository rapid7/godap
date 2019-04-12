# GODAP: The Data Analysis Pipeline
(a port of the ruby-based DAP: https://github.com/rapid7/dap)

[![Gem Version](https://badge.fury.io/rb/godap.svg)](http://badge.fury.io/rb/godap)
[![Build Status](https://travis-ci.org/rapid7/godap.svg?branch=master)](https://travis-ci.org/rapid7/godap)

DAP was created to transform text-based data on the command-line, specializing in transforms that are annoying or difficult to do with existing tools.

DAP reads data using an input plugin, transforms it through a series of filters, and prints it out again using an output plugin. Every record is treated as a document (aka: hash/dict) and filters are used to reduce, expand, and transform these documents as they pass through. Think of DAP as a mashup between sed, awk, grep, csvtool, and jq, with map/reduce capabilities.

DAP was written to process terabyte-sized public scan datasets, such as those provided by https://opendata.rapid7.com/. This go version of dap supports parallel processing of data. Results are forwarded to stdout and consistency of ordering is not guaranteed (and are highly likely to be out of order when compared to the input data stream).

## Usage

### Quick Setup for GeoIP Lookups

Note: The documentation below assumes you've properly setup $GOPATH and $PATH (usually $GOPATH/bin:$PATH) per the official golang documentation.

```
$ go get github.com/rapid7/godap
$ sudo bash
# mkdir -p /var/lib/geoip && cd /var/lib/geoip && wget http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz && gunzip GeoLiteCity.dat.gz && mv GeoLiteCity.dat geoip.dat
```

```
$  echo 8.8.8.8 | godap lines + geo_ip line + json
{"line":"8.8.8.8","line.country_code":"US","line.country_code3":"USA","line.country_name":"United States","line.latitude":"38.0","line.longitude":"-97.0"}
```

Where dap gets fun is doing transforms, like just grabbing the country code:
```
$  echo 8.8.8.8 | godap lines + geo_ip line + select line.country_code3 + lines
USA
```

## Inputs, filters and outputs
The general syntax when calling godap is ```godap <input> + (<filter +> <filter +> <filter +> ...) + <output>```, where input, filter and output correspond to one of the supported features below. Filters are optional, though an input and output are required. Each feature component is separated by ```+```. Component options are specified immediately after the component declaration. For example, streaming from a wifi adapter and spitting out json documents would resemble: ```godap pcap iface=en0 rfmon=true + json```. Component options with spaces or other complexities can be specified using shell-like quoting. For example, for a bpf pcap filter on the pcap component: ```godap pcap iface=en0 'filter="tcp port 80"' + json```.

## Inputs

 * <a id="pcap">pcap</a>

  Specifies that the input stream is a packet capture. Currently supports streaming in from a file or interface.

  | Option  | Description                                                                           | Value               | Default |
  |---------|---------------------------------------------------------------------------------------|---------------------|---------|
  | iface   | The interface to read packets from. If iface is specified, file must not be specified | string interface id | ```<none>```    |
  | file    | The pcap file to read from. If file is specified, iface must not be specified         | string filename     | ```<none>```    |
  | promisc | Whether to capture in promiscuous mode                                                | boolean             | false   |
  | timeout | The capture timeout                                                                   | integer             | -1 (inf)|
  | snaplen | The snap length                                                                       | integer             | 65536   |
  | rfmon   | Whether to capture in monitor mode (applicable only to adapters which support it)     | boolean             | false   |

  Example:
  * Pull packets in from monitor mode: ```godap pcap iface=en0 rfmon=true + json```
  * Read pcap (or pcap-ng) file contents and convert to json: ```godap pcap file=foo.pcap + json```
  * Live capture in promiscuous mode: ```godap pcap iface=en0 promisc=true + json```

 * <a id="json">json</a>

  Specifies that the input stream is represented as JSON data.

  | Option  | Description                                                                           | Value               | Default |
  |---------|---------------------------------------------------------------------------------------|---------------------|---------|
  | file    | The file to stream from. If not specified, stdin is assumed. Can also be - for stdin. | string filename     | stdin   |

  Example:
   ```
$  echo '{"a":2}' | godap json + lines
2
   ```

 * <a id="lines">lines</a>

  Specifies that the input stream is represented as newline-terminated plaintext.

  | Option  | Description                                                                           | Value               | Default |
  |---------|---------------------------------------------------------------------------------------|---------------------|---------|
  | file    | The file to stream from. If not specified, stdin is assumed. Can also be - for stdin. | string filename     | stdin   |

  Example:
   ```
$  echo hello world | godap lines + json
{"line":"hello world"}
   ```
## Filters

 * <a id="rename">rename</a>

  Renames a document field.

  | Option               | Description                | Value                          | Default |
  |----------------------|----------------------------|--------------------------------|---------|
  | ```<document key>``` | The document key to rename | string ```<destination key>``` | ```<none>```    |

  Example:
   ```
$  echo world | godap lines + rename line=hello + json
{"hello":"world"}
   ```
 * <a id="not_exists">not_exists</a>

  Prevents a document from further processing if a specified key is not present

  | Option               | Description                | Value                          | Default         |
  |----------------------|----------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The document key that must not exist | ```<none>```                   | ```<none>```    |

  Example:
   ```
$  echo '{"foo":"bar"}' | godap json + not_exists foo + json

   ```

 * <a id="split_comma">split_comma</a>

  Extracts comma-separated fields in the specified key's value into new documents

  | Option               | Description                | Value                          | Default         |
  |----------------------|----------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The document key which will be split | ```<none>```                   | ```<none>```    |

  Example:
   ```
$  echo '{"foo":"bar,baz"}' | godap json + split_comma foo + json
{"foo":"bar,baz","foo.word":"bar"}
{"foo":"bar,baz","foo.word":"baz"}
   ```

 * field_split_line

  Extracts fields separated by a newline from the source key's value into new fields of the same document. Each new key is named ```<origkey>.f###``` where ### is an incremental integer indicating the original field position from left to right.

  | Option               | Description                | Value                          | Default         |
  |----------------------|----------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The document key which will be split | ```<none>```                   | ```<none>```    |

  Example:
   ```
$  echo '{"foo":"bar\nbaz"}' | godap json + field_split_line foo + json
{"foo":"bar\nbaz","foo.f1":"bar","foo.f2":"baz"}
   ```

 * <a id="not_empty">not_empty</a>

  Filters out a document if the value for a given key is empty

  | Option               | Description                   | Value                          | Default         |
  |----------------------|-------------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The document key to filter on | ```<none>```                   | ```<none>```    |

  Example:
   ```
$  echo '{"foo":"bar,baz"}' | godap json + not_empty foo + json
{"foo":"bar,baz"}
   ```

 * <a id="field_split_tab">field_split_tab</a>

  Splits a key into multiple new subkeys each containing a field from the original value split by ```\t```. Each new key is named ```<origkey>.f###``` where ### is an incremental integer indicating the original field position from left to right.

  | Option               | Description                   | Value                          | Default         |
  |----------------------|-------------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The document key to split     | ```<none>```                   | ```<none>```    |

  Example:
   ```
$  echo '{"foo":"bar\tbaz"}' | godap json + field_split_tab foo + json
{"foo":"bar\tbaz","foo.f1":"bar","foo.f2":"baz"}
   ```
 * <a id="truncate">truncate</a>

  Sets the value of the specified key to the empty string

  | Option               | Description                   | Value                          | Default         |
  |----------------------|-------------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The key to truncate           | ```<none>```                   | ```<none>```    |

  Example:
   ```
$  echo '{"foo":"bar\tbaz"}' | godap json + truncate foo + json
{"foo":""}
   ```
 * <a id="flatten">flatten</a>

  Flattens pesky nested json into "key.name" properties of the top-level document. The original key/value are left
  untouched (you can remove them using the `remove` filter).

  | Option               | Description                   | Value                          | Default         |
  |----------------------|-------------------------------|--------------------------------|-----------------|
  | ```<document key>``` | The key to flatten            | ```<none>```                   | ```<none>```    |

  Example:
   ```
$ echo '{"foo":{"bar": "baz"}}' | godap json + flatten foo + json
{"foo":{"bar":"baz"},"foo.bar":"baz"}
   ```
 * <a id="insert">insert</a>

  Adds a new value to the document

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to truncate           | ```<document value>```         | ```<empty string>```    |

  Example:
   ```
$  echo '{"foo":"bar\tbaz"}' | godap json + insert a=b + json
{"a":"b","foo":"bar\tbaz"}
   ```
 * <a id="field_split_array">field_split_array</a>

  Splits a field that contains an array data type value into multiple new fields. Each new key is named ```<origkey>.f###``` where ### is an incremental integer indicating the original field position from left to right. The array can contain multiple different data types.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to split              | ```<none>```                   | ```<none>```            |

  Example:
   ```
$  echo '{"foo":["a",2]}' | godap json + field_split_array foo + json
{"foo":["a",2],"foo.f1":"a","foo.f2":2}
   ```
 * <a id="exists">exists</a>

  Ensures the specified key exists in the source document. If it does not, the document is removed from the pipeline.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key that must exist       | ```<none>```                   | ```<none>```            |

  Example:
   ```
$  echo '{"foo":"bar\tbaz"}' | godap json + exists a + json
$  echo '{"foo":"bar\tbaz"}' | godap json + exists foo + json
{"foo":"bar\tbaz"}

   ```
 * <a id="split_line">split_line</a>

  Splits a given key's value into multiple new documents with the same key name, each document containing a field extracted from the source key's value separated by a newline.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to split              | ```<none>```                   | ```<none>```            |

  Example:
   ```
$  echo '{"foo":"bar\nbaz"}' | godap json + split_line foo + json
{"foo":"bar\nbaz","foo.line":"bar"}
{"foo":"bar\nbaz","foo.line":"baz"}
   ```
 * <a id="select">select</a>

  Keeps only the specified keys in the resulting document. Multiple key names can be specified.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to keep               | ```<none>```                   | ```<none>```            |

  Example:
   ```
$  echo '{"foo":"bar", "baz":"qux", "a":"b"}' | godap json + select foo + json
{"foo":"bar"}
$  echo '{"foo":"bar", "baz":"qux", "a":"b"}' | godap json + select foo baz + json
{"baz":"qux","foo":"bar"}
   ```
 * <a id="remove">remove</a>

  Removes the specified keys from the source document.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to remove             | ```<none>```                   | ```<none>```            |

  Example:
   ```
$  echo '{"foo":"bar", "baz":"qux", "a":"b"}' | godap json + remove foo baz + json
{"a":"b"}
   ```
 * <a id="include">include</a>

  Ensures a document key includes a specified string.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to remove             | ```string contains_str```      | ```<none>```            |

  Example:
   ```
$  echo '{"foo":"bar", "baz":"qux", "a":"b"}' | godap json + include a=c + json
$  echo '{"foo":"bar", "baz":"qux", "a":"b"}' | godap json + include a=b + json
{"a":"b","baz":"qux","foo":"bar"}
   ```
 * <a id="transform">transform</a>

  Ensures a document key includes a specified string.

  | Option               | Description                   | Value                          | Default                 |
  |----------------------|-------------------------------|--------------------------------|-------------------------|
  | ```<document key>``` | The key to transform          | ```utf8encode``` or ```ascii``` or ```base64encode``` or ```base64decode``` or ```upcase``` or ```downcase``` or ```hexencode``` or ```ascii```         | ```<none>```            |


  Example:
   ```
$  echo '{"foo":"bar"}' | godap json + transform foo=base64encode + json
{"foo":"YmFy"}

   ```
 * <a id="field_array_join_whitespace">field_array_join_whitespace</a>
 * <a id="digest">digest</a>
 * <a id="geo_ip">geo_ip</a>
 * <a id="annotate">annotate</a>
 * <a id="split_word">split_word</a>
 * <a id="field_split_comma">field_split_comma</a>
 * <a id="field_array_join_comma">field_array_join_comma</a>
 * <a id="exclude">exclude</a>
 * <a id="where">where</a>
 * <a id="split_tab">split_tab</a>
 * <a id="split_array">split_array</a>
 * <a id="field_split_word">field_split_word</a>

## Outputs

 * <a id="json">json</a>
 * <a id="lines">lines</a>

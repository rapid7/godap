# GODAP: The Data Analysis Pipeline 
(a port of the ruby-based DAP: https://github.com/rapid7/dap)

DAP was created to transform text-based data on the command-line, specializing in transforms that are annoying or difficult to do with existing tools.

DAP reads data using an input plugin, transforms it through a series of filters, and prints it out again using an output plugin. Every record is treated as a document (aka: hash/dict) and filters are used to reduce, expand, and transform these documents as they pass through. Think of DAP as a mashup between sed, awk, grep, csvtool, and jq, with map/reduce capabilities.

DAP was written to process terabyte-sized public scan datasets, such as those provided by https://scans.io/. This go version of dap supports parallel processing of data. Results are forwarded to stdout and consistency of ordering is not guaranteed (and are highly likely to be out of order when compared to the input data stream).

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
$  echo 8.8.8.8 | godap + lines + geo_ip line + select line.country_code3 + lines
USA
```


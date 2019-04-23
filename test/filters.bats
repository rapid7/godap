#!/usr/bin/env bats

load ./test_common

@test "rename" {
  run bash -c 'echo world | $DAP_EXECUTABLE lines + rename line=hello + json'
  assert_success
  assert_output '{"hello":"world"}'
}

@test "not_exists" {
  run bash -c "echo '{\"foo\":\"bar\"}' | $DAP_EXECUTABLE json + not_exists foo + json"
  assert_success
  assert_output ''
  run bash -c "echo '{\"bar\":\"bar\"}' | $DAP_EXECUTABLE json + not_exists foo + json"
  assert_success
  assert_output '{"bar":"bar"}'
}

@test "split_comma" {
  run bash -c "echo '{\"foo\":\"bar,baz\"}' | $DAP_EXECUTABLE json + split_comma foo + json"
  assert_success
  assert_line --index 0 '{"foo":"bar,baz","foo.word":"bar"}'
  assert_line --index 1 '{"foo":"bar,baz","foo.word":"baz"}'
}

@test "field_split_line" {
  run bash -c "echo '{\"foo\":\"bar\nbaz\"}' | $DAP_EXECUTABLE json + field_split_line foo + json"
  assert_success
  assert_output '{"foo":"bar\nbaz","foo.f1":"bar","foo.f2":"baz"}'
}

@test "not_empty" {
  run bash -c "echo '{\"foo\":\"bar,baz\"}' | $DAP_EXECUTABLE json + not_empty foo + json"
  assert_success
  assert_output '{"foo":"bar,baz"}'
}

@test "field_split_tab" {
  run bash -c "echo '{\"foo\":\"bar\tbaz\"}' | $DAP_EXECUTABLE json + field_split_tab foo + json"
  assert_success
  assert_output '{"foo":"bar\tbaz","foo.f1":"bar","foo.f2":"baz"}'
}

@test "truncate" {
  run bash -c "echo '{\"foo\":\"bar\tbaz\"}' | $DAP_EXECUTABLE json + truncate foo + json"
  assert_success
  assert_output '{"foo":""}'
}

@test "insert" {
  run bash -c "echo '{\"foo\":\"bar\tbaz\"}' | $DAP_EXECUTABLE json + insert a=b + json"
  assert_success
  assert_output '{"a":"b","foo":"bar\tbaz"}'
}

@test "field_split_array" {
  run bash -c "echo '{\"foo\":[\"a\",2]}' | $DAP_EXECUTABLE json + field_split_array foo + json"
  assert_success
  assert_output '{"foo":["a",2],"foo.f1":"a","foo.f2":2}'
}

@test "exists" {
  run bash -c "echo '{\"foo\":\"bar\tbaz\"}' | $DAP_EXECUTABLE json + exists a + json"
  assert_success
  assert_output ''
  run bash -c "echo '{\"foo\":\"bar\tbaz\"}' | $DAP_EXECUTABLE json + exists foo + json"
  assert_success
  assert_output '{"foo":"bar\tbaz"}'
}

@test "split_line" {
  run bash -c "echo '{\"foo\":\"bar\nbaz\"}' | $DAP_EXECUTABLE json + split_line foo + json"
  assert_success
  assert_line --index 0 '{"foo":"bar\nbaz","foo.line":"bar"}'
  assert_line --index 1 '{"foo":"bar\nbaz","foo.line":"baz"}'
}

@test "select" {
  run bash -c "echo '{\"foo\":\"bar\", \"baz\":\"qux\", \"a\":\"b\"}' | $DAP_EXECUTABLE json + select foo + json"
  assert_success
  assert_output '{"foo":"bar"}'
  run bash -c "echo '{\"foo\":\"bar\", \"baz\":\"qux\", \"a\":\"b\"}' | $DAP_EXECUTABLE json + select foo baz + json"
  assert_success
  assert_output '{"baz":"qux","foo":"bar"}'
}

@test "remove" {
  run bash -c "echo '{\"foo\":\"bar\", \"baz\":\"qux\", \"a\":\"b\"}' | $DAP_EXECUTABLE json + remove foo baz + json"
  assert_success
  assert_output '{"a":"b"}'
}

@test "include" {
  run bash -c "echo '{\"foo\":\"bar\", \"baz\":\"qux\", \"a\":\"b\"}' | $DAP_EXECUTABLE json + include a=c + json"
  assert_success
  assert_output ''
  run bash -c "echo '{\"foo\":\"bar\", \"baz\":\"qux\", \"a\":\"b\"}' | $DAP_EXECUTABLE json + include a=b + json"
  assert_success
  assert_output '{"a":"b","baz":"qux","foo":"bar"}'
}

@test "transform" {
  run bash -c "echo '{\"foo\":\"bar\"}' | $DAP_EXECUTABLE json + transform foo=base64encode + json"
  assert_success
  assert_output '{"foo":"YmFy"}'
}

@test "recog_match" {
  run bash -c "echo '9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2' | $DAP_EXECUTABLE lines + recog line=dns.versionbind + json | jq -Sc ."
  assert_success
  assert_output '{"line":"9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2","line.recog.os.cpe23":"cpe:/o:redhat:enterprise_linux:6","line.recog.os.family":"Linux","line.recog.os.product":"Enterprise Linux","line.recog.os.vendor":"Red Hat","line.recog.os.version":"6","line.recog.os.version.version":"9","line.recog.service.cpe23":"cpe:/a:isc:bind:9.8.2rc1","line.recog.service.family":"BIND","line.recog.service.product":"BIND","line.recog.service.vendor":"ISC","line.recog.service.version":"9.8.2rc1"}'
}

@test "recog_nomatch" {
  run bash -c "echo 'should not match' | $DAP_EXECUTABLE lines + recog line=dns.versionbind + json | jq -Sc ."
  assert_success
  assert_output '{"line":"should not match"}'
}

@test "recog_invalid_arg" {
  run bash -c "echo 'test' | $DAP_EXECUTABLE lines + recog + json"
  assert_failure
}

@test "geo_ip yields valid fields" {
  run bash -c "echo 66.92.181.240 | GEOIP_CITY_DATABASE_PATH=./test/test_data/geoip/GeoIPCity.dat $DAP_EXECUTABLE lines + geo_ip line + json | jq -Sc ."
  assert_success
  assert_output '{"line":"66.92.181.240","line.area_code":"510","line.city":"Fremont","line.country_code":"US","line.country_code3":"USA","line.country_name":"United States","line.dma_code":"807","line.latitude":"37.50790023803711","line.longitude":"-121.95999908447266","line.postal_code":"94538","line.region":"CA","line.region_name":"California"}'
}

@test "geo_ip_org yields valid fields" {
  run bash -c "echo 12.87.118.0 | GEOIP_ORG_DATABASE_PATH=./test/test_data/geoip/GeoIPOrg.dat $DAP_EXECUTABLE lines + geo_ip_org line + json | jq -Sc -r ."
  assert_success
  assert_output '{"line":"12.87.118.0","line.org":"AT&T Worldnet Services"}'
}

@test "geo_ip_asn" {
  run bash -c "echo 1.128.0.0 | GEOIP_ASN_DATABASE_PATH=./test/test_data/geoip/GeoIPASNum.dat $DAP_EXECUTABLE lines + geo_ip_asn line + json | jq -Sc -r ."
  assert_success
  assert_output '{"line":"1.128.0.0","line.asn":"AS1221"}'
}

@test "geo_ip2_city yields valid fields" {
  run bash -c "echo 81.2.69.142 | GEOIP2_LANGUAGE=en GEOIP2_CITY_DATABASE_PATH=./test/test_data/geoip2/GeoIP2-City-Test.mmdb $DAP_EXECUTABLE lines + geo_ip2_city line + json | jq -Sc ."
  assert_success
  assert_output '{"line":"81.2.69.142","line.geoip2.city.geoname_id":"2643743","line.geoip2.city.name":"London","line.geoip2.continent.code":"EU","line.geoip2.continent.geoname_id":"6255148","line.geoip2.continent.name":"Europe","line.geoip2.country.geoname_id":"2635167","line.geoip2.country.is_eu":"true","line.geoip2.country.iso_code":"GB","line.geoip2.country.name":"United Kingdom","line.geoip2.location.accuracy_raidus":"10","line.geoip2.location.latitude":"51.5142","line.geoip2.location.longitude":"-0.0931","line.geoip2.location.metro_code":"0","line.geoip2.location.time_zone":"Europe/London","line.geoip2.postal.code":"","line.geoip2.registered_country.geoname_id":"6252001","line.geoip2.registered_country.is_eu":"false","line.geoip2.registered_country.iso_code":"US","line.geoip2.registered_country.name":"United States","line.geoip2.represented_country.geoname_id":"0","line.geoip2.represented_country.is_eu":"false","line.geoip2.represented_country.iso_code":"","line.geoip2.represented_country.name":"","line.geoip2.represented_country.type":"","line.geoip2.subdivisions.0.geoname_id":"6269131","line.geoip2.subdivisions.0.iso_code":"ENG","line.geoip2.subdivisions.0.name":"England","line.geoip2.subdivisions.length":"1","line.geoip2.traits.is_anon_proxy":"false","line.geoip2.traits.is_satellite":"false"}'
}

@test "geo_ip2_isp yields valid fields" {
  run bash -c "echo 1.128.0.0 | GEOIP2_LANGUAGE=en GEOIP2_ISP_DATABASE_PATH=./test/test_data/geoip2/GeoIP2-ISP-Test.mmdb $DAP_EXECUTABLE lines + geo_ip2_isp line + json | jq -Sc ."
  assert_success
  assert_output '{"line":"1.128.0.0","line.geoip2.isp.asn":"1221","line.geoip2.isp.asn_org":"Telstra Pty Ltd","line.geoip2.isp.isp":"Telstra Internet","line.geoip2.isp.org":"Telstra Internet"}'
}
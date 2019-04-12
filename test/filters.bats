#!/usr/bin/env bats

load ./test_common

@test "rename" {
  run bash -c 'echo world | dap lines + rename line=hello + json'
  assert_success
  assert_output '{"hello":"world"}'
}

@test "not_exists" {
  run bash -c 'echo '{"foo":"bar"}' | dap json + not_exists foo + json'
  assert_success
  assert_output ''
  run bash -c 'echo '{"bar":"bar"}' | dap json + not_exists foo + json'
  assert_success
  assert_output '{"bar":"bar"}'
}

@test "split_comma" {
  run bash -c 'echo '{"foo":"bar,baz"}' | dap json + split_comma foo + json'
  assert_success
  assert_output '{"foo":"bar,baz","foo.word":"bar"}\n{"foo":"bar,baz","foo.word":"baz"}'
}

@test "field_split_line" {
  run bash -c 'echo -e '{"foo":"bar\nbaz"}' | dap json + field_split_line foo + json'
  assert_success
  assert_output '{"foo":"bar\nbaz","foo.f1":"bar","foo.f2":"baz"}'
}

@test "not_empty" {
  run bash -c 'echo '{"foo":"bar,baz"}' | dap json + not_empty foo + json'
  assert_success
  assert_output '{"foo":"bar,baz"}'
}

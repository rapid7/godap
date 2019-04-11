#!/usr/bin/env bats

load ./test_common

@test "reads json" {
  run bash -c 'echo "{\"foo\": 1 }" | godap json + json'
  assert_success
  assert_output '{"foo":1}'
}

@test "reads lines" {
  run bash -c 'echo hello world | godap lines + json'
  assert_success
  assert_output '{"line":"hello world"}'
}

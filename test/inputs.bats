#!/usr/bin/env bats

load ./test_common

@test "reads json" {
  run bash -c 'echo "{\"foo\": 1 }" | $DAP_EXECUTABLE json + json'
  assert_success
  assert_output '{"foo":1}'
}

@test "reads lines" {
  run bash -c 'echo hello world | $DAP_EXECUTABLE lines + json'
  assert_success
  assert_output '{"line":"hello world"}'
}

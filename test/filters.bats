#!/usr/bin/env bats

load ./test_common

@test "rename" {
  run bash -c 'echo world | dap lines + rename line=hello + json'
  assert_success
  assert_output '{"hello":"world"}'
}

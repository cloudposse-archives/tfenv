function setup() {
  export FOOBAR=123
  export Blah=true
  export _something=good
  export TF_ENV=${TF_ENV:-../release/tfenv}
}

function teardown() {
  unset FOOBAR
  unset Blah
  unset _something
  unset TF_ENV
}

@test "TF_VAR_foobar works" {
  which ${TF_ENV}
  ${TF_ENV} printenv TF_VAR_foobar >&2
  [ "$(${TF_ENV} printenv TF_VAR_foobar)" != "" ]
  [ "$(${TF_ENV} printenv TF_VAR_foobar)" == "${FOOBAR}" ]
}

@test "TF_VAR_blah works" {
  which ${TF_ENV}
  ${TF_ENV} printenv TF_VAR_blah >&2
  [ "$(${TF_ENV} printenv TF_VAR_blah)" != "" ]
  [ "$(${TF_ENV} printenv TF_VAR_blah)" == "${Blah}" ]
}

@test "TF_VAR_something works" {
  which ${TF_ENV}
  ${TF_ENV} printenv TF_VAR_something >&2
  [ "$(${TF_ENV} printenv TF_VAR_something)" != "" ]
  [ "$(${TF_ENV} printenv TF_VAR_something)" == "${_something}" ]
}

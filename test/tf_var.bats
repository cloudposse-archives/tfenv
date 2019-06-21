function setup() {
  export FOOBAR=123
  export Blah=true
  export _something=good
  export _PATH=$PATH
  export PATH="../release:${PATH}"
}

function teardown() {
  unset FOOBAR
  unset Blah
  unset _something
  export PATH=$_PATH
}

@test "TF_VAR_foobar works" {
  which tfenv
  [ "$(tfenv printenv TF_VAR_foobar)" != "" ]
  [ "$(tfenv printenv TF_VAR_foobar)" == "${FOOBAR}" ]
}

@test "TF_VAR_blah works" {
  which tfenv
  [ "$(tfenv printenv TF_VAR_blah)" != "" ]
  [ "$(tfenv printenv TF_VAR_blah)" == "${Blah}" ]
}

@test "TF_VAR_something works" {
  which tfenv
  [ "$(tfenv printenv TF_VAR_something)" != "" ]
  [ "$(tfenv printenv TF_VAR_something)" == "${_something}" ]
}

function setup() {
  export TF_CLI_INIT_FROM_MODULE="git::https://https://github.com/cloudposse/terraform-null-label?ref=master"
  export TF_CLI_INIT_BACKEND=false
  export TF_CLI_INIT="module/"
  export TF_ENV=${TF_ENV:-../release/tfenv}
}

function teardown() {
  unset TF_CLI_INIT_FROM_MODULE
  unset TF_CLI_INIT_BACKEND
  unset TF_CLI_INIT
  unset TF_ENV
}

@test "TF_CLI_ARGS_init works" {
  which ${TF_ENV}
  ${TF_ENV} printenv TF_CLI_ARGS_init >&2
  [ "$(${TF_ENV} printenv TF_CLI_ARGS_init)" != "" ]
  [ "$(${TF_ENV} printenv TF_CLI_ARGS_init)" == "-backend=${TF_CLI_INIT_BACKEND} -from-module=${TF_CLI_INIT_FROM_MODULE} ${TF_CLI_INIT}" ]
}

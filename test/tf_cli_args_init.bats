function setup() {
  export TF_CLI_INIT_FROM_MODULE="git::https://https://github.com/cloudposse/terraform-null-label?ref=master"
  export TF_CLI_INIT_BACKEND=false
  export TF_CLI_INIT="module/"
  export _PATH=$PATH
  export PATH="../release:${PATH}"
}

function teardown() {
  unset TF_CLI_INIT_FROM_MODULE
  unset TF_CLI_INIT
  export PATH=$_PATH
}

@test "TF_CLI_ARGS_init works" {
  which tfenv
  [ "$(tfenv printenv TF_CLI_ARGS_init)" != "" ]
  [ "$(tfenv printenv TF_CLI_ARGS_init)" == "-backend=${TF_CLI_INIT_BACKEND} -from-module=${TF_CLI_INIT_FROM_MODULE} ${TF_CLI_INIT}" ]
}

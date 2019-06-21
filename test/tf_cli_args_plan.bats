function setup() {
  export TF_CLI_PLAN_AUTO_APPROVE="true"
  export TF_CLI_PLAN_NO_COLOR="true"
  export TF_CLI_PLAN_OUT="plan.txt"
  export TF_CLI_PLAN="module/"
  export _PATH=$PATH
  export PATH="../release:${PATH}"
}

function teardown() {
  unset TF_CLI_PLAN_AUTO_APPROVE
  unset TF_CLI_PLAN_NO_COLOR
  unset TF_CLI_PLAN_OUT
  unset TF_CLI_PLAN
  export PATH=$_PATH
}

@test "TF_CLI_ARGS_plan works" {
  which tfenv
  tfenv printenv TF_CLI_ARGS_plan >&2
  [ "$(tfenv printenv TF_CLI_ARGS_plan)" != "" ]
  [ "$(tfenv printenv TF_CLI_ARGS_plan)" == "-no-color -out=${TF_CLI_PLAN_OUT} -auto-approve $TF_CLI_PLAN" ]
}

function setup() {
  export TF_CLI_PLAN_AUTO_APPROVE="true"
  export TF_CLI_PLAN_NO_COLOR="true"
  export TF_CLI_PLAN_OUT="plan.txt"
  export TF_CLI_PLAN="module/"
  export TF_ENV=${TF_ENV:-../release/tfenv}
}

function teardown() {
  unset TF_CLI_PLAN_AUTO_APPROVE
  unset TF_CLI_PLAN_NO_COLOR
  unset TF_CLI_PLAN_OUT
  unset TF_CLI_PLAN
  unset TF_ENV
}

@test "TF_CLI_ARGS_plan works" {
  which ${TF_ENV}
  ${TF_ENV} printenv TF_CLI_ARGS_plan >&2
  [ "$(${TF_ENV} printenv TF_CLI_ARGS_plan)" != "" ]
  [ "$(${TF_ENV} printenv TF_CLI_ARGS_plan)" == "-no-color -out=${TF_CLI_PLAN_OUT} -auto-approve $TF_CLI_PLAN" ]
}

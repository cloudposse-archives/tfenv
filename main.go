package main

import (
	"fmt"
	"github.com/taskcluster/shell"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

func init() {
	// make sure we only have one process and that it runs on the main thread
	// (so that ideally, when we Exec, we keep our user switches and stuff)
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	log.SetFlags(0) // no timestamps on our logs

	// Args that we pass to exec
	var args []string

	// List of environment variables that will be passed to the executable.
	var env = []string{}

	// Prefix used for all normalized environment variables
	var tfenvPrefix = getEnv("TFENV_PREFIX", "TF_VAR_")

	// Whitelist of allowed environment variables. Processed *after* blacklist.
	var tfenvWhitelist = getEnv("TFENV_WHITELIST", ".*")

	// Blacklist of excluded environment variables. Processed *before* whitelist.
	var tfenvBlacklist = getEnv("TFENV_BLACKLIST", "^(AWS_ACCESS_KEY_ID|AWS_SECRET_ACCESS_KEY)$")

	// Args that we pass to TF_CLI_ARGS_init
	var tfCliArgsInit []string
	var tfCliArgsPlan []string
	var tfCliArgsApply []string
	var tfCliArgsDestroy []string
	var tfCliArgs []string

	reTfCliInitBackend := regexp.MustCompile("^TF_CLI_INIT_BACKEND_CONFIG_(.*)")
	reTfCliCommand := regexp.MustCompile("^TF_CLI_(INIT|PLAN|APPLY|DESTROY)_(.*)")
	reTfCliDefault := regexp.MustCompile("^TF_CLI_DEFAULT_(.*)")

	reTfVar := regexp.MustCompile("^" + tfenvPrefix)
	reTrim := regexp.MustCompile("(^_+|_+$)")
	reDashes := regexp.MustCompile("-+")
	reUnderscores := regexp.MustCompile("_+")
	reWhitelist := regexp.MustCompile(tfenvWhitelist)
	reBlacklist := regexp.MustCompile(tfenvBlacklist)

	for _, e := range os.Environ() {
		// Preserve the original environment variable
		env = append(env, e)

		// Begin normalization of environment variable
		pair := strings.SplitN(e, "=", 2)

		originalEnvName := pair[0]

		// `TF_CLI_ARGS_init`: Map `TF_CLI_INIT_BACKEND_CONFIG_FOO=value` to `-backend-config=foo=value`
		if reTfCliInitBackend.MatchString(pair[0]) {
			match := reTfCliInitBackend.FindStringSubmatch(pair[0])

			// Replace all underscores with dashes
			arg := reUnderscores.ReplaceAllString(match[1], "-")

			// Lowercase parameters for terraform
			arg = strings.ToLower(arg)

			// Convert things like `role-arn` to `role_arn`
			arg = reDashes.ReplaceAllString(arg, "_")

			// Combine parameters into something like `-backend-config=role_arn=xxx`
			arg = "-backend-config=" + arg + "=" + pair[1]
			tfCliArgsInit = append(tfCliArgsInit, arg)
		} else if reTfCliCommand.MatchString(pair[0]) {
			// `TF_CLI_ARGS_plan`: Map `TF_CLI_PLAN_SOMETHING=value` to `-something=value`
			match := reTfCliCommand.FindStringSubmatch(pair[0])
			cmd := reUnderscores.ReplaceAllString(match[1], "-")
			cmd = strings.ToLower(cmd)

			param := reUnderscores.ReplaceAllString(match[2], "-")
			param = strings.ToLower(param)
			arg := "-" + param + "=" + pair[1]
			switch cmd {
			case "init":
				tfCliArgsInit = append(tfCliArgsInit, arg)
			case "plan":
				tfCliArgsPlan = append(tfCliArgsPlan, arg)
			case "apply":
				tfCliArgsApply = append(tfCliArgsApply, arg)
			case "destroy":
				tfCliArgsDestroy = append(tfCliArgsDestroy, arg)
			}
		} else if reTfCliDefault.MatchString(pair[0]) {
			// `TF_CLI_ARGS`: Map `TF_CLI_DEFAULT_SOMETHING=value` to `-something=value`
			match := reTfCliDefault.FindStringSubmatch(pair[0])
			param := reUnderscores.ReplaceAllString(match[1], "-")
			param = strings.ToLower(param)
			arg := "-" + param + "=" + pair[1]
			tfCliArgs = append(tfCliArgs, arg)
		} else if !reBlacklist.MatchString(pair[0]) && reWhitelist.MatchString(pair[0]) {
			// Process the blacklist for exclusions, then the whitelist for inclusions
			// Strip off TF_VAR_ prefix so we can simplify normalization
			pair[0] = reTfVar.ReplaceAllString(pair[0], "")

			// downcase key
			pair[0] = strings.ToLower(pair[0])

			// trim leading and trailing underscores
			pair[0] = reTrim.ReplaceAllString(pair[0], "")

			// remove consecutive underscores
			pair[0] = reUnderscores.ReplaceAllString(pair[0], "_")

			// prepend TF_VAR_, if not there already
			if len(pair[0]) > 0 {
				pair[0] = tfenvPrefix + pair[0]
				if strings.Compare(pair[0], originalEnvName) != 0 {
					envvar := pair[0] + "=" + pair[1]
					//fmt.Println(envvar)
					env = append(env, envvar)
				}
			}
		}
	}

	if len(tfCliArgsInit) > 0 {
		env = append(env, "TF_CLI_ARGS_init="+strings.Join(tfCliArgsInit, " "))
	}

	if len(tfCliArgsPlan) > 0 {
		env = append(env, "TF_CLI_ARGS_plan="+strings.Join(tfCliArgsPlan, " "))
	}

	if len(tfCliArgsApply) > 0 {
		env = append(env, "TF_CLI_ARGS_apply="+strings.Join(tfCliArgsApply, " "))
	}

	if len(tfCliArgsDestroy) > 0 {
		env = append(env, "TF_CLI_ARGS_destroy="+strings.Join(tfCliArgsDestroy, " "))
	}

	if len(tfCliArgs) > 0 {
		env = append(env, "TF_CLI_ARGS="+strings.Join(tfCliArgs, " "))
	}

	sort.Strings(env)

	// The command that was executed
	cmd := os.Args[0]

	if len(os.Args) < 2 {
		for _, envvar := range env {
			// Begin normalization of environment variable
			pair := strings.SplitN(envvar, "=", 2)
			fmt.Printf("export %v=%v\n", pair[0], shell.Escape(pair[1]))
		}
	} else {
		// The command that will be executed
		exe := os.Args[1]

		// The command + any arguments
		args = append(args, os.Args[1:]...)

		// Lookup path for executable
		binary, binaryPathErr := exec.LookPath(exe)
		if binaryPathErr != nil {
			log.Fatalf("error: %v failed to find executable `%v`: %v", cmd, exe, binaryPathErr)
		}

		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			log.Fatalf("error: %v exec failed: %v", cmd, execErr)
		}
	}
}

package main

import (
	//	"fmt"
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
	var tfenv_prefix = getEnv("TFENV_PREFIX", "TF_VAR_")

	// Whitelist of allowed environment variables. Processed *after* blacklist.
	var tfenv_whitelist = getEnv("TFENV_WHITELIST", ".*")

	// Blacklist of excluded environment variables. Processed *before* whitelist.
	var tfenv_blacklist = getEnv("TFENV_BLACKLIST", "^(AWS_ACCESS_KEY_ID|AWS_SECRET_ACCESS_KEY)$")

	// The command that was executed
	cmd := os.Args[0]

	re_tf_var := regexp.MustCompile("^" + tfenv_prefix)
	re_trim := regexp.MustCompile("(^_+|_+$)")
	re_dedupe := regexp.MustCompile("_+")
	re_whitelist := regexp.MustCompile(tfenv_whitelist)
	re_blacklist := regexp.MustCompile(tfenv_blacklist)

	for _, e := range os.Environ() {
		// Preserve the original environment variable
		env = append(env, e)

		// Begin normalization of environment variable
		pair := strings.Split(e, "=")

		// Process the blacklist for exclusions, then the whitelist for inclusions
		if !re_blacklist.MatchString(pair[0]) && re_whitelist.MatchString(pair[0]) {
			// Strip off TF_VAR_ prefix so we can simplify normalization
			pair[0] = re_tf_var.ReplaceAllString(pair[0], "")

			// downcase key
			pair[0] = strings.ToLower(pair[0])

			// trim leading and trailing underscopres
			pair[0] = re_trim.ReplaceAllString(pair[0], "")

			// remove consequtive underscopres
			pair[0] = re_dedupe.ReplaceAllString(pair[0], "_")

			// prepend TF_VAR_, if not there already
			if len(pair[0]) != 0 {
				pair[0] = tfenv_prefix + pair[0]
				envvar := pair[0] + "=" + pair[1]
				//fmt.Println(envvar)
				env = append(env, envvar)
			}
		}
	}
	sort.Strings(env)

	if len(os.Args) < 2 {
		log.Fatalf("error: %v command args...", cmd)
	}

	// The command that will be executed
	exe := os.Args[1]

	// The command + any arguments
	args = append(args, os.Args[1:]...)

	// Lookup path for executable
	binary, binaryPathErr := exec.LookPath(exe)
	if binaryPathErr != nil {
		log.Fatalf("error: find to find executable `%v`: %v", exe, binaryPathErr)
	}

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		log.Fatalf("error: exec failed: %v", execErr)
	}
}

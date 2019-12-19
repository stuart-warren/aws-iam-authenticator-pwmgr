package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Entry struct {
	ID       string
	Name     string
	Username string
	Password string
}

func mergeEnv(original, extra []string) []string {
	m := map[string]string{}
	if len(extra) == 0 {
		return original
	}
	for _, e := range append(original, extra...) {
		p := strings.SplitN(e, "=", 2)
		m[p[0]] = p[1]
	}
	var out []string
	for k, v := range m {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}

func run(command string, args ...string) ([]byte, error) {
	return runWithEnv([]string{}, command, args...)
}

func runWithEnv(env []string, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	cmd.Env = mergeEnv(cmd.Env, env)
	err := cmd.Run()
	if errout.String() != "" {
		err = fmt.Errorf("%s", errout.String())
	}
	return out.Bytes(), err
}

func log(msg ...interface{}) {
	fmt.Println(msg...)
}

func fatal(msg ...interface{}) {
	log(msg...)
	os.Exit(1)
}

func main() {
	name := os.Getenv("AWS_PROFILE")
	if name == "" {
		fatal("Required env var not set: AWS_PROFILE")
	}
	out, err := run("lpass", "show", "--sync=auto", "--json", name)
	if err != nil {
		fatal(err)
	}
	dec := json.NewDecoder(bytes.NewBuffer(out))
	var entries []Entry
	err = dec.Decode(&entries)
	if err != nil {
		log("Unexpected output from lpass, perhaps multiple entries with same name? Please be more specific")
		fatal(string(out))
	}
	accessKey := entries[0].Username
	secretKey := entries[0].Password
	env := []string{
		"AWS_ACCESS_KEY_ID=" + accessKey,
		"AWS_SECRET_ACCESS_KEY=" + secretKey,
		"AWS_PROFILE=",
	}
	out, err = runWithEnv(env, "aws-iam-authenticator", os.Args[1:]...)
	if err != nil {
		fatal(err)
	}
	fmt.Print(string(out))
}

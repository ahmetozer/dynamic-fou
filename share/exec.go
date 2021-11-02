package share

import (
	"bytes"
	"os/exec"
)

func (e Env) Exec(Command string, arg ...string) (string, string) {
	cmd := exec.Command(Command, arg...)

	var Stdout bytes.Buffer
	var Stderr bytes.Buffer
	cmd.Stdout = &Stdout
	cmd.Stderr = &Stderr

	cmd.Env = e
	err := cmd.Run()

	if err != nil {
		return Stdout.String(), err.Error()
	}
	return Stdout.String(), Stderr.String()
}

type Env []string

package main

import (
	"bytes"
	"os/exec"
)

type JadeConverter struct {
	HtmlConverter
}

func (c JadeConverter) Convert(src []byte) ([]byte, error) {
	echo_src := exec.Command("echo", string(src))
	stdout, err := echo_src.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	jade := exec.Command("jade")
	jade.Stdin = stdout
	jade.Stdout = &out

	cmds := []*exec.Cmd{echo_src, jade}
	for _, c := range cmds {
		if err := c.Start(); err != nil {
			return nil, err
		}
	}
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			return nil, err
		}
	}

	return out.Bytes(), nil
}

func (c JadeConverter) IsAvailable() bool {
	err := exec.Command("jade", "--version").Run()
	return err == nil
}

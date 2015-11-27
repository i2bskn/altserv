package converter

import (
	"bytes"
	"os/exec"
)

type StylusConverter struct {
	CSSConverter
}

func (c StylusConverter) Convert(src []byte) ([]byte, error) {
	echoSrc := exec.Command("echo", string(src))
	stdout, err := echoSrc.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	stylus := exec.Command("stylus", "--print")
	stylus.Stdin = stdout
	stylus.Stdout = &out

	cmds := []*exec.Cmd{echoSrc, stylus}
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

func (c StylusConverter) IsAvailable() bool {
	err := exec.Command("stylus", "--version").Run()
	return err == nil
}

package main

import (
	"bytes"
	"os/exec"
)

type Converter interface {
	Convert(src []byte) ([]byte, error)
	IsAvailable() bool
	ConvertedExt() string
}

type AvailableConverters struct {
	Converters map[string]Converter
	ConvertMap map[string][]string
}

func newAvailableConverters() *AvailableConverters {
	converters := allConverters()
	for from_ext, c := range converters {
		if !c.IsAvailable() {
			delete(converters, from_ext)
		}
	}

	convert_map := make(map[string][]string)
	for from_ext, c := range converters {
		to_ext := c.ConvertedExt()
		convert_map[to_ext] = append(convert_map[to_ext], from_ext)
	}
	return &AvailableConverters{
		Converters: converters,
		ConvertMap: convert_map,
	}
}

func allConverters() map[string]Converter {
	jade := new(JadeConverter)
	return map[string]Converter{
		".jade": jade,
	}
}

func (c *AvailableConverters) Convert(src []byte, t string) []byte {
	if len(t) == 0 {
		return src
	}

	if converter, exist := c.Converters[t]; exist {
		if converted, err := converter.Convert(src); err == nil {
			return converted
		}
	}

	return src
}

type HtmlConverter struct{}

func (c HtmlConverter) ConvertedExt() string {
	return ".html"
}

type CssConverter struct{}

func (c CssConverter) ConvertedExt() string {
	return ".css"
}

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
	_, err := exec.Command("jade", "--version").Output()
	return err == nil
}

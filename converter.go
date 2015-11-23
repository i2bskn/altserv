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
	for ext, c := range converters {
		if !c.IsAvailable() {
			delete(converters, ext)
		}
	}

	convert_map := make(map[string][]string)
	for ext, c := range converters {
		conv_ext := c.ConvertedExt()
		_, exist := convert_map[conv_ext]
		if !exist {
			convert_map[conv_ext] = make([]string, 0, 0)
		}
		convert_map[conv_ext] = append(convert_map[conv_ext], ext)
	}
	return &AvailableConverters{
		Converters: converters,
		ConvertMap: convert_map,
	}
}

func allConverters() map[string]Converter {
	var jade JadeConverter
	return map[string]Converter{
		".jade": jade,
	}
}

func (c *AvailableConverters) Convert(src []byte, t string) []byte {
	if len(t) == 0 {
		return src
	}

	converter, exist := c.Converters[t]
	if exist {
		converted, err := converter.Convert(src)
		if err == nil {
			return converted
		}
	}

	return src
}

type HtmlConverter struct{}

func (c HtmlConverter) ConvertedExt() string {
	return ".html"
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
		err := c.Start()
		if err != nil {
			return nil, err
		}
	}
	for _, c := range cmds {
		err := c.Wait()
		if err != nil {
			return nil, err
		}
	}

	return out.Bytes(), nil
}

func (c JadeConverter) IsAvailable() bool {
	_, err := exec.Command("jade", "--version").Output()
	return err == nil
}

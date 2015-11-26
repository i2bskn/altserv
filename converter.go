package main

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
	stylus := new(StylusConverter)

	return map[string]Converter{
		".jade": jade,
		".styl": stylus,
	}
}

func (c *AvailableConverters) Convert(src []byte, t string) ([]byte, string) {
	if len(t) == 0 {
		return src, ""
	}

	if converter, exist := c.Converters[t]; exist {
		if converted, err := converter.Convert(src); err == nil {
			return converted, converter.ConvertedExt()
		}
	}

	return src, ""
}

type HtmlConverter struct{}

func (c HtmlConverter) ConvertedExt() string {
	return ".html"
}

type CssConverter struct{}

func (c CssConverter) ConvertedExt() string {
	return ".css"
}

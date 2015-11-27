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
	for fromExt, c := range converters {
		if !c.IsAvailable() {
			delete(converters, fromExt)
		}
	}

	convertMap := make(map[string][]string)
	for fromExt, c := range converters {
		toExt := c.ConvertedExt()
		convertMap[toExt] = append(convertMap[toExt], fromExt)
	}
	return &AvailableConverters{
		Converters: converters,
		ConvertMap: convertMap,
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

type HTMLConverter struct{}

func (c HTMLConverter) ConvertedExt() string {
	return ".html"
}

type CSSConverter struct{}

func (c CSSConverter) ConvertedExt() string {
	return ".css"
}

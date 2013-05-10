package GoVideoParser

type DefinitionType uint8

const (
	DT_NORMAL = (DefinitionType)(0)
	DT_HIGH   = (DefinitionType)(1)
	DT_SUPER  = (DefinitionType)(2)
)

type ParserType uint8

const (
	PT_YOUKU = (ParserType)(0)
)

type IParserError interface {
	Error() string
}

type VideoParseResult interface {
	GetURLS() []string
	GetTitle() string
	GetFileType() string
}

type IParser interface {
	Parse(url string, defi DefinitionType) (VideoParseResult, error)
	GetType() ParserType
}

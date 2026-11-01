package parser

type ParsedFiled struct {
	Original string
	Context  string
	Part     int
	Index    int
	Version  int
	Slug     string
	IsPart   bool
}

package tools

type Position struct {
	Idx  int
	Ln   int
	Col  int
	Fn   string
	Ftxt string
}

func (position *Position) AdvanceN(berapa int) {
	position.Idx += berapa
	position.Col += berapa
}

func (position *Position) Advance(CurrentChar string) *Position {
	position.Idx++
	position.Col++

	if CurrentChar == "\n" {
		position.Ln++
		position.Col = 0
	}

	return position
}

func (position *Position) Copy() *Position {
	return &Position{
		Idx:  position.Idx,
		Ln:   position.Ln,
		Col:  position.Col,
		Fn:   position.Fn,
		Ftxt: position.Ftxt,
	}
}

func GetComparison(hasil bool) int {
	if hasil {
		return 1
	}

	return 0
}

var SemuaBuiltInFunction map[string]string = map[string]string{
	"print": "Print",
	"PRINT": "Print",
	"write": "Print",
	"WRITE": "Print",
	"read":  "Input",
	"READ":  "Input",
	"input": "Input",
	"INPUT": "Input",
}

func ApakahBuiltinFunction(s string) bool {
	return SemuaBuiltInFunction[s] != ""
}

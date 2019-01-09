package scanner

type Mode int

const (
	ScanComments Mode = 1 << iota
	DontInsertSemis
)

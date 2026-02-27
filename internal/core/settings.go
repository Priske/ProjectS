package core

type Settings struct {
	BoardW, BoardH int
	CellSize       int
}

const (
	DefaultBoardW   = 10
	DefaultBoardH   = 10
	DefaultCellSize = 64
	VirtualW        = 1280
	VirtualH        = 720
	MenuButtonH     = 50
	MenuButtonW     = 200
)

func DefaultSettings() Settings {
	return Settings{
		BoardW:   DefaultBoardW,
		BoardH:   DefaultBoardH,
		CellSize: DefaultCellSize,
	}
}

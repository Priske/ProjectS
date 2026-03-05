package core

// MousePositionX,MousePositionY
type Input struct {
	MX, MY int

	LeftPressed  bool
	LeftClicked  bool
	RightClicked bool
	RightPressed bool

	RuneBuffer []rune
	Backspace  bool
	Escape     bool
}

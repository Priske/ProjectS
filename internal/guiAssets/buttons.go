package guiassets

type Button struct {
	X, Y, W, H int
	Text       string
	OnClick    func()
}

func (b Button) Contains(mx, my int) bool {
	return mx >= b.X && mx < b.X+b.W && my >= b.Y && my < b.Y+b.H
}

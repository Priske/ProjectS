package interaction

type DragState struct {
	Active bool

	// Where the drag started (optional but useful for boards)
	FromX, FromY int

	// Optional payload (unit pointer, unit ID, widget ID, etc.)
	Payload any

	// Mouse offset inside the grabbed object at drag-start
	GrabOffX, GrabOffY int

	// Live cursor for drawing
	MX, MY int
}

func (d *DragState) Begin(fromX, fromY int, payload any, mx, my, grabOffX, grabOffY int) {
	d.Active = true
	d.FromX, d.FromY = fromX, fromY
	d.Payload = payload
	d.MX, d.MY = mx, my
	d.GrabOffX, d.GrabOffY = grabOffX, grabOffY
}

func (d *DragState) UpdateMouse(mx, my int) {
	d.MX, d.MY = mx, my
}

func (d *DragState) End() {
	*d = DragState{}
}

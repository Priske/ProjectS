package core

type Formation struct {
	Name  string
	W, H  int              // e.g. 5x5 editor (optional)
	Wants map[Pos]UnitType // only cells the player reserved
}
type Deployment struct {
	At map[int]Pos // UnitId -> absolute board pos
}

// Missing entry: which cell wanted what type but couldn't be filled
type MissingSlot struct {
	Rel  Pos
	Want UnitType
}

func (f Formation) BuildDeployment(units []*Unit, anchor Pos) (Deployment, []MissingSlot) {
	used := map[int]bool{}
	deployment := Deployment{At: map[int]Pos{}}
	missing := make([]MissingSlot, 0)

	for rel, wantType := range f.Wants {
		u := findAndMark(units, wantType, used)
		if u == nil {
			missing = append(missing, MissingSlot{Rel: rel, Want: wantType})
			continue
		}

		abs := Pos{X: anchor.X + rel.X, Y: anchor.Y + rel.Y}
		deployment.At[u.UnitId] = abs
	}

	return deployment, missing
}

func findAndMark(units []*Unit, want UnitType, used map[int]bool) *Unit {
	for _, u := range units {
		if u == nil {
			continue
		}
		if u.Type != want {
			continue
		}
		if used[u.UnitId] {
			continue
		}
		used[u.UnitId] = true
		return u
	}
	return nil
}

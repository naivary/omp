//go:generate go tool go-enum --marshal --sql --nocomments
package v1

// ENUM(Right, Left)
type StrongFoot int

type Player struct {
	ID           int64
	TeamID       int64
	Name         string
	JerseyNumber uint
	StrongFoot   StrongFoot
}

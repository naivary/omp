//go:generate go tool go-enum --marshal --sql --nocomments
package v1

// ENUM(Right, Left)
type StrongFoot int

type PlayerProfile struct {
	ID           int64      `json:"id"`
	TeamID       int64      `json:"teamID"`
	Name         string     `json:"name"`
	JerseyNumber uint       `json:"jerseyNumber"`
	StrongFoot   StrongFoot `json:"strongFoot"`
	Position     string     `json:"position"`
}

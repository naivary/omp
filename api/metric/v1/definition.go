//go:generate go tool go-enum --marshal --sql --nocomments
package v1

// ENUM(Counter, Gauge)
type Type int

// ENUM(Club, Team, Coach, Player)
type Scope int

type Definition struct {
	// ID unique to this definition
	ID int64
	// Name is a user-friendly name which is used on the GUI. It is unique
	// within the choosen scope.
	Name string
	// Type of metric (counter, gauge etc.)
	Type Type
	// Scope defines who has access to this definition.
	Scope Scope
	// ClubID is the owner of the club
	ClubID int64
	// Description of the metric, it's usage and purpose.
	Description string
}

//go:generate go tool go-enum --marshal --sql --nocomments
package v1

// ENUM(Counter, Gauge)
type Type int

// ENUM(Club, Team, TeamPrivate)
type Scope int

type Definition struct {
	// ID unique to this definition
	ID int64
	// Name is a user-friendly name which is used on the GUI. It is unique
	// within the choosen scope.
	Name string
	// Type of metric (counter, gauge etc.)
	Type Type
	// Scope defines who has access to this definition. The following options
	// are defined:
	//
	// 	Club: The metric definition is used to track something on club level
	// 	like amount of teams
	//
	// 	Team: The metric definition is accessible by every team in the club.
	//
	// 	TeamPrivate: The metric definition is only useable by the owner of the
	// 	metric definition.
	//
	// The Club and Team scope can only be created by a Club administration
	// where TeamPrivate can only be created by a Team administrator.
	Scope Scope
	// Owner is the creator of the metric definition. It's either a team ID or
	// club ID. Depending on the scope the owner ID can be inferred. If the
	// Scope is TeamPrivate then the owner has to be a team ID. Anything other
	// it's a club ID. The Owner is still needed for better queries and
	// management.
	Owner int64
	// Description of the metric, it's usage and purpose.
	Description string
}

type CreateMetricRequest struct {
	Name        string
	Type        Type
	Scope       Scope
	Description string
}

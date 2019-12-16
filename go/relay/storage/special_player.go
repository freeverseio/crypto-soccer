package storage

// SpecialPlayerAction is the 'special_player_action' enum type from schema 'public'.
type SpecialPlayerAction uint16

const (
	// SpecialPlayerActionCreate is the 'create' SpecialPlayerAction.
	SpecialPlayerActionCreate = SpecialPlayerAction(1)

	// SpecialPlayerActionDestroy is the 'destroy' SpecialPlayerAction.
	SpecialPlayerActionDestroy = SpecialPlayerAction(2)
)

// String returns the string value of the SpecialPlayerAction.
func (spa SpecialPlayerAction) String() string {
	var enumVal string

	switch spa {
	case SpecialPlayerActionCreate:
		enumVal = "create"

	case SpecialPlayerActionDestroy:
		enumVal = "destroy"
	}

	return enumVal
}

// SpecialPlayer represents a row from 'public.special_players'.
type SpecialPlayer struct {
	SpecialPlayerID string              `json:"special_player_id"` // special_player_id
	Action          SpecialPlayerAction `json:"action"`            // action
}

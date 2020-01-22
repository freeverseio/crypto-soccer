package v1

type OrderStatus int

const (
	DRAFT OrderStatus = iota
	PENDING_PUBLISH
	FAILURE
	PUBLISHED
	CANCELLED
	ACCEPTED
	REJECTED
	PENDING_VALIDATE
	PENDING_RELEASE
	BLOCKED_RELEASE
	RELEASED
)

func (s OrderStatus) String() string {
	names := [...]string{
		"DRAFT",
		"PENDING_PUBLISH",
		"FAILURE",
		"PUBLISHED",
		"CANCELLED",
		"ACCEPTED",
		"REJECTED",
		"PENDING_VALIDATE",
		"PENDING_RELEASE",
		"BLOCKED_RELEASE",
		"RELEASED"}

	if s < DRAFT || s > RELEASED {
		return "UNKNOWN"
	}
	return names[s]
}

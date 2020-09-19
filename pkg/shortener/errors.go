package shortener

// Error represents a generic error
type Error string

// Application domain errors
var (
	ErrLinkNotFound Error = Error("Link not found")
	ErrLinkExists   Error = Error("Link's slug already exists")
	ErrInvalidLink  Error = Error("Link is not valid")
)

func (e Error) Error() string {
	return string(e)
}

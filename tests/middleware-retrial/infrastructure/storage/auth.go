package storage

// Auth defines the interface for authentication storage operations.
type Auth interface {
	Login(userToken string) (bool, error)
}

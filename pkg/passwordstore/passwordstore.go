package passwordstore

//go:generate mockgen -source=passwordstore.go -destination=mocks/mock.go
type PasswordStore interface {
	Add(string)
	Exists(string) bool
	Get() []string
}

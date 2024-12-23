package resourceloader

type ResourceLoader interface {
	Load() ([]string, error)
	ValidateSource() bool
}

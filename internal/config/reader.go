package config

type validatedReader interface {
	ReadValidated() (map[string]string, error)
}

type unvalidatedReader interface {
	ReadUnvalidated() (map[string]string, error)
}

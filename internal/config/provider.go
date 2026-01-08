package config

type provider interface {
	UnmarshalInto(target any) error
}

package resolve

type Config interface {
	Validate() map[string]string
	RequiredKeys() []string
	FillInKeys() error
}

package handlers

type Handler interface {
	TargetToStruct(target string) (any, error)
}

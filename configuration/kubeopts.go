package configuration

type CommandFlag int

const (
	Create CommandFlag = iota
	Apply
	Replace
	Delete
)

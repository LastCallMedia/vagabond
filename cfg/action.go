package cfg

import(
)

type ConfigAction interface {
	GetName() string
	NeedsRun() (bool, error)
	Run() error
}
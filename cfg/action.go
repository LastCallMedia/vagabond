package cfg

import(
)

type ConfigAction interface {
	NeedsRun() (bool, error)
	Run(interface{}) error
}


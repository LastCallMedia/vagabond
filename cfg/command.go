package cfg

type CommandConfigAction struct {
	Condition string
	Command string
	Target  string
	IgnoreReturn bool
}

func (a CommandConfigAction) GetName() string {
	return a.Command
}

func (c *CommandConfigAction)GetIo() VagabondIo {
	if c.Target == "" {
		return VagabondIoLocal{}
	}
	return VagabondIoMachine{c.Target}
}

func (c CommandConfigAction)NeedsRun() (bool, error) {
	if c.Condition == "" {
		return true, nil
	}
	err := c.GetIo().Exec(c.Condition)
	return err != nil, nil
}

func (c CommandConfigAction)Run() error {
	err := c.GetIo().Exec(c.Command)
	if c.IgnoreReturn {
		return nil
	}
	return err
}




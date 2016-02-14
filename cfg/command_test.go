package cfg

import(
	"testing"
)

func TestRun(t *testing.T) {
	c := CommandConfigAction{
		Command: "exit 0",
	}
	err := c.Run()
	if err != nil {
		t.Error(
			"expected", "0",
			"got", err,
		)
	}

	c = CommandConfigAction{
		Command: "exit 1",
	}
	err = c.Run()
	if err == nil {
		t.Error(
			"expected", 1,
			"got", 0,
		)
	}
}

func TestNeedsRun(t *testing.T) {
	c := CommandConfigAction{
		Condition: "exit 0",
	}
	nr, err := c.NeedsRun()
	if nr {
		t.Error("Needs run test 1 failed")
	}
	if err != nil {
		t.Error(
			"expected", nil,
			"got", err,
		)
	}

	c = CommandConfigAction{
		Condition: "exit 1",
	}
	nr, err = c.NeedsRun()
	if !nr {
		t.Error(
			"expected", true,
			"got", nr,
		)
	}
	if err != nil {
		t.Error(
			"expected", nil,
			"got", err,
		)
	}

}
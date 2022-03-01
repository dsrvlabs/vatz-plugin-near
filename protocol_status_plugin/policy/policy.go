package policy

import (
	"log"

	"github.com/shirou/gopsutil/v3/process"
)

type arguments []string

func (args arguments) Contain(values []string) bool {

	findRet := map[string]bool{}

	for _, value := range values {
		for _, a := range args {
			if a == value {
				findRet[value] = true
				break
			}
		}
	}

	if len(findRet) == len(values) {
		return true
	}

	return false
}

// Executor provides interfaces for testing policy.
type Executor interface {
	IsNearUp() (bool, error)
}

type nearExecutor struct {
}

func (e *nearExecutor) IsNearUp() (bool, error) {
	log.Println("IsNearUp")
	return isProcessRunning("IsNearUp", []string{"start"})
}

func isProcessRunning(pName string, mustHaveArgs []string) (bool, error) {
	ps, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range ps {
		var (
			pName     string
			isRunning bool
			err       error
		)

		if pName, err = p.Name(); err != nil {
			continue
		}

		if pName != pName {
			continue
		}

		args, err := p.CmdlineSlice()
		if err != nil {
			continue
		}

		if !arguments(args).Contain(mustHaveArgs) {
			continue
		}

		if isRunning, err = p.IsRunning(); err != nil {
			return false, err
		}

		return isRunning, nil
	}

	return false, nil
}

// NewExecutor returns new executor.
func NewExecutor() Executor {
	return &nearExecutor{}
}

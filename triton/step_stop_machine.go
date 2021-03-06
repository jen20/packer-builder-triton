package triton

import (
	"fmt"
	"time"

	"github.com/joyent/gosdc/cloudapi"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// StepStopMachine stops the machine with the given Machine ID, and waits
// for it to reach the stopped state.
type StepStopMachine struct{}

func (s *StepStopMachine) Run(state multistep.StateBag) multistep.StepAction {
	sdcClient := state.Get("client").(*cloudapi.Client)
	ui := state.Get("ui").(packer.Ui)

	machineID := state.Get("machine").(string)

	ui.Say(fmt.Sprintf("Stopping source machine (%s)...", machineID))
	err := sdcClient.StopMachine(machineID)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem stopping source machine: %s", err))
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Waiting for source machine to stop (%s)...", machineID))
	err = waitForMachineState(sdcClient, machineID, "stopped", 10*time.Minute)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem waiting for source machine to stop: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepStopMachine) Cleanup(state multistep.StateBag) {
	// Explicitly don't clean up here as StepCreateSourceMachine will do it if necessary
	// and there is no real meaning to cleaning this up.
}

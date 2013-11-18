package cloudstack

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"
)

type stepDeployVirtualMachine struct {
	id string
}

func (s *stepDeployVirtualMachine) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*CloudStackClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)
	sshKeyName := state.Get("ssh_key_name")

	ui.Say("Creating virtual machine...")

	// Some random virtual machine name as it's temporary
	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())

	// Create the droplet based on configuration
	id, err := client.DeployVirtualMachine(name, c.ServiceOfferingId, c.TemplateId, c.ZoneId, sshKeyName)
	if err != nil {
		err := fmt.Errorf("Error deploying Virtual Machine: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// We use this in cleanup
	s.id = id

	// Store the virtual machine id for later use
	state.Put("virtual_machine_id", id)

	return multistep.ActionContinue
}

func (s *stepCreateDroplet) Cleanup(state multistep.StateBag) {
	// If the virtual machine id isn't there, we probably never created it
	if s.id == 0 {
		return
	}

	client := state.Get("client").(*CloudStackClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)

	// Destroy the droplet we just created
	ui.Say("Destroying virtual machine...")

	err := client.DestroyVirtualMachine(s.id)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying droplet. Please destroy it manually."))
	}
}

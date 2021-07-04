package kamatera

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"packer-plugin-kamatera/httpclient"
)

type stepPowerOffServer struct{}

func (s *stepPowerOffServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	kamateraClient := state.Get("kamateraClient").(*httpclient.Kamatera)

	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Shutting down server...")

	serverName := state.Get("server_name").(string)

	result, err := kamateraClient.Request("POST", "service/server/terminate", struct {
		Name string `json:"name"`
	}{
		serverName,
	})
	if err != nil {
		err := fmt.Errorf("error power off server: %v", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	commandIds := result.([]interface{})
	_, err = kamateraClient.WaitCommand(commandIds[0].(string))
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepPowerOffServer) Cleanup(state multistep.StateBag) {

}

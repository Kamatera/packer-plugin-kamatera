package kamatera

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/Kamatera/packer-plugin-kamatera/httpclient"
)

type stepGetHduuid struct{}

func (s *stepGetHduuid) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	kamateraClient := state.Get("kamateraClient").(*httpclient.Kamatera)

	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Getting hduuid ...")

	response, err := kamateraClient.Request("POST", "server/hdlib", struct {
		Name string `json:"name"`
	}{
		state.Get("server_name").(string),
	}, false)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	servers := response.([]interface{})
	if len(servers) != 1 {
		err := fmt.Errorf("wrong number of server, got %v", len(servers))
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	server := servers[0].(map[string]interface{})
	hduuid := server["uuid"].(string)

	state.Put("hduuid", hduuid)

	return multistep.ActionContinue
}

func (s *stepGetHduuid) Cleanup(state multistep.StateBag) {

}

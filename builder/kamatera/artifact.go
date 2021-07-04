package kamatera

import (
	"fmt"
	"log"
)

// packersdk.Artifact implementation
type Artifact struct {
	// The name of the snapshot
	snapshotName string

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (*Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.snapshotName
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A snapshot was created: '%v'", a.snapshotName)
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %s", a.snapshotName)
	// TODO: implement
	//_, err := a.hcloudClient.Image.Delete(context.TODO(), &hcloud.Image{ID: a.snapshotId})
	return nil
}

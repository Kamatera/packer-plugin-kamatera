package kamatera

import (
	"fmt"
	"log"
)

type Artifact struct {
	imageName string

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (*Artifact) Files() []string {
	log.Fatal("Downloading image is not supported")
	return nil
}

func (a *Artifact) Id() string {
	return a.imageName
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A image was created: '%v'", a.imageName)
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	//log.Printf("Destroying image: %s", a.imageName)
	// TODO: implement
	log.Fatal("Destroying image is not supported")
	return nil
}

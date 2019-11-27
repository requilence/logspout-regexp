package stderrtransport

import (
	"fmt"
	"os"
)

type StdErrTransport struct{}

func New(_ map[string]string) (*StdErrTransport, error){
	return &StdErrTransport{}, nil
}

func (tgn *StdErrTransport) Name() string {
	return "stderr"
}

func (tgn *StdErrTransport) Write(containerId, containerName, matchedString, re string) error {
	_, err := fmt.Fprintf(os.Stderr, "container: %s(%s) match '%s': %s", containerId, containerName, re, matchedString)

	return err
}


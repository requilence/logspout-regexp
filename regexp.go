package regexp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	re "regexp"
	"strings"

	"github.com/gliderlabs/logspout/router"

	"github.com/requilence/logspout-regexp/transports/stderrtransport"
	"github.com/requilence/logspout-regexp/transports/tgtransport"
)

const defaultRegexpsFile = "logspout_regexps.txt"

func init() {
	router.AdapterFactories.Register(New, "regexp")
}

type Transport interface {
	Name() string
	Write(containerId, containerName, matchedString, re string) error
}

// RegexpAdapter is an adapter that match logs with regexp in file
type RegexpAdapter struct {
	route     *router.Route
	regexps   []*re.Regexp
	transport Transport

	hideMatchedString bool
}

// New creates a RegexpAdapter that looks for regexps in logs
func New(route *router.Route) (router.LogAdapter, error) {
	regexpFile := route.Options["file"]
	if regexpFile == "" {
		regexpFile = defaultRegexpsFile
	}

	file, err := os.Open(regexpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rna := &RegexpAdapter{
		route: route,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		re, err := re.Compile(scanner.Text())
		if err != nil {
			return nil, err
		}

		rna.regexps = append(rna.regexps, re)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var transport string
	adapterParts := strings.Split(route.Adapter, "+")

	if len(adapterParts) < 2 {
		transport = "stderr"
	} else {
		transport = adapterParts[1]
	}

	switch transport {
	case "tg":
		rna.transport, err = tgtransport.New(route.Options)
		if err != nil {
			return nil, err
		}
	case "stderr":
		rna.transport, err = stderrtransport.New(route.Options)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown transport: %s", transport)
	}

	hideMatch := route.Options["hide_matched_string"]
	if hideMatch == "1" {
		rna.hideMatchedString = true
	}

	return rna, nil
}

func (a *RegexpAdapter) Match(s string) (matched bool, re string) {
	for _, singlere := range a.regexps {
		if singlere.MatchString(s) {
			return true, singlere.String()
		}
	}
	return false, ""
}

// Stream implements the router.LogAdapter interface.
func (a *RegexpAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		if matched, re := a.Match(m.Data); matched {
			if a.hideMatchedString {
				m.Data = "***"
			}

			err := a.transport.Write(m.Container.ID, m.Container.Name, m.Data, re)
			if err != nil {
				log.Printf("failed to write %s: %s", a.transport.Name(), err.Error())
			}
		}
	}
}

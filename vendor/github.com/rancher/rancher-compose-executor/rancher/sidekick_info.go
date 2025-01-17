package rancher

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/rancher-compose-executor/project"
)

type SidekickInfo struct {
	primariesToSidekicks map[string][]string
	primaries            map[string]bool
	sidekickToPrimaries  map[string][]string
}

func NewSidekickInfo(project *project.Project) (*SidekickInfo, error) {
	result := &SidekickInfo{
		primariesToSidekicks: map[string][]string{},
		primaries:            map[string]bool{},
		sidekickToPrimaries:  map[string][]string{},
	}

	for _, name := range project.ServiceConfigs.Keys() {
		config, _ := project.ServiceConfigs.Get(name)

		sidekicks := []string{}

		for key, value := range config.Labels {
			if key != "io.rancher.sidekicks" {
				continue
			}

			for _, part := range strings.Split(strings.TrimSpace(value), ",") {
				part = strings.TrimSpace(part)
				result.primaries[name] = true

				sidekicks = append(sidekicks, part)

				list, ok := result.sidekickToPrimaries[part]
				if !ok {
					list = []string{}
				}
				result.sidekickToPrimaries[part] = append(list, name)
			}
		}
		for sidekick, primaries := range result.sidekickToPrimaries {
			if len(primaries) > 1 {
				return nil, errors.Errorf("can't have more than one primary service %v referencing the same sidekick [%v]", primaries, sidekick)
			}
		}

		result.primariesToSidekicks[name] = sidekicks
	}

	return result, nil
}

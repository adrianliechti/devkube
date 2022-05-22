package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

type ListOptions struct {
	All bool

	Filter []string
}

func List(ctx context.Context, options ListOptions) ([]Container, error) {
	var result []Container

	tool, _, err := Tool(ctx)

	if err != nil {
		return result, err
	}

	listArgs := []string{
		"ps",
		"--no-trunc",
		"--format", "{{json .}}",
	}

	if options.All {
		listArgs = append(listArgs, "--all")
	}

	for _, filter := range options.Filter {
		listArgs = append(listArgs, "--filter", filter)
	}

	list := exec.CommandContext(ctx, tool, listArgs...)

	out, err := list.CombinedOutput()

	if err != nil {
		return result, errors.New(string(out))
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))

	type entry struct {
		ID string `json:"ID"`

		Image string `json:"Image"`

		// Command      string `json:"Command"`
		// CreatedAt    string `json:"CreatedAt"`

		Labels string `json:"Labels"`
		// LocalVolumes string `json:"LocalVolumes"`
		// Mounts       string `json:"Mounts"`
		Names string `json:"Names"`
		// Networks     string `json:"Networks"`
		// Ports        string `json:"Ports"`
		// RunningFor   string `json:"RunningFor"`
		// Size         string `json:"Size"`
		// State        string `json:"State"`
		// Status       string `json:"Status"`
	}

	for scanner.Scan() {
		var e entry

		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return result, err
		}

		names := strings.Split(e.Names, ",")

		labels := map[string]string{}

		for _, l := range strings.Split(e.Labels, ",") {
			if key, value, ok := strings.Cut(l, "="); ok {
				labels[key] = value
			}
		}

		c := Container{
			ID: e.ID,

			Names:  names,
			Labels: labels,

			Image: e.Image,
		}

		result = append(result, c)
	}

	return result, nil
}

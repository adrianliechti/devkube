package eksctl

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
)

func List(ctx context.Context, opt ...Option) ([]string, error) {
	output := new(bytes.Buffer)

	e := New(opt...)
	e.stdout = output

	type listType []struct {
		Name   string `json:"Name"`
		Region string `json:"Region"`
		Owned  string `json:"Owned"`
	}

	args := []string{
		"get", "cluster",

		"--region", e.region,
		"--output", "json",
	}

	if err := e.Invoke(ctx, args...); err != nil {
		return nil, err
	}

	var list listType

	if err := json.NewDecoder(output).Decode(&list); err != nil {
		return nil, err
	}

	names := make([]string, 0)

	for _, cluster := range list {
		if !strings.EqualFold(cluster.Owned, "true") {
			continue
		}

		names = append(names, cluster.Name)
	}

	return names, nil
}

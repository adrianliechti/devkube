package gatekeeper

import (
	"context"
	"embed"
	"errors"

	"github.com/adrianliechti/devkube/pkg/apply"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	namespace = "gatekeeper-system"

	version = "v3.16.3"
)

var (
	//go:embed policies/*
	policies embed.FS
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := apply.ApplyURL(ctx, client, namespace, "https://raw.githubusercontent.com/open-policy-agent/gatekeeper/"+version+"/deploy/gatekeeper.yaml"); err != nil {
		return err
	}

	var result error

	entries, _ := policies.ReadDir("policies")

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		f, err := policies.Open("policies/" + e.Name())

		if err != nil {
			result = errors.Join(result, err)
			continue
		}

		defer f.Close()

		if err := apply.Apply(ctx, client, namespace, f); err != nil {
			result = errors.Join(result, err)
			continue
		}
	}

	return result
}

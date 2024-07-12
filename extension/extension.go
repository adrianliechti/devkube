package extension

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	"github.com/adrianliechti/devkube/extension/argocd"
	"github.com/adrianliechti/devkube/extension/tekton"
)

type Extension struct {
	Name string

	Title  string
	Ensure EnsureFunc
}

type EnsureFunc = func(ctx context.Context, client kubernetes.Client) error

var Optional []Extension = []Extension{
	{
		Name:   "argocd",
		Title:  "Argo CD",
		Ensure: argocd.Ensure,
	},
	{
		Name:   "tekton",
		Title:  "Tekton",
		Ensure: tekton.Ensure,
	},
}

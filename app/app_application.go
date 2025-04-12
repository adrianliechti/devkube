package app

import (
	"context"
	"errors"
	"strings"

	"github.com/adrianliechti/go-cli"
	"github.com/adrianliechti/loop/pkg/kubernetes"
	"github.com/adrianliechti/loop/pkg/kubernetes/resource"
)

func SelectApplication(ctx context.Context, client kubernetes.Client, namespace string) (*resource.Application, error) {
	apps, err := resource.Apps(ctx, client, namespace)

	if err != nil {
		return nil, err
	}

	var items []string

	if err != nil {
		return nil, err
	}

	for _, i := range apps {
		items = append(items, strings.Join([]string{i.Namespace, i.Name}, "/"))
	}

	if len(items) == 0 {
		return nil, errors.New("no application(s) found")
	}

	i, _, err := cli.Select("select application", items)

	if err != nil {
		return nil, err
	}

	app := apps[i]
	return &app, nil
}

func MustApplication(ctx context.Context, client kubernetes.Client, namespace string) resource.Application {
	app, err := SelectApplication(ctx, client, namespace)

	if err != nil {
		cli.Fatal(err)
	}

	return *app
}

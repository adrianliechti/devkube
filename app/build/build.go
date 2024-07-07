package build

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/loop/pkg/kubernetes"
	"github.com/adrianliechti/loop/pkg/to"

	"github.com/google/uuid"
	"github.com/moby/moby/pkg/archive"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build image and copy to registry",

		Action: func(c *cli.Context) error {
			client := app.MustClient(c)

			if c.Args().Len() != 2 {
				return errors.New("needs two arguments: image and context path")
			}

			image := "registry.default/" + c.Args().Get(0)

			path, err := filepath.Abs(c.Args().Get(1))

			if err != nil {
				return err
			}

			return Run(c.Context, client, "", image, path, "")
		},
	}
}

func Run(ctx context.Context, client kubernetes.Client, namespace, image, dir, dockerfile string) error {
	if namespace == "" {
		namespace = client.Namespace()
	}

	if dir == "" || dir == "." {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		dir = wd
	}

	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}

	f, err := archive.TarWithOptions(dir, &archive.TarOptions{})

	if err != nil {
		return err
	}

	name := "loop-buildkit-" + uuid.New().String()[0:7]

	cli.Infof("Starting BuildKit pod (%s/%s)...", namespace, name)
	pod, err := startPod(ctx, client, namespace, name, "")

	if err != nil {
		return err
	}

	defer func() {
		cli.Infof("Stopping BuildKit pod (%s/%s)...", pod.Namespace, pod.Name)
		stopPod(context.Background(), client, pod.Namespace, pod.Name)
	}()

	cli.Infof("Copy Context...")

	sessionPath := "/tmp/" + uuid.New().String()

	if err := client.PodExec(ctx, pod.Namespace, pod.Name, "buildkitd", []string{"mkdir", "-p", sessionPath}, false, nil, io.Discard, io.Discard); err != nil {
		return err
	}

	if err := client.PodExec(ctx, namespace, name, "buildkitd", []string{"tar", "xf", "-", "-C", sessionPath}, false, f, io.Discard, io.Discard); err != nil {
		return err
	}

	build := []string{
		"buildctl",
		"build",

		"--frontend", "dockerfile.v0",

		"--local", "context=" + sessionPath,
		"--local", "dockerfile=" + filepath.Dir(filepath.Join(sessionPath, dockerfile)),

		"--output", "type=image,push=true,registry.insecure=true,name=" + image,
	}

	if err := client.PodExec(ctx, namespace, name, "buildkitd", build, false, f, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func startPod(ctx context.Context, client kubernetes.Client, namespace, name, image string) (*corev1.Pod, error) {
	if image == "" {
		image = "moby/buildkit"
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "buildkitd",

					Image:           image,
					ImagePullPolicy: corev1.PullAlways,

					SecurityContext: &corev1.SecurityContext{
						Privileged: to.Ptr(true),
					},

					Args: []string{
						"--addr",
						"unix:///run/buildkit/buildkitd.sock",
					},

					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"buildctl",
									"debug",
									"workers",
								},
							},
						},

						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},

					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"buildctl",
									"debug",
									"workers",
								},
							},
						},

						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
				},
			},

			TerminationGracePeriodSeconds: to.Ptr(int64(10)),
		},
	}

	if _, err := client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return client.WaitForPod(ctx, namespace, name)
}

func stopPod(ctx context.Context, client kubernetes.Client, namespace, name string) error {
	return client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{
		GracePeriodSeconds: to.Ptr(int64(0)),
	})
}

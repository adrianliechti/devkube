package kubectl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/Masterminds/semver"
)

var (
	minimalVersion = semver.MustParse("1.19.0")

	errNotFound = errors.New("kubectl not found. see https://kubernetes.io/docs/tasks/tools/install-kubectl")
	errOutdated = errors.New("kubectl is outdated. see https://kubernetes.io/docs/tasks/tools/install-kubectl")
)

func Info(ctx context.Context) (string, *semver.Version, error) {
	return path(ctx)
}

func path(ctx context.Context) (string, *semver.Version, error) {
	name := "kubectl"

	if runtime.GOOS == "windows" {
		name = "kubectl.exe"
	}

	// verify global tool
	if path, err := exec.LookPath(name); err == nil {
		if version, err := version(ctx, path); err == nil {
			if !version.LessThan(minimalVersion) {
				return path, version, nil
			}

			return path, version, errOutdated
		}

		return path, nil, errOutdated
	}

	return "", nil, errNotFound
}

func version(ctx context.Context, path string) (*semver.Version, error) {
	type versionType struct {
		ClientVersion struct {
			BuildDate    string `json:"buildDate"`
			Compiler     string `json:"compiler"`
			GitCommit    string `json:"gitCommit"`
			GitTreeState string `json:"gitTreeState"`
			GitVersion   string `json:"gitVersion"`
			GoVersion    string `json:"goVersion"`
			Major        string `json:"major"`
			Minor        string `json:"minor"`
			Platform     string `json:"platform"`
		} `json:"clientVersion"`
	}

	cmd := exec.CommandContext(ctx, path, "version", "--client", "-o", "json")
	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var version versionType

	if err := json.Unmarshal(data, &version); err != nil {
		return nil, err
	}

	return semver.NewVersion(version.ClientVersion.GitVersion)
}

type Option func(h *Kubectl)

type Kubectl struct {
	kubeconfig string

	context   string
	namespace string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func New(options ...Option) *Kubectl {
	k := &Kubectl{}

	for _, option := range options {
		option(k)
	}

	return k
}

func WithKubeconfig(kubeconfig string) Option {
	return func(k *Kubectl) {
		k.kubeconfig = kubeconfig
	}
}

func WithContext(context string) Option {
	return func(k *Kubectl) {
		k.context = context
	}
}

func WithNamespace(namespace string) Option {
	return func(k *Kubectl) {
		k.namespace = namespace
	}
}

func WithOutput(stdout, stderr io.Writer) Option {
	return func(k *Kubectl) {
		k.stdout = stdout
		k.stderr = stderr
	}
}

func WithDefaultOutput() Option {
	return WithOutput(os.Stdout, os.Stderr)
}

func (k *Kubectl) Invoke(ctx context.Context, arg ...string) error {
	path, _, err := Info(ctx)

	if err != nil {
		return err
	}

	if k.kubeconfig != "" {
		arg = append(arg, "--kubeconfig", k.kubeconfig)
	}

	if k.context != "" {
		arg = append(arg, "--context", k.context)
	}

	if k.namespace != "" {
		arg = append(arg, "--namespace", k.namespace)
	}

	cmd := exec.CommandContext(ctx, path, arg...)
	cmd.Stdin = k.stdin
	cmd.Stdout = k.stdout
	cmd.Stderr = k.stderr

	return cmd.Run()
}

func Invoke(ctx context.Context, args []string, opt ...Option) error {
	k := New(opt...)

	return k.Invoke(ctx, args...)
}

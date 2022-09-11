package eksctl

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
	minimalVersion = semver.MustParse("0.111.0")

	errNotFound = errors.New("eksctl not found. see https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html")
	errOutdated = errors.New("eksctl is outdated. see https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html")
)

func Info(ctx context.Context) (string, *semver.Version, error) {
	return path(ctx)
}

func path(ctx context.Context) (string, *semver.Version, error) {
	name := "eksctl"

	if runtime.GOOS == "windows" {
		name = "eksctl.exe"
	}

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
		Version string `json:"Version"`
		// PreReleaseID string `json:"PreReleaseID"`
		// Metadata     struct {
		// 	BuildDate time.Time `json:"BuildDate"`
		// 	GitCommit string    `json:"GitCommit"`
		// } `json:"Metadata"`
		// EKSServerSupportedVersions []string `json:"EKSServerSupportedVersions"`
	}

	cmd := exec.CommandContext(ctx, path, "version", "-o", "json")
	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var version versionType

	if err := json.Unmarshal(data, &version); err != nil {
		return nil, err
	}

	return semver.NewVersion(version.Version)
}

type Option func(h *Eksctl)

type Eksctl struct {
	accessKey    string
	accessSecret string

	region string

	stdout io.Writer
	stderr io.Writer
}

func New(options ...Option) *Eksctl {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	accessSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}

	if region == "" {
		region = "eu-west-2"
	}

	k := &Eksctl{
		accessKey:    accessKey,
		accessSecret: accessSecret,

		region: region,
	}

	for _, option := range options {
		option(k)
	}

	return k
}

func WithRegion(region string) Option {
	return func(e *Eksctl) {
		e.region = region
	}
}

func WithOutput(stdout, stderr io.Writer) Option {
	return func(k *Eksctl) {
		k.stdout = stdout
		k.stderr = stderr
	}
}

func WithDefaultOutput() Option {
	return WithOutput(os.Stdout, os.Stderr)
}

func (e *Eksctl) Invoke(ctx context.Context, arg ...string) error {
	path, _, err := Info(ctx)

	if err != nil {
		return err
	}

	env := os.Environ()

	if e.region != "" {
		env = append(env, "AWS_REGION="+e.region)
	}

	if e.accessKey != "" {
		env = append(env, "AWS_ACCESS_KEY_ID="+e.accessKey)
	}

	if e.accessSecret != "" {
		env = append(env, "AWS_SECRET_ACCESS_KEY="+e.accessSecret)
	}

	cmd := exec.CommandContext(ctx, path, arg...)
	cmd.Env = env
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr

	return cmd.Run()
}

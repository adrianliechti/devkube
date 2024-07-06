package helm

import (
	"errors"
	"net/http"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func loadChart(repoURL, chartName, chartVersion string) (*chart.Chart, error) {
	settings := cli.New()

	path, err := repo.FindChartInRepoURL(repoURL, chartName, chartVersion, "", "", "", getter.All(settings))

	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		return loader.LoadArchive(resp.Body)
	}

	//return loader.Load(path)
	return nil, errors.ErrUnsupported
}

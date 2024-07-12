package kubeconfig

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Merge(config ...[]byte) ([]byte, error) {
	result := api.NewConfig()

	for _, data := range config {
		c, err := clientcmd.Load(data)

		if err != nil {
			return nil, err
		}

		for key, cluster := range c.Clusters {
			result.Clusters[key] = cluster
		}

		for key, authInfo := range c.AuthInfos {
			result.AuthInfos[key] = authInfo
		}

		for key, context := range c.Contexts {
			result.Contexts[key] = context
		}

		result.CurrentContext = c.CurrentContext
	}

	return clientcmd.Write(*result)
}

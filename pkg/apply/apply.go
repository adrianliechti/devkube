package apply

import (
	"context"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/adrianliechti/loop/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/restmapper"
)

func Apply(ctx context.Context, client kubernetes.Client, namespace string, reader io.Reader) error {
	if namespace == "" {
		namespace = client.Namespace()
	}

	data, err := io.ReadAll(reader)

	if err != nil {
		return err
	}

	docs, err := splitDocuments(data)

	if err != nil {
		return err
	}

	discovery, err := discovery.NewDiscoveryClientForConfig(client.Config())

	if err != nil {
		return err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discovery))

	serializer := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	for _, doc := range docs {
		obj := &unstructured.Unstructured{}
		_, gvk, err := serializer.Decode(doc, nil, obj)

		if err != nil {
			if runtime.IsMissingKind(err) {
				continue
			}

			return err
		}

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)

		if err != nil {
			return err
		}

		if obj.GetNamespace() == "" && mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			obj.SetNamespace(namespace)
		}

		i := client.Resource(mapping.Resource).Namespace(obj.GetNamespace())

		if _, err := i.Apply(ctx, obj.GetName(), obj, metav1.ApplyOptions{
			Force:        true,
			FieldManager: "devkube",
		}); err != nil {
			return err
		}
	}

	return nil
}

func ApplyFile(ctx context.Context, client kubernetes.Client, namespace string, path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer f.Close()

	return Apply(ctx, client, namespace, f)
}

func ApplyURL(ctx context.Context, client kubernetes.Client, namespace string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return Apply(ctx, client, namespace, resp.Body)
}

func splitDocuments(data []byte) ([][]byte, error) {
	re := regexp.MustCompile(`(?m)^---$`)

	var results [][]byte

	for _, doc := range re.Split(string(data), -1) {
		results = append(results, []byte(doc))
	}

	return results, nil
}

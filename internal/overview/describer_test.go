package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestListDescriber(t *testing.T) {
	thePath := "/"
	key := cache.Key{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	c := newSpyCache()

	object := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name":              "name",
			"creationTimestamp": "2019-01-14T13:34:56+00:00",
		},
	}

	retrieveKey := cache.Key{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}
	c.spyRetrieve(retrieveKey, []*unstructured.Unstructured{{Object: object}}, nil)

	listType := func() interface{} {
		return &corev1.PodList{}
	}

	objectType := func() interface{} {
		return &corev1.Pod{}
	}

	d := NewListDescriber(thePath, "list", key, listType, objectType)

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:   c,
		Fields:  fields,
		Printer: printer.NewResource(c),
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	list := component.NewList("list", nil)

	tableCols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable("/v1, Kind=PodList", tableCols)
	table.Add(component.TableRow{
		"Age":    component.NewTimestamp(time.Unix(1547472896, 0)),
		"Labels": component.NewLabels(nil),
		"Name":   component.NewText("name"),
	})
	list.Add(table)

	expected := component.ContentResponse{
		ViewComponents: []component.ViewComponent{list},
	}

	assert.Equal(t, expected, cResponse)

	assert.True(t, c.isSatisfied(), "cache was not satisfied")
}

func TestObjectDescriber(t *testing.T) {
	thePath := "/"
	key := cache.Key{APIVersion: "v1", Kind: "Pod"}
	namespace := "default"
	fields := map[string]string{}

	c := newSpyCache()

	object := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name":      "pod",
			"namespace": "default",
		},
		"spec": map[string]interface{}{
			"containers": []interface{}{
				map[string]interface{}{
					"name": "one",
				},
				map[string]interface{}{
					"name": "two",
				},
			},
		},
	}

	retrieveKey := cache.Key{Namespace: namespace, APIVersion: "v1", Kind: "Pod"}
	c.spyRetrieve(retrieveKey, []*unstructured.Unstructured{{Object: object}}, nil)

	objectType := func() interface{} {
		return &corev1.Pod{}
	}

	fn := DefaultLoader(key)

	d := NewObjectDescriber(thePath, "object", fn, objectType, true)

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	p := printer.NewResource(c)
	err = p.Handler(func(*corev1.Pod, printer.Options) (component.ViewComponent, error) {
		return component.NewText("*v1.Pod"), nil
	})
	require.NoError(t, err)

	options := DescriberOptions{
		Cache:   c,
		Fields:  fields,
		Printer: p,
	}

	ctx := context.Background()
	cResponse, err := d.Describe(ctx, "/path", namespace, clusterClient, options)
	require.NoError(t, err)

	summary := component.NewText("*v1.Pod")
	summary.SetAccessor("summary")

	yaml := component.NewYAML(component.TitleFromString("YAML"),
		"---\napiVersion: v1\nkind: Pod\nmetadata:\n  creationTimestamp: null\n  name: pod\n  namespace: default\nspec:\n  containers:\n  - name: one\n    resources: {}\n  - name: two\n    resources: {}\nstatus: {}\n")
	yaml.SetAccessor("yaml")

	logs := component.NewLogs("default", "pod", []string{"one", "two"})
	logs.SetAccessor("logs")

	expected := component.ContentResponse{
		Title: component.Title(component.NewText("object"), component.NewText("pod")),
		ViewComponents: []component.ViewComponent{
			summary,
			yaml,
			logs,
		},
	}
	assert.Equal(t, expected, cResponse)
	assert.True(t, c.isSatisfied())
}

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	c := cache.NewMemoryCache()

	scheme := runtime.NewScheme()
	objects := []runtime.Object{}
	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	options := DescriberOptions{
		Cache: c,
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		d        *SectionDescriber
		expected component.ContentResponse
	}{
		{
			name: "general",
			d: NewSectionDescriber(
				"/section",
				"section",
				newStubDescriber("/foo"),
			),
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("section")),
				ViewComponents: []component.ViewComponent{
					component.NewList("section", nil),
				},
			},
		},
		{
			name: "empty",
			d: NewSectionDescriber(
				"/section",
				"section",
				newEmptyDescriber("/foo"),
				newEmptyDescriber("/bar"),
			),
			expected: component.ContentResponse{
				Title: component.Title(component.NewText("section")),
				ViewComponents: []component.ViewComponent{
					component.NewList("section", nil),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got, err := tc.d.Describe(ctx, "/prefix", namespace, clusterClient, options)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

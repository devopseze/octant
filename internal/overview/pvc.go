package overview

import (
	"context"
	"github.com/heptio/developer-dash/internal/cache"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type PersistentVolumeClaimSummary struct{}

var _ View = (*PersistentVolumeClaimSummary)(nil)

func NewPersistentVolumeClaimSummary(prefix, namespace string, c clock.Clock) View {
	return &PersistentVolumeClaimSummary{}
}

func (js *PersistentVolumeClaimSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	secret, err := retrievePersistentVolumeClaim(object)
	if err != nil {
		return nil, err
	}

	detail, err := printPersistentVolumeClaimSummary(secret)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

func retrievePersistentVolumeClaim(object runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	rc, ok := object.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil, errors.Errorf("expected object to be a Persistent Volume Claim, it was %T", object)
	}

	return rc, nil
}

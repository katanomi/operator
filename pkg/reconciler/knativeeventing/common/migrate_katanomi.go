/*
Copyright 2022 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/operator/pkg/client/clientset/versioned"
	"knative.dev/pkg/logging"
)

var (
	annotationDisableMigrateAutoUpgrade = "katanomi.dev/disable-migrate-auto-upgrade"
)

// MigrateEventingSpecVersion will set `spec.version` to empty for all KnativeEventings
// if KnativeEventing has an annotation katanomi.dev/disable-migrate-auto-upgrade, it will just skip migration
func MigrateEventingSpecVersion(ctx context.Context, operatorClient versioned.Interface) error {
	logger := logging.FromContext(ctx).With("name", "migrate-eventing-specversion")
	logger.Info("Migrating the existing KnativeEventing instances")
	eventings, err := operatorClient.OperatorV1beta1().KnativeEventings("").List(ctx, metav1.ListOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		return err
	}
	if len(eventings.Items) == 0 {
		logger.Debugw("there is no knative eventings need to migrate spec.version")
		return nil
	}

	for _, eventing := range eventings.Items {
		if eventing.Annotations != nil {
			if _, ok := eventing.Annotations[annotationDisableMigrateAutoUpgrade]; ok {
				logger.Debugw("skip eventing", "eventing", eventing.Namespace+"/"+eventing.Name)
				continue
			}
		}
		if eventing.Spec.Version == "" {
			continue
		}

		patch := `{"spec":{"version": ""}}`
		_, err = operatorClient.OperatorV1beta1().KnativeEventings(eventing.Namespace).Patch(
			ctx, eventing.Name, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
		if err != nil {
			logger.Errorw("error to update spec.version to empty", "err", err)
			continue
		}
		logger.Infow("migrated spec.version to empty", "eventing", eventing.Namespace+"/"+eventing.Name)
	}

	return nil
}

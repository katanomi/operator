/*
Copyright 2020 The Knative Authors
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

	mf "github.com/manifestival/manifestival"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"knative.dev/operator/pkg/apis/operator/v1alpha1"
	"knative.dev/pkg/logging"
)

// transformers that are common to all components.
func transformers(ctx context.Context, obj v1alpha1.KComponent) []mf.Transformer {
	logger := logging.FromContext(ctx)
	return []mf.Transformer{
		injectOwner(obj),
		mf.InjectNamespace(obj.GetNamespace()),
		JobTransform(obj),
		HighAvailabilityTransform(obj, logger),
		ImageTransform(obj.GetSpec().GetRegistry(), logger),
		ConfigMapTransform(obj.GetSpec().GetConfig(), logger),
		ResourceRequirementsTransform(obj, logger),
		DeploymentsTransform(obj, logger),
	}
}

func injectOwner(owner mf.Owner) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetNamespace() != "" {
			u.SetOwnerReferences([]v1.OwnerReference{*v1.NewControllerRef(owner, owner.GroupVersionKind())})
		}
		return nil
	}
}

// Transform will mutate the passed-by-reference manifest with one
// transformed by platform, common, and any extra passed in
func Transform(ctx context.Context, manifest *mf.Manifest, instance v1alpha1.KComponent, extra ...mf.Transformer) error {
	logger := logging.FromContext(ctx)
	logger.Debug("Transforming manifest")

	transformers := transformers(ctx, instance)
	transformers = append(transformers, extra...)

	m, err := manifest.Transform(transformers...)
	if err != nil {
		instance.GetStatus().MarkInstallFailed(err.Error())
		return err
	}
	*manifest = m
	return nil
}

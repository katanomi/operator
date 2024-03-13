package common

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/operator/pkg/apis/operator/base"
	"knative.dev/operator/pkg/apis/operator/v1beta1"
	fakeclientset "knative.dev/operator/pkg/client/clientset/versioned/fake"
	"knative.dev/pkg/logging"
	"testing"
)

func TestMigrateEventingSpecVersion(t *testing.T) {

	var tests = []struct {
		desc      string
		eventings []runtime.Object

		expectedSpecVersion map[string]string
	}{
		{
			desc:                "no eventings is cluster",
			expectedSpecVersion: map[string]string{},
		},

		{
			desc: "only one eventings",
			eventings: []runtime.Object{
				&v1beta1.KnativeEventing{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "default",
						Namespace: "ns1",
					},
					Spec: v1beta1.KnativeEventingSpec{
						CommonSpec: base.CommonSpec{
							Version: "v1.0",
						},
					},
				},
			},
			expectedSpecVersion: map[string]string{
				"ns1/default": "",
			},
		},

		{
			desc: "more than one eventings",
			eventings: []runtime.Object{
				&v1beta1.KnativeEventing{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-empty",
						Namespace: "ns1",
					},
					Spec: v1beta1.KnativeEventingSpec{
						CommonSpec: base.CommonSpec{
							Version: "v1.0",
						},
					},
				},
				&v1beta1.KnativeEventing{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "empty",
						Namespace: "ns2",
					},
					Spec: v1beta1.KnativeEventingSpec{
						CommonSpec: base.CommonSpec{
							Version: "",
						},
					},
				},
				&v1beta1.KnativeEventing{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-empty",
						Namespace: "ns2",
					},
					Spec: v1beta1.KnativeEventingSpec{
						CommonSpec: base.CommonSpec{
							Version: "v2.0",
						},
					},
				},
				&v1beta1.KnativeEventing{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "empty",
						Namespace: "skip-migrate-ns",
						Annotations: map[string]string{
							annotationDisableMigrateAutoUpgrade: "true",
						},
					},
					Spec: v1beta1.KnativeEventingSpec{
						CommonSpec: base.CommonSpec{
							Version: "",
						},
					},
				},
				&v1beta1.KnativeEventing{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-empty",
						Namespace: "skip-migrate-ns",
						Annotations: map[string]string{
							annotationDisableMigrateAutoUpgrade: "true",
						},
					},
					Spec: v1beta1.KnativeEventingSpec{
						CommonSpec: base.CommonSpec{
							Version: "v1.0",
						},
					},
				},
			},
			expectedSpecVersion: map[string]string{
				"ns1/not-empty":             "",
				"ns2/empty":                 "",
				"ns2/not-empty":             "",
				"skip-migrate-ns/empty":     "",
				"skip-migrate-ns/not-empty": "v1.0",
			},
		},
	}

	for _, item := range tests {

		t.Run(item.desc, func(t *testing.T) {
			ctx := context.Background()
			ctx = logging.WithLogger(ctx, log)

			objs := []runtime.Object{}
			if item.eventings != nil {
				objs = append(objs, item.eventings...)
			}
			client := fakeclientset.NewSimpleClientset(objs...)

			err := MigrateEventingSpecVersion(ctx, client)
			if err != nil {
				t.Errorf("should migrate succeed: %s", err.Error())
			}

			all, err := client.OperatorV1beta1().KnativeEventings("").List(ctx, metav1.ListOptions{ResourceVersion: "0"})
			if err != nil {
				t.Errorf("list eventings error: %s", err.Error())
			}

			actualSpecVersion := map[string]string{}
			for _, eventing := range all.Items {
				actualSpecVersion[eventing.Namespace+"/"+eventing.Name] = eventing.Spec.Version
			}

			if len(actualSpecVersion) != len(item.expectedSpecVersion) {
				t.Errorf("length should be equal, actual: %d, expected: %d", len(actualSpecVersion), len(item.expectedSpecVersion))
			}
			for key, actualValue := range actualSpecVersion {
				var ok bool
				actualValue, ok = item.expectedSpecVersion[key]
				if !ok {
					t.Errorf("should not contains knativeeventing `%s`", key)
				}

				if actualValue != item.expectedSpecVersion[key] {
					t.Errorf("expect `spec.version` of knativeeventing '%s' should be '%s' but got '%s' ", key, item.expectedSpecVersion[key], actualValue)
				}
			}
		})
	}
}

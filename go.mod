module knative.dev/operator

go 1.14

require (
	github.com/emicklei/go-restful v2.11.1+incompatible // indirect
	github.com/go-logr/zapr v0.1.1
	github.com/google/go-cmp v0.5.4
	github.com/google/go-github/v32 v32.1.0
	github.com/manifestival/client-go-client v0.4.0
	github.com/manifestival/manifestival v0.6.1
	go.uber.org/zap v1.16.0
	golang.org/x/mod v0.3.0
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.18.12
	knative.dev/caching v0.0.0-20210120064853-cc4d11bd4274
	knative.dev/eventing v0.20.1-0.20210120133253-5dc363dcb5be
	knative.dev/hack v0.0.0-20210114150620-4422dcadb3c8
	knative.dev/pkg v0.0.0-20210119162123-1bbf0a6436c3
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/code-generator => k8s.io/code-generator v0.18.8
)

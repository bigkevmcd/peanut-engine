module github.com/bigkevmcd/peanut-engine

go 1.17

require (
	github.com/argoproj/gitops-engine v0.4.2
	github.com/argoproj/pkg v0.11.0
	github.com/bigkevmcd/peanut v0.0.0-20211103065456-91b7d2995f2d
	github.com/go-git/go-git/v5 v5.4.2
	github.com/google/go-cmp v0.5.6
	github.com/julienschmidt/httprouter v1.3.0
	github.com/manifestival/manifestival v0.7.0
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	k8s.io/apimachinery v0.21.5
	k8s.io/client-go v11.0.0+incompatible
	knative.dev/pkg v0.0.0-20211129195804-438776b3c87c
	sigs.k8s.io/kustomize/api v0.10.0
)

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.6.0

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20170808103936-bb23615498cd // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/chai2010/gettext-go v0.0.0-20160711120539-c6fed771bfd5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20210610120745-9d4ed1856297 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/russross/blackfriday v2.0.0+incompatible // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	golang.org/x/crypto v0.0.0-20210920023735-84f357641f63 // indirect
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3 // indirect
	golang.org/x/oauth2 v0.0.0-20211028175245-ba495a64dcb5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.21.5 // indirect
	k8s.io/apiextensions-apiserver v0.21.5 // indirect
	k8s.io/apiserver v0.21.5 // indirect
	k8s.io/cli-runtime v0.21.5 // indirect
	k8s.io/component-base v0.21.5 // indirect
	k8s.io/component-helpers v0.21.5 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/kube-aggregator v0.21.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e // indirect
	k8s.io/kubectl v0.21.0 // indirect
	k8s.io/kubernetes v1.21.0 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.4.2
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway => github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/improbable-eng/grpc-web => github.com/improbable-eng/grpc-web v0.0.0-20181111100011-16092bd1d58a
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

	google.golang.org/grpc => google.golang.org/grpc v1.15.0

	k8s.io/api => k8s.io/api v0.21.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.5
	k8s.io/apiserver => k8s.io/apiserver v0.21.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.5
	k8s.io/client-go => k8s.io/client-go v0.21.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.5
	k8s.io/code-generator => k8s.io/code-generator v0.21.5
	k8s.io/component-base => k8s.io/component-base v0.21.5
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.5
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.5
	k8s.io/cri-api => k8s.io/cri-api v0.21.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.5
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.5
	k8s.io/kubectl => k8s.io/kubectl v0.21.5
	k8s.io/kubelet => k8s.io/kubelet v0.21.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.5
	k8s.io/metrics => k8s.io/metrics v0.21.5
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.5
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.21.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.5
)

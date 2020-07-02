module github.com/bigkevmcd/peanut-engine

go 1.14

require (
	github.com/argoproj/gitops-engine v0.1.3-0.20200620112536-6657adfcfde4
	github.com/argoproj/pkg v0.0.0-20200424003221-9b858eff18a1
	github.com/bigkevmcd/peanut v0.0.0-20200617060744-b104905993ee
	github.com/go-git/go-billy/v5 v5.0.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/google/go-cmp v0.5.0
	github.com/jenkins-x/go-scm v1.5.79
	github.com/julienschmidt/httprouter v1.2.0
	github.com/manifestival/manifestival v0.5.0
	github.com/prometheus/client_golang v1.5.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190816222228-6d55c1b1f1ca+incompatible
	knative.dev/pkg v0.0.0-20200616232624-ffb929374a39
	sigs.k8s.io/kustomize v2.0.3+incompatible
)

replace (
	// https://github.com/kubernetes/kubernetes/issues/79384#issuecomment-505627280
	k8s.io/api => k8s.io/api v0.16.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.6 // indirect
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.6 // indirect
	k8s.io/apiserver => k8s.io/apiserver v0.16.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.6
	k8s.io/client-go => k8s.io/client-go v0.16.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.6
	k8s.io/code-generator => k8s.io/code-generator v0.16.6
	k8s.io/component-base => k8s.io/component-base v0.16.6
	k8s.io/cri-api => k8s.io/cri-api v0.16.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.6
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.6
	k8s.io/kubectl => k8s.io/kubectl v0.16.6
	k8s.io/kubelet => k8s.io/kubelet v0.16.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.6
	k8s.io/metrics => k8s.io/metrics v0.16.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.6
)

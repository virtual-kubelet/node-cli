module github.com/virtual-kubelet/node-cli

go 1.19

require (
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/virtual-kubelet/virtual-kubelet v1.8.0
	go.opencensus.io v0.23.0
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.26.2
	k8s.io/apimachinery v0.26.2
	k8s.io/apiserver v0.25.0
	k8s.io/client-go v0.25.0
	k8s.io/klog v1.0.0
)

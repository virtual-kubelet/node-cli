module github.com/virtual-kubelet/node-cli

go 1.12

require (
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/virtual-kubelet/virtual-kubelet v1.2.1-0.20200320212802-3ec3b14e49d0
	go.opencensus.io v0.21.0
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // indirect
	golang.org/x/text v0.3.1-0.20181227161524-e6919f6577db // indirect
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0-20181213151703-3ccfe8365421 // indirect
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v0.3.1
	k8s.io/kubernetes v1.15.2
	k8s.io/utils v0.0.0-20180801164400-045dc31ee5c4 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190521190702-177766529176

replace k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628

replace k8s.io/kubernetes => k8s.io/kubernetes v1.13.7

package root

import (
	"context"
	"testing"

	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/node-cli/provider/mock"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestRunRootCommand(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := opts.New()
	providerInitFunc := func(cfg provider.InitConfig) (provider.Provider, error) {
		mockConfig := mock.Config{
			CPU:    "1",
			Memory: "128M",
			Pods:   "120",
		}
		return mock.NewProviderConfig(mockConfig, cfg.NodeName, cfg.OperatingSystem, cfg.InternalIP, cfg.DaemonPort)
	}
	fakeClient := fake.NewSimpleClientset()
	errCh := make(chan error)
	go func() {
		errCh <- runRootCommandWithProviderAndClient(ctx, providerInitFunc, fakeClient, opts)
	}()

	watch, err := fakeClient.CoreV1().Nodes().Watch(ctx, metav1.ListOptions{})
	assert.NilError(t, err)
	defer watch.Stop()
	for ev := range watch.ResultChan() {
		node := ev.Object.(*corev1.Node)
		t.Logf("Node registered: %+v", node)
		break
	}
}

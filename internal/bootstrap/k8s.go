package bootstrap

import (
	"fmt"
	"path/filepath"

	appcfg "gpu-scheduler-platform/internal/config"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClients struct {
	RestConfig *rest.Config
	Clientset  kubernetes.Interface
	Dynamic    dynamic.Interface
}

func NewKubernetesClients(cfg appcfg.KubernetesConfig) (*K8sClients, error) {
	restCfg, err := buildK8sRestConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes rest config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes clientset: %w", err)
	}

	dyn, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes dynamic client: %w", err)
	}

	return &K8sClients{
		RestConfig: restCfg,
		Clientset:  clientset,
		Dynamic:    dyn,
	}, nil
}

func buildK8sRestConfig(cfg appcfg.KubernetesConfig) (*rest.Config, error) {
	var (
		restCfg *rest.Config
		err     error
	)

	switch {
	case cfg.InCluster:
		restCfg, err = rest.InClusterConfig()
	default:
		kubeconfig := cfg.Kubeconfig
		if kubeconfig == "" {
			kubeconfig = defaultKubeconfigPath()
		}
		restCfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, err
	}

	restCfg.QPS = float32(cfg.QPS)
	restCfg.Burst = cfg.Burst
	return restCfg, nil
}

func defaultKubeconfigPath() string {
	home, err := homeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

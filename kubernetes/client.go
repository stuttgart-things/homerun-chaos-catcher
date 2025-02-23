/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateKubernetesClient(pathToKubeconfig string) (clientset *kubernetes.Clientset) {

	// LOAD KUBECONFIG FROM DEFAULT LOCATION (~/.KUBE/CONFIG)
	config, err := clientcmd.BuildConfigFromFlags("", pathToKubeconfig)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		return
	}

	// CREATE THE KUBERNETES CLIENTSET
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		return
	}

	return clientset
}

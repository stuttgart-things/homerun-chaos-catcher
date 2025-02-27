/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package kubernetes

import (
	"fmt"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
	"k8s.io/client-go/kubernetes"
)

func CorndonUncordonNode(clientset *kubernetes.Clientset) {
	workers, err := sthingsK8s.GetNodesByRole(clientset, "worker")
	if err != nil {
		fmt.Println("Error getting workers:", err)
		return
	}

	fmt.Println("Found workers:", len(workers))
	fmt.Println(workers)
}

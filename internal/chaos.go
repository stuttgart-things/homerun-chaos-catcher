/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	k8s "github.com/stuttgart-things/homerun-chaos-catcher/kubernetes"

	"k8s.io/client-go/kubernetes"
)

func CreateChaos(kind string, count int, operation string, k8sClient *kubernetes.Clientset) {

	// SWITCH CASE FOR KIND
	fmt.Println("Creating chaos for", kind)
	fmt.Println("Operation", operation)
	fmt.Println("Count", count)

	switch kind {
	case "pod":
		k8s.DeleteRandomPods(k8sClient, "", count)

	case "node":
		fmt.Println("Creating chaos for nodes")
		k8s.CorndonUncordonNode(k8sClient)

	default:
		fmt.Println("Unknown kind")
		return
	}

}

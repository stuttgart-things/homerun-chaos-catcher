/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package kubernetes

import (
	"context"
	"fmt"
	"time"

	"math/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// func GetRandomPods(clientset *kubernetes.Clientset, namespace string, count int) ([]v1.Pod, error) {
// 	// List all pods in the given namespace
// 	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		return nil, fmt.Errorf("could not list pods: %v", err)
// 	}

// 	rand.Seed(uint64(time.Now().UnixNano()))

// 	// Shuffle the pod slice
// 	rand.Shuffle(len(pods.Items), func(i, j int) { pods.Items[i], pods.Items[j] = pods.Items[j], pods.Items[i] })

// 	return pods.Items[:count], nil
// }

func DeleteRandomPods(clientset *kubernetes.Clientset, namespace string, countPods int) error {
	// Get the list of all pods (in a single namespace or across all)
	podList, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	pods := podList.Items

	// Print pod names
	for _, pod := range pods {
		fmt.Printf("Found pod: %s in namespace: %s\n", pod.Name, pod.Namespace)
	}

	// If there are fewer pods than the requested count, adjust the count
	if len(pods) < countPods {
		countPods = len(pods)
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Randomly select pods and delete them
	for i := 0; i < countPods; i++ {
		randomIndex := rand.Intn(len(pods))
		podToDelete := pods[randomIndex]

		// Delete the selected pod using its actual namespace
		err := clientset.CoreV1().Pods(podToDelete.Namespace).Delete(
			context.Background(), podToDelete.Name, metav1.DeleteOptions{},
		)
		if err != nil {
			fmt.Printf("Error deleting pod %s in namespace %s: %v\n", podToDelete.Name, podToDelete.Namespace, err)
		} else {
			fmt.Printf("Successfully deleted pod: %s in namespace: %s\n", podToDelete.Name, podToDelete.Namespace)
		}

		// Remove the deleted pod from the list to avoid re-selection
		pods = append(pods[:randomIndex], pods[randomIndex+1:]...)
	}

	return nil
}

/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package kubernetes

import (
	"context"
	"fmt"
	"time"

	"math/rand"

	"os"

	homerun "github.com/stuttgart-things/homerun-library"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	dt              = time.Now()
	redisConnection = map[string]string{
		"addr":     os.Getenv("REDIS_SERVER"),
		"port":     os.Getenv("REDIS_PORT"),
		"password": os.Getenv("REDIS_PASSWORD"),
		"stream":   os.Getenv("REDIS_STREAM"),
		"group":    os.Getenv("REDIS_CONSUMER_GROUP"),
	}
)

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
			message := homerun.Message{
				Title:           "Test",
				Message:         fmt.Sprintf("Successfully deleted pod: %s in namespace: %s", podToDelete.Name, podToDelete.Namespace),
				Severity:        "INFO",
				Author:          "homerun-chaos-catcher",
				Timestamp:       dt.Format("01-02-2006 15:04:05"),
				System:          "homerun-chaos-catcher",
				Tags:            "chaos, pods",
				AssigneeAddress: "",
				AssigneeName:    "",
				Artifacts:       "",
				Url:             "",
			}

			homerun.EnqueueMessageInRedisStreams(message, redisConnection)

			fmt.Printf("Successfully deleted pod: %s in namespace: %s\n", podToDelete.Name, podToDelete.Namespace)
		}

		// Remove the deleted pod from the list to avoid re-selection
		pods = append(pods[:randomIndex], pods[randomIndex+1:]...)
	}

	return nil
}

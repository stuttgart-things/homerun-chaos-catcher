package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"math/rand"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetPodsInNamespace fetches all pods in the specified namespace.
func GetRandomPods(clientset *kubernetes.Clientset, namespace string, count int) ([]v1.Pod, error) {
	// List all pods in the given namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list pods: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	// Shuffle the pod slice
	rand.Shuffle(len(pods.Items), func(i, j int) { pods.Items[i], pods.Items[j] = pods.Items[j], pods.Items[i] })

	return pods.Items[:count], nil
}

func main() {
	namespace := "" // Replace with the desired namespace

	clientset := CreateKubernetesClient("/home/sthings/.kube/config")

	DeleteRandomPods(clientset, namespace, 4)

	// pods, err := GetRandomPods(clientset, namespace, 2)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // LOOP OVER PODS AND PRINT NAMES
	// for _, pod := range pods {
	// 	fmt.Println(pod.Name)
	// }

	// nodes, err := GetNodes(clientset, namespace, 4)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println("NOODES", nodes)

}

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

func GetNodes(clientset *kubernetes.Clientset, namespace string, count int) ([]v1.Node, error) {

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing nodes: %v", err)
	}

	fmt.Println("Nodes in cluster:")
	for _, node := range nodes.Items {
		fmt.Println("- ", node.Name)
	}

	return nodes.Items, nil

}

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

// GetPodsInNamespace fetches all pods in the specified namespace.
func GetPodsInNamespace(clientset *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	// List all pods in the given namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list pods: %v", err)
	}
	return pods.Items, nil
}

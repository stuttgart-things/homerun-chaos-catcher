package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Specify your kubeconfig file path
	kubeconfig := filepath.Join("/home/sthings/.kube", "config")

	// Load the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Example: List nodes in the cluster
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing nodes: %v", err)
	}

	for _, node := range nodes.Items {
		fmt.Println("Node Name:", node.Name)
	}

	// cordonNode(clientset, "kind-worker2")
	// evictPods(clientset, "kind-worker2")
	restartAllDeployments(clientset)
	// uncordonNode
}

// cordonNode marks the node as unschedulable
func cordonNode(clientset *kubernetes.Clientset, nodeName string) error {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	node.Spec.Unschedulable = true // Mark node as unschedulable

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Node %s is now cordoned\n", nodeName)
	return nil
}

// evictPods evicts all non-DaemonSet and non-mirror pods from the node
func evictPods(clientset *kubernetes.Clientset, nodeName string) error {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		// Skip DaemonSet pods
		if isDaemonSetPod(pod) {
			continue
		}

		// Create eviction
		eviction := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			DeleteOptions: &metav1.DeleteOptions{
				GracePeriodSeconds: int64Ptr(30), // Set a grace period
			},
		}

		err := clientset.PolicyV1().Evictions(eviction.Namespace).Evict(context.TODO(), eviction)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			log.Printf("Failed to evict pod %s: %v\n", pod.Name, err)
		} else {
			fmt.Printf("Evicted pod: %s/%s\n", pod.Namespace, pod.Name)
		}

		time.Sleep(1 * time.Second) // Avoid API throttling
	}

	return nil
}

// isDaemonSetPod checks if a pod is managed by a DaemonSet
func isDaemonSetPod(pod v1.Pod) bool {
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "DaemonSet" {
			return true
		}
	}
	return false
}

// Helper function to get pointer to int64
func int64Ptr(i int64) *int64 {
	return &i
}

func uncordonNode(clientset *kubernetes.Clientset, nodeName string) error {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	node.Spec.Unschedulable = false // Mark node as schedulable

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Node %s is now uncordoned (schedulable again)\n", nodeName)
	return nil
}

// restartAllDeployments restarts all deployments in every namespace
func restartAllDeployments(clientset *kubernetes.Clientset) error {
	// Get all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ns := range namespaces.Items {
		namespace := ns.Name
		fmt.Printf("Restarting deployments in namespace: %s\n", namespace)

		// Get all deployments in the namespace
		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error getting deployments in namespace %s: %v", namespace, err)
			continue
		}

		for _, deploy := range deployments.Items {
			fmt.Printf("Restarting deployment: %s/%s\n", namespace, deploy.Name)

			// Add an annotation to trigger a rolling restart
			if deploy.Spec.Template.Annotations == nil {
				deploy.Spec.Template.Annotations = make(map[string]string)
			}
			deploy.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

			_, err := clientset.AppsV1().Deployments(namespace).Update(context.TODO(), &deploy, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("Failed to restart deployment %s/%s: %v", namespace, deploy.Name, err)
			}
		}
	}

	return nil
}

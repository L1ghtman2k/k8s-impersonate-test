package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

func main() {
	// Path to the service account token.
	// Usually found in /var/run/secrets/kubernetes.io/serviceaccount/token
	serviceAccountTokenFile := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	// Read the service account token.
	token, err := os.ReadFile(serviceAccountTokenFile)
	if err != nil {
		panic(err)
	}

	// Configure the Kubernetes client using the service account token.
	config := &rest.Config{
		Host:        "https://kubernetes.default.svc.cluster.local",
		BearerToken: string(token),
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // Set this to true if using an insecure connection (not recommended)
		},
	}

	config.Impersonate = rest.ImpersonationConfig{
		UserName: "system:serviceaccount:service-team-1:sa-with-limited-role",
	}

	// Create the Kubernetes client.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// List all namespaces.
	pods, err := clientset.CoreV1().Pods("service-team-1").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Namespaces:\n")
	for _, namespace := range pods.Items {
		fmt.Printf("- %s\n", namespace.Name)
	}
}

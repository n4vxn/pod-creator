package main

import (
	"context"
	"net/http"
	"os/user"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"                       // k8s pod and related APIs
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // k8s metadata types
	"k8s.io/client-go/kubernetes"                 // k8s client
	"k8s.io/client-go/tools/clientcmd"            // for kubeconfig
)

// CreatePod handles the creation of a pod
func CreatePod(c *gin.Context) {
	usr, _ := user.Current()
	kubeconfig := usr.HomeDir + "/.kube/config"
	
	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load kubeconfig"})
		return
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create kubernetes client"})
		return
	}

	// Define the pod
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-first-pod",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}

	// Create the pod in the default namespace
	_, err = clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pod created successfully"})
}

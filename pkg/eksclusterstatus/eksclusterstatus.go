package eksclusterstatus

import (

	"context"
	"fmt"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func EKSClusterStatus(clientset *kubernetes.Clientset, clusterName, region string) (string, error) {

	var sb strings.Builder

	// api path to check the health of Kubernetes api server
	path := "/healthz"
	content, err := clientset.Discovery().RESTClient().Get().AbsPath(path).DoRaw(context.TODO())
	if err != nil {
		return "", fmt.Errorf("Error in cluster health check: %w", err)
	}
	fmt.Fprintf(&sb, "Cluster %s Health Status : %s\n", clusterName, string(content))

	// Api call to check if all worker nodes are associated with cluster
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("Error getting EKS nodes: %w", err)
	}
	fmt.Fprintf(&sb, "There are %d nodes associated with cluster %s\n", len(nodes.Items), clusterName)

	// Api call to list the number of pods in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	fmt.Fprintf(&sb, "There are %d pods in the cluster\n", len(pods.Items))
    


	// Code to check if a deployment object exists in the cluster - For Future Reference
	// 	deployment := "external-dns"
	// 	namespace = "kube-external-dns"
	// 	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	// 	_, err = deploymentsClient.Get(context.TODO(), deployment, metav1.GetOptions{})
	// 	if err != nil {
	// 		panic(err.Error())
	// 	} else {
	// 		 	fmt.Printf("Found deployment %s in namespace %s\n", deployment, namespace)
	//  }
	
	return sb.String(), nil
}

package podstatus

import (
	
	"context"
	"fmt"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

)

// Function to check status of the pods
func CheckPodsStatus(clientset *kubernetes.Clientset, clusterName, region, name, labelSelector, namespace string) (string, error) {

	var sb strings.Builder
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         100,
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)

	if err != nil {
		  fmt.Fprintf(&sb,"\n#############%s#############\n", name)
		  return "", fmt.Errorf("Error occured when retrieving %s pods in %s namespace\n", name, namespace)
	} else if len(pods.Items) == 0 {
		fmt.Fprintf(&sb,"\n#############%s#############\n", name)
		return "", fmt.Errorf("%s pods with label selector %s not found in %s namespace\n", name, labelSelector, namespace)
	} else {
		fmt.Fprintf(&sb,"\n#############%s#############\n", name)
		for _, pod := range pods.Items {
			fmt.Fprintf(&sb, "%s pod in %s namespace is %v\n", pod.Name, namespace, pod.Status.Phase)		
		}
		fmt.Fprintf(&sb, "\n")
	}
	return sb.String(),nil
}

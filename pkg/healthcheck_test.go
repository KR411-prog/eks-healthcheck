package healthcheck

import (
	"healthcheck/pkg/podstatus"
	"healthcheck/pkg/clientset"
	"healthcheck/pkg/eksclusterstatus"
	"context"
	"fmt"
	"os"

	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

// Test function to check EKS cluster health status, Pods status
func TestEKSComplete(t *testing.T) {

	clusterName, clusterExists := os.LookupEnv("clusterName")
	if !clusterExists {
		t.Skip("clusterName environment variable  not set")
	}
	region, regionExists := os.LookupEnv("region")
	if !regionExists {
		t.Skip("region environment variable  not set")
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	eksSvc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}
	result, err := eksSvc.DescribeCluster(input)
	assert.NoError(t, err)

	clientset, err := clientset.NewClientset(result.Cluster)
	if err != nil {
		t.Fatalf("Error creating clientset: %v", err)
	}
    

	output, err := eksclusterstatus.EKSClusterStatus(clientset, clusterName, region)
	if err != nil {
		t.Fatalf("Expected EKSComplete(%q, %q) to have a non nil error value; Got %v", clusterName, region, err)
	}

	fmt.Println("########################################################")
	fmt.Println(" 		    CLUSTER STATUS 			    ")
	fmt.Println("########################################################")
	fmt.Println(output)

	fmt.Println("########################################################")
	fmt.Println(" 		    PODS STATUS 			    ")
	fmt.Println("########################################################")
	// Function call to check if external-dns pods are Running
	name := "external-dns"
	namespace := "kube-external-dns"
	labelSelector := "app.kubernetes.io/name=external-dns"
	output, err = podstatus.CheckPodsStatus(clientset, clusterName, region, name, labelSelector, namespace)
	fmt.Println(output)
	if err != nil {
		t.Errorf("Expected EKSComplete(%q, %q) to have a non nil error value; Got %v", clusterName, region, err)
	}
	time.Sleep(10 * time.Second)

	// // Function call to check if ingress pods are Running
	name = "ingress"
	namespace = "kube-ingress"
	labelSelector = "app.kubernetes.io/name=ingress-nginx"
	output, err = podstatus.CheckPodsStatus(clientset, clusterName, region, name, labelSelector, namespace)
	fmt.Println(output)
	if err != nil {
		t.Errorf("Expected EKSComplete(%q, %q) to have a non nil error value; Got %v", clusterName, region, err)
	}
	time.Sleep(10 * time.Second)


	// Function call to check if ebs-csi pod is Running
	name = "ebs-csi"
	namespace = "kube-system"
	labelSelector = "app.kubernetes.io/name=aws-ebs-csi-driver"
	output, err = podstatus.CheckPodsStatus(clientset, clusterName, region, name, labelSelector, namespace)
	fmt.Println(output)
	if err != nil {
		t.Errorf("Expected EKSComplete(%q, %q) to have a non nil error value; Got %v", clusterName, region, err)
	}
	time.Sleep(10 * time.Second)
}

// Function to verify if the resources contain required tags
func TestTagging(t *testing.T) {


	region, regionExists := os.LookupEnv("region")
	if !regionExists {
		t.Skip("region environment variable  not set")
	}

	account_id, accountIdExists := os.LookupEnv("account_id")
	if !accountIdExists {
		t.Skip("account id  not set")
	}

	clusterName, clusterExists := os.LookupEnv("clusterName")
	if !clusterExists {
		t.Skip("clusterName environment variable  not set")
	}


	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		t.Fatalf("unable to load SDK config, %v", err)
	}

	cluster_arn := "arn:aws:eks:us-west-2:" + account_id + ":cluster/" + clusterName

	// Using the Config value, create the Resource Groups Tagging API client
	svc := resourcegroupstaggingapi.NewFromConfig(cfg)
	params := &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: []string{
			cluster_arn,
			//redis_arn, (example to show how arn of resources can be listed)
		},
	}

	resp, err := svc.GetResources(context.TODO(), params)
	if err != nil {
		t.Fatalf("failed to list resources, %v", err)
	}
	for _, res := range resp.ResourceTagMappingList {
		tags_map := map[string]int{"Application": 0, "Customer": 0, "Department": 0, "Environment": 0, "Level": 0, "ManagedBy": 0, "Name": 0}
		fmt.Println("\n####################")
		fmt.Println(*res.ResourceARN)
		for _, tag := range res.Tags {
			fmt.Println(*tag.Key, *tag.Value)
			if _, ok := tags_map[*tag.Key]; ok {
				tags_map[*tag.Key] = 1
			}

		}
		for k, v := range tags_map {
			if v == 0 {
				t.Errorf("This resource does not contain the tag %s\n", k)
			}
		}
		if err != nil {
			os.Exit(1)
		}
	}
}

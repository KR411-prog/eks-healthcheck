package clientset

import (
	
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Function to create client to access eks cluster
func NewClientset(cluster *eks.Cluster) (*kubernetes.Clientset, error) {
	log.Printf("%+v", cluster)
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        aws.StringValue(cluster.Endpoint),
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

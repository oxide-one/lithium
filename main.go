package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

// Control plane node information that we actually care about
type nodeNetworkInformation struct {
	NodeName string
	NodeIP   string
}

type kubeClient struct {
	clientset kubernetes.Interface
}

func grabControlPlaneNodes(roleName *string, k8s *kubeClient) []nodeNetworkInformation {
	// Format a string to contain the controlPlaneRoleName variable
	formattedRoleName := fmt.Sprintf("node-role.kubernetes.io/%s=", *roleName)

	// Grab a list of nodes (machines)
	nodeInformationList, err := k8s.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: formattedRoleName})
	if err != nil {
		// Panic on error
		panic(err.Error())
	}

	// If No nodes are returned, panic and crash
	if len(nodeInformationList.Items) == 0 {
		klog.Error("Unable to get any nodes back from the API. Searching with role name: " + *roleName)
		panic(err.Error())
	}

	// Create a list of nodes
	nodesNetworkInformation := []nodeNetworkInformation{}

	for _, node := range nodeInformationList.Items {
		var nodeIP string
		var nodeName string
		for _, Address := range node.Status.Addresses {
			if Address.Type == "InternalIP" {
				nodeIP = Address.Address
			}
			if Address.Type == "Hostname" {
				nodeName = Address.Address
			}
		}
		//var nodeInfo = &controlPlaneNodes{nodeName: nodeName, nodeIP: nodeIP}
		nodesNetworkInformation = append(nodesNetworkInformation, nodeNetworkInformation{NodeName: nodeName, NodeIP: nodeIP})

	}
	klog.Info(fmt.Sprintf("Nodes returned: %s", nodesNetworkInformation))
	return nodesNetworkInformation
}

func connectivityCheck(k8s *kubeClient) {
	// Run a basic connectivity check to get the server version
	version, err := k8s.clientset.Discovery().ServerVersion()
	if err != nil {
		// Panic if unable to connect
		klog.Error("Unable to connect to the endpoint")
		panic(err.Error())
	}

	// Log the server version to the endpoint
	klog.Info("Connection successful. Server version: " + version.String())

}

// createClient initalises a kubernetes config.
func createClient(kubeConfig *string, alternateEndpoint *string) *kubeClient {
	// Build a kubernetes configuration from flags
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		klog.Error("Unable to create build a config with the given kubeconfig")
		panic(err.Error())
	}

	// If the alternate endpoint is specified, change the TLS server name to the original
	// and then change the endpoint.
	// This is to prevent a TLS issue
	if *alternateEndpoint != "" {
		endpointURL, err := url.Parse(config.Host)
		if err != nil {
			panic(err.Error())
		}
		config.TLSClientConfig.ServerName = endpointURL.Hostname()
		config.Host = *alternateEndpoint
	}
	client := kubeClient{}
	// create the clientset
	client.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &client
}

func main() {
	// Initalize Klog
	klog.InitFlags(nil)

	// Create a variable called kubeconfig that is a string
	var kubeConfig *string
	// Variable to store optional role name
	var controlPlaneRoleName *string

	controlPlaneRoleName = flag.String("roleName", "control-plane", "(optional) the role name of the control plane nodes")

	// Set an alternate API Endpoint
	alternateEndpoint := flag.String("endpoint", "", "(optional) an alternative endpoint for the kubernetes API")

	// If a home directory is not detected
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// Set the klog flag for verbosity
	flag.Set("v", "3")

	// Parse CLI args
	flag.Parse()

	client := createClient(kubeConfig, alternateEndpoint)

	connectivityCheck(client)

	grabControlPlaneNodes(controlPlaneRoleName, client)
}

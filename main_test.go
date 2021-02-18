package main

import (
	"fmt"
	"math/rand"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createNodes(count int) v1.NodeList {
	// Create a list of Nodes to put into the Nodelist later
	nodes := []v1.Node{}
	// Iterate over the count needed
	for i := 1; i < count; i++ {
		// Set a hostname and IP address
		nodeHostName := "control-plane-node-" + RandStringBytes(8)
		nodeIPAddress := fmt.Sprintf("10.0.0.%d", i+10)
		//
		node := v1.Node{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    "Hostname",
						Address: nodeHostName,
					},
					{
						Type:    "InternalIP",
						Address: nodeIPAddress,
					},
				},
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: nodeHostName,
				Labels: map[string]string{
					"node-role.kubernetes.io/control-plane": "",
				},
			},
		}

		nodes = append(nodes, node)
	}
	nodeList := v1.NodeList{Items: nodes}
	return nodeList
}

func TestGrabControlPlaneNodes(t *testing.T) {

	nodes := createNodes(10)

	client := kubeClient{}
	client.clientset = fake.NewSimpleClientset(&nodes)
	//fmt.Println(fak8s.Discovery().ServerVersion())
	roleName := "control-plane"
	grabControlPlaneNodes(&roleName, &client)
}

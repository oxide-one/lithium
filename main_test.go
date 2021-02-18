package main

import (
	"fmt"
	"math/rand"
	"testing"

	v1 "k8s.io/api/core/v1"
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

		nodeAddresses := []v1.NodeAddress{}
		nodeAddresses = append(nodeAddresses, v1.NodeAddress{Type: "Hostname", Address: nodeHostName})
		nodeAddresses = append(nodeAddresses, v1.NodeAddress{Type: "InternalIP", Address: nodeIPAddress})

		nodes = append(nodes, v1.Node{Status: v1.NodeStatus{Addresses: nodeAddresses}})
	}
	fmt.Println(nodes)
	nodeList := v1.NodeList{Items: nodes}
	return nodeList
}

func TestGrabControlPlaneNodes(t *testing.T) {

	createNodes(10)
	//node := v1.Node{Status: v1.NodeStatus{Addresses: []}}

	//fak8s := fake.NewSimpleClientset(&nodes)
	//fmt.Println(fak8s.Discovery().ServerVersion())
	//grabControlPlaneNodes("master")
}

package main

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
}
type DistributedHashMap struct {
	// Add fields for the distributed hash map implementation
	Nodes map[string]*Node
}
type Node struct {
	Name string
	Data map[string]User
}

func main() {
	fmt.Println("hi starting build distrivuted hasing using golang")

	// hashmap := make(map[string]User)

	// hashmap["user1"] = User{Name: "Sanket", Age: 20}

	// fmt.Println(hashmap["user1"])

	dhm := DistributedHashMap{
		Nodes: make(map[string]*Node),
	}
	dhm.Nodes["node1"] = &Node{
		Name: "node1",
		Data: make(map[string]User),
	}
	fmt.Println("intial", dhm)

}

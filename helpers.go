package main

import (
	"log"
)

// Helper function to find Host object by Name attribute in an array of Host objects
func getHostByName(name string, hosts []Host) Host {
	var host Host

	for h := range hosts {
		if hosts[h].Name == name {
			host = hosts[h]
			break
		}
	}

	return host
}

// Helper function to check if string exists within array of strings
// https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Started as copy paste from https://gobyexample.com/
func check(e error) {
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
}

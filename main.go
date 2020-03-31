package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	cluster := flag.String("cluster", "", "The name of the cluster you want to check")
	region := flag.String("region", "", "The name of the AWS region to use. If not provided, the default region will be used.")
	flag.Parse()
	if *cluster == "" {
		log.Fatalln("Need a cluster to operate on to continue, sorry :(")
	}
	fmt.Println(*cluster, *region)
}

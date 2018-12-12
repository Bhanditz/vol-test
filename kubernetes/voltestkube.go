/*
Copyright 2018 Docker, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"log"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
)

type Config struct {
	KubeConfigFile string
}

func getConfig() (config Config) {
	// Get path to kube config.yaml
	filePtr := flag.String("config", "kube.yml", "path to Kubernetes client config yaml")
	flag.Parse()
	config.KubeConfigFile = *filePtr

	return config
}

func printPVCs(pvcs *v1.PersistentVolumeClaimList) {
	template := "%-32s%-8s%-8s\n"
	fmt.Printf(template, "NAME", "STATUS", "CAPACITY")
	for _, pvc := range pvcs.Items {
		quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]
		fmt.Printf(
			template,
			pvc.Name,
			string(pvc.Status.Phase),
			quant.String())
	}
}

func main() {
	configVars := getConfig()
	var err error

	config, err := clientcmd.BuildConfigFromFlags("", configVars.KubeConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Test coverage:

	// Confirm Kube version

	version, err := c.Discovery().ServerVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version is %s\n", version)
	fmt.Println("foo")

	// Confirm test pod exists

	namespace := "default"
	pod := "voltest-0"
	// foo := Pod.new()
	//	api := c.CoreV1()
	_, err = c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}

	// Confirm test pod is running

	p, err := c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})

	fmt.Printf("Pod %s is %s\n", pod, p.Status.Phase)

	fmt.Println(p.Status.PodIP)

	// Confirm that status page of container is happy

	resp, err := http.Get("http://" + p.Status.PodIP + "/status")
	if err != nil {
		fmt.Printf("Response: %s\n", resp.Body)
	}

	// Clear storage data

	// Initialize storage data

	// Confirm textfile

	// Confirm binfile

	// Reschedule container

	// Confirm textfile on rescheduled container

	// confirm binfile on rescheduled container

	//

	// Get Pod by name

	//	label := ""
	//	field := ""

	//	listOptions := metav1.ListOptions{
	//		LabelSelector: label,
	//		FieldSelector: field,
	//	}

	//pvcs, err := c.PersistentVolumeClaims(namespace).List(listOptions)
	//if err != nil {
	//		log.Fatal(err)
	//	}
	//	printPVCs(pvcs)

	// Print its creation time
	//fmt.Println(foo.  )
}

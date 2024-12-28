/**
1. Obtain the Kubeconfig File:
You typically get the kubeconfig file from your Kubernetes cluster administrator
 or the service hosting your cluster (e.g., Google Kubernetes Engine (GKE),
 Amazon EKS, Azure AKS, etc.).

For Cloud Services (e.g., GKE, EKS, AKS):

GKE: You can use the gcloud CLI to get the kubeconfig by running:

gcloud container clusters get-credentials <cluster-name> --region <region> --project <project-id>
EKS: Use the AWS CLI to fetch the kubeconfig:

aws eks update-kubeconfig --region <region> --name <cluster-name>
AKS: Use the Azure CLI to fetch the kubeconfig:

az aks get-credentials --resource-group <resource-group> --name <cluster-name>

For Minikube: You can start Minikube and then automatically set the kubeconfig
using:

minikube start
Manual Generation: If you're setting up your own cluster or working with a
self-hosted Kubernetes setup, you might need to manually create a kubeconfig
file by configuring it with your cluster's details
(API server URL, certificates, tokens, etc.)
*/

/*
*
go get k8s.io/api@v0.28.2
go get k8s.io/apimachinery@v0.28.2
go get k8s.io/client-go@v0.28.2

go get k8s.io/api/core/v1@v0.28.2
go get k8s.io/apimachinery/pkg/api/resource@v0.28.2

--------------------

go mod tidy: This command automatically adds missing modules and
updates the go.sum file

verify if all dependencies are installed properly
go mod verify

Clear and Reinstall Modules (if errors persist):
If the above steps don't resolve the issue,
try clearing the Go module cache and reinstalling:

go clean -modcache
go mod tidy

If you suspect network issues, ensure you have internet access
and that your proxy settings (if any) allow access to Go modules. Use:
go env -w GOPROXY=https://proxy.golang.org,direct
*/
package main

/**
context:
The context package is used to carry deadlines, cancelation signals, and other
request-scoped values across API boundaries and goroutines.
In this code, it helps manage the lifecycle of the pod watch.
Specifically, it allows the program to listen for shutdown signals and
cancel the pod watch process when needed.

flag:
The flag package provides a way to define and parse command-line flags.
Here, it is used to handle the kubeconfig (path to the Kubernetes config file)
and namespace (namespace to monitor) flags.

fmt:
The fmt package is used for formatted I/O, such as printing to the console.
It is used to output log messages and errors (e.g., fmt.Printf()
and fmt.Errorf()).

os:
The os package provides functions for interacting with the operating system.
Here, it is used to handle system signals and create the necessary channels
(e.g., os.Signal) to listen for shutdown signals.

os/signal:
The os/signal package allows the program to capture OS signals, such as
SIGINT (interrupt) or SIGTERM (terminate), which are sent when you stop the
program or shut down the system.
This is crucial for gracefully shutting down the Kubernetes pod monitoring.

syscall:
The syscall package provides low-level primitives to interact with the
operating system's system calls.
In this case, it is used to specify the types of signals (SIGINT, SIGTERM)
that the program should handle to allow for graceful shutdown.

-----------------

v1 "k8s.io/api/core/v1":
This package contains the core API types for Kubernetes, specifically in the
core/v1 group.
Here, it is used to import the Pod type (v1.Pod), which represents the
Kubernetes pod resource that we are monitoring.

metav1 "k8s.io/apimachinery/pkg/apis/meta/v1":
The metav1 package includes metadata-related APIs in Kubernetes, which
are used across most Kubernetes resources.
In this case, it is used to pass metav1.ListOptions{} when watching pods,
defining options related to the pod listing, such as filters.

watch "k8s.io/apimachinery/pkg/watch":
The watch package allows the program to watch for changes (events) to
Kubernetes resources.
It provides types and mechanisms for receiving events like "Added",
"Modified", and "Deleted", which are used to handle pod events in the program.

kubernetes "k8s.io/client-go/kubernetes":
This is the main Kubernetes client Go package that provides access to various
Kubernetes resources (pods, services, etc.).
In this code, it is used to interact with the Kubernetes API server to create
a clientset (clientset), which allows you to communicate with the Kubernetes
cluster.

clientcmd "k8s.io/client-go/tools/clientcmd":
The clientcmd package provides tools for loading Kubernetes configuration,
typically from a kubeconfig file.
It is used here to read the kubeconfig file (provided by the user via the -kubeconfig flag)
and create a config object that is used to authenticate and communicate
with the Kubernetes API server.
*/
import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/*
*
Flags: This segment defines command-line flags:

kubeconfig: Specifies the path to the kubeconfig file that provides the
necessary information to connect to a Kubernetes cluster.
Default is "C:/Users/ethan/.kube/config".

namespace: Defines the Kubernetes namespace in which to monitor pods.
The default namespace is default.

The flag.Parse() reads and processes the flags from the command line.
*/
func main() {
	// Parse kubeconfig flag
	kubeconfig := flag.String("kubeconfig", "C:/Users/ethan/.kube/config", "Path to the kubeconfig file")
	namespace := flag.String("namespace", "default", "Namespace to monitor pods in")
	flag.Parse()

	// Build config from kubeconfig path
	/**
	Config Creation: The clientcmd.BuildConfigFromFlags() function creates the
	Kubernetes client configuration (config) from the kubeconfig file.
	This config contains connection details like the cluster API endpoint,
	authentication credentials, and more.
	*/
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(fmt.Errorf("error building kubeconfig: %v", err))
	}

	// Create a clientset
	/**
	Clientset Creation: kubernetes.NewForConfig(config) creates a clientset,
	which is an object used to interact with various Kubernetes resources
	(like pods, services, deployments, etc.).
	*/
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("error creating clientset: %v", err))
	}

	// Start monitoring pods
	/**
	Context Setup: It creates a context (ctx) for managing the lifecycle of
	the pod watcher. The cancel function allows you to cancel the context,
	which will stop the pod watch when the program exits or receives a
	shutdown signal.
	defer cancel() ensures that the cancel() function is called when the main
	function finishes, cleaning up resources.
	*/
	fmt.Printf("Starting to monitor pods in namespace: %s\n", *namespace)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	/**
	Graceful Shutdown Handling: A goroutine (go handleShutdown(cancel))
	is launched to listen for OS signals (like SIGINT or SIGTERM),
	which allows the program to gracefully terminate when it receives
	a termination signal.
	*/
	go handleShutdown(cancel)

	/**
	The watchPods() function is called to start watching pod events
	in the specified namespace.
	*/
	watchPods(ctx, clientset, *namespace)
}

/*
*
Watcher Creation: clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{}) creates a watcher that listens for pod events (add, modify, delete) in the specified namespace.
Error Handling: If there’s an error in creating the watcher,
the program panics.
defer watcher.Stop() ensures that the watcher is stopped when the
function returns.
*/
func watchPods(ctx context.Context, clientset *kubernetes.Clientset, namespace string) {
	watcher, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("error creating pod watcher: %v", err))
	}
	defer watcher.Stop()

	/**
	Event Loop: The program enters an infinite loop, listening for events from
	the watcher.ResultChan() channel, which delivers pod events
	(such as addition, modification, deletion).

	Event Handling: When an event is received, the program checks if it's an
	error. If it's not an error, it passes the event to handlePodEvent() to
	handle the event further.

	Context Cancellation: If the context (ctx) is canceled
	(for example, when the program shuts down),
	the program prints a shutdown message and exits.
	*/
	for {
		select {
		case event := <-watcher.ResultChan():
			if event.Type == watch.Error {
				fmt.Println("Error occurred while watching pods")
				return
			}
			handlePodEvent(event)
		case <-ctx.Done():
			fmt.Println("Shutting down pod monitor")
			return
		}
	}
}

/*
*
Type Assertion: The event.Object contains the resource affected by the event.
In this case, it expects a Pod. event.Object.(*v1.Pod) performs a type
assertion to ensure the event is related to a pod. If it’s not,
the program prints an error message.
*/
func handlePodEvent(event watch.Event) {
	pod, ok := event.Object.(*v1.Pod)
	if !ok {
		fmt.Println("Unexpected type received from watcher")
		return
	}

	switch event.Type {
	case watch.Added:
		fmt.Printf("Pod added: %s\n", pod.Name)
	case watch.Modified:
		fmt.Printf("Pod modified: %s (Status: %s)\n", pod.Name, pod.Status.Phase)
	case watch.Deleted:
		fmt.Printf("Pod deleted: %s\n", pod.Name)
	}
}

func handleShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
}

/**
Microservices refer to an architectural style in which a software application is
composed of a collection of loosely coupled, independently deployable services.

Benefits of Microservices:
Flexibility: Different teams can work on different services using the best
technology suited for each service.

Scalability: Only the services that need scaling can be scaled, rather than
the entire application.

Resilience: Since services are independent, a failure in one service can be
isolated, making the system more robust.

Faster Development and Deployment: Teams can work independently on services,
enabling continuous integration and continuous deployment (CI/CD).

Example of Microservices in Practice:
Imagine an e-commerce platform. It could be split into the following microservices:

User Service: Handles user registration, login, and authentication.
Product Service: Manages product details and inventory.
Order Service: Manages customer orders and payments.
Shipping Service: Handles the logistics of shipping orders.
*/

package main

/**
The import statement imports the required packages:
fmt for formatting strings and printing output.
net/http for making HTTP requests (to check the health status of microservices).
time for adding delays (waiting 10 seconds between health checks).
*/
import (
	"fmt"
	"net/http"
	"time"
)

// HealthCheck function simulates checking the health of a microservice
/**
healthCheck function: This function simulates checking the health of a microservice
by sending an HTTP GET request to the provided URL and checking the response.

serviceName: The name of the service being checked (e.g., "Service A").

url: The URL where the health check can be accessed (e.g., http://localhost:8081/health).

http.Get(url) sends a GET request to the given URL.

Error Handling: If the request fails (err != nil), it returns a message indicating the
service is "DOWN" along with the error message.

Response Handling: If the request succeeds, it checks the HTTP status code:
If the status code is 200 OK, it returns the message: <serviceName> is UP.
If the status code is anything other than 200, it returns the status message
indicating the service is "DOWN".
*/
func healthCheck(serviceName string, url string) string {
	// Send a GET request to the service
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s is DOWN: %s", serviceName, err)
	}
	defer resp.Body.Close()

	// Return status
	if resp.StatusCode == http.StatusOK {
		return fmt.Sprintf("%s is UP", serviceName)
	}
	return fmt.Sprintf("%s is DOWN: %s", serviceName, resp.Status)
}

func main() {
	// List of microservices and their URLs to check
	/**
	A map called services is defined, which holds the name of each service
	(e.g., "Service A") and the corresponding health check URL
	(e.g., http://localhost:8081/health).
	*/
	services := map[string]string{
		"Service A": "http://localhost:8081/health",
		"Service B": "http://localhost:8082/health",
		"Service C": "http://localhost:8083/health",
	}

	// Simulate checking the health of each service every 10 seconds
	/**
	The program enters an infinite loop (for { ... }), where it continually checks
	the health of each service in the services map.

	for name, url := range services iterates over each service in the map:

	name is the name of the service (e.g., "Service A").

	url is the URL of the health check endpoint (e.g., http://localhost:8081/health).

	status := healthCheck(name, url) calls the healthCheck function to check the
	health of the service.

	fmt.Println(status) prints the health status of each service
	(whether it is "UP" or "DOWN").

	time.Sleep(10 * time.Second) pauses the program for 10 seconds before
	checking the services again.
	*/
	for {
		for name, url := range services {
			status := healthCheck(name, url)
			fmt.Println(status)
		}
		fmt.Println("Waiting for next check...")
		time.Sleep(10 * time.Second)
	}
}

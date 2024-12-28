package main

/**
The net/http package in Go provides HTTP client and server implementations,
allowing you to work with HTTP requests and responses.
*/
import (
	"fmt"
	"net/http"
	"time"
)

// Function to check the HTTP status of a URL
func checkStatus(url string) {
	// Set a timeout for the HTTP request
	client := http.Client{
		Timeout: 10 * time.Second, // 10 seconds timeout
	}

	// Send the HTTP GET request
	resp, err := client.Get(url)
	if err != nil {
		// If there's an error, print the error message
		fmt.Printf("Error checking URL %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Print the status code for the URL
	fmt.Printf("URL: %s, Status Code: %d\n", url, resp.StatusCode)
}

func main() {
	// List of URLs to check
	urls := []string{
		"https://www.google.com",
		"https://www.pixabay.com",
		"https://www.github.com",
	}

	// Check the status of each URL
	for _, url := range urls {
		checkStatus(url)
	}
}

/**
1xx: Informational Responses
100: Continue – The server has received the request headers, and the client should proceed to send the request body.
101: Switching Protocols – The server is switching protocols, as requested by the client.
102: Processing – The server is processing the request, but no response is yet available.

2xx: Success
200: OK – The request was successful, and the server returned the requested content.
201: Created – The request was successful, and a new resource was created.
202: Accepted – The request has been accepted for processing, but the processing is not complete.
203: Non-Authoritative Information – The returned metadata is from a cached copy or a different source than the original server.
204: No Content – The server successfully processed the request, but there is no content to return.
205: Reset Content – The server processed the request and instructed the client to reset the document view.
206: Partial Content – The server is sending only part of the resource (e.g., in response to a range request).

3xx: Redirection
300: Multiple Choices – There are multiple options for the resource that the client can follow.
301: Moved Permanently – The resource has been permanently moved to a new URL.
302: Found (Previously "Moved Temporarily") – The resource is temporarily located at a different URL.
303: See Other – The response to the request can be found under a different URL using the GET method.
304: Not Modified – The resource has not been modified since the last request, so the client can use its cached copy.
305: Use Proxy – The requested resource must be accessed through a proxy.
307: Temporary Redirect – Similar to 302, but the client must use the same method (e.g., POST) for the subsequent request.
308: Permanent Redirect – Similar to 301, but the client must use the same method for the subsequent request.

4xx: Client Errors
400: Bad Request – The request could not be understood or was malformed.
401: Unauthorized – The request requires authentication, and the client has not provided valid credentials.
402: Payment Required – This code is reserved for future use, but it is intended for situations where payment is needed.
403: Forbidden – The server understands the request but refuses to authorize it (client lacks permission).
404: Not Found – The requested resource could not be found on the server.
405: Method Not Allowed – The method specified in the request (e.g., GET, POST) is not allowed for the resource.
406: Not Acceptable – The resource is not able to generate content acceptable according to the Accept headers sent in the request.
407: Proxy Authentication Required – The client must first authenticate with a proxy before accessing the resource.
408: Request Timeout – The server timed out waiting for the client’s request.
409: Conflict – The request could not be processed because of a conflict in the current state of the resource.
410: Gone – The resource is no longer available and will not be available again.
411: Length Required – The server requires the Content-Length header for the request.
412: Precondition Failed – One of the conditions specified by the client in the request header fields failed.
413: Payload Too Large – The request is larger than the server is willing or able to process.
414: URI Too Long – The URI provided is too long for the server to process.
415: Unsupported Media Type – The server does not support the media type of the request.
416: Range Not Satisfiable – The server cannot supply the requested range of a resource.
417: Expectation Failed – The server cannot meet the requirements of the Expect header field.
418: I'm a teapot – An April Fools' joke from the HTTP 1.1 spec. The server is a teapot and cannot brew coffee.
421: Misdirected Request – The request was directed to a server that is not able to produce a response.
422: Unprocessable Entity – The server understands the request, but it cannot process the instructions due to semantic errors.
423: Locked – The resource that is being accessed is locked.
424: Failed Dependency – The request failed due to a failure of a previous request.
425: Too Early – The server is unwilling to risk processing a request that may be replayed.
426: Upgrade Required – The client must upgrade to a different protocol (e.g., HTTP/2).
427: Unassigned – Not yet assigned.
428: Precondition Required – The server requires that the request be conditional (e.g., using If-Match).
429: Too Many Requests – The user has sent too many requests in a given amount of time.
431: Request Header Fields Too Large – The server is unwilling to process the request because the header fields are too large.
451: Unavailable For Legal Reasons – The resource is unavailable for legal reasons (e.g., censorship).

5xx: Server Errors
500: Internal Server Error – The server encountered an unexpected condition that prevented it from fulfilling the request.
501: Not Implemented – The server does not support the functionality required to fulfill the request.
502: Bad Gateway – The server, while acting as a gateway, received an invalid response from the upstream server.
503: Service Unavailable – The server is currently unable to handle the request due to temporary overloading or maintenance.
504: Gateway Timeout – The server, while acting as a gateway, did not receive a timely response from the upstream server.
505: HTTP Version Not Supported – The server does not support the HTTP protocol version used in the request.
506: Variant Also Negotiates – The server has an internal configuration error that prevents it from fulfilling the request.
507: Insufficient Storage – The server is unable to store the representation needed to complete the request.
508: Loop Detected – The server detected an infinite loop while processing a request.
510: Not Extended – The server requires further extensions to fulfill the request.
511: Network Authentication Required – The client needs to authenticate to gain network access.

Summary:
2xx status codes indicate successful requests.
3xx status codes indicate that the client needs to take additional action, such as following a redirect.
4xx status codes indicate client-side errors, such as malformed requests or lack of permissions.
5xx status codes indicate server-side errors, such as internal failures or resource unavailability.
*/

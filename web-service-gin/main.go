/**
Steps
Open a command prompt and change to your home directory.

Using the command prompt, create a directory for your code called web-service-gin.

 mkdir web-service-gin
 cd web-service-gin

Create a module in which you can manage dependencies.
Run the go mod init command, giving it the path of the module your code will be in.

 go mod init example/web-service-gin

Code

use go get to add the github.com/gin-gonic/gin module as a dependency for your module.
Use a dot argument to mean “get dependencies for code in the current directory.

go get .

From the command line in the directory containing main.go, run the code. Use a dot
argument to mean “run code in the current directory.”

go run .

From a new command line window, use curl to make a request to your running web service.

curl http://localhost:8080/albums

For retrieving through id
curl http://localhost:8080/albums/2

New command window for POST

Invoke-WebRequest -Uri http://localhost:8080/albums `
    -Headers @{ "Content-Type" = "application/json" } `
    -Method POST `
    -Body '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'

*/

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/** store album data in memory.

Struct tags such as json:"artist" specify what a field’s name should be when the struct’s
contents are serialized into JSON. Without them, the JSON would use the struct’s capitalized
field names – a style not as common in JSON.*/

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

/*
*Initialize a Gin router using Default.

Use the GET function to associate the GET HTTP method and /albums path with a handler function.

Note that you’re passing the name of the getAlbums function. This is different from passing
the result of the function, which you would do by passing getAlbums() (note the parenthesis).

*/

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

/**gin.Context is the most important part of Gin. It carries request details, validates and
serializes JSON, and more. (Despite the similar name, this is different from Go’s built-in
context package.)

Call Context.IndentedJSON to serialize the struct into JSON and add it to the response.

The function’s first argument is the HTTP status code you want to send to the client.
Here, you’re passing the StatusOK constant from the net/http package to indicate 200 OK.
*/

//GET endpoint

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

//POST enpoint

/**
Use Context.BindJSON to bind the request body to newAlbum.
Append the album struct initialized from the JSON to the albums slice.
Add a 201 status code to the response, along with JSON representing the album you added.
*/

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

/**This getAlbumByID function will extract the ID in the request path, then locate
an album that matches.

Use Context.Param to retrieve the id path parameter from the URL. When you map this handler
to a path, you’ll include a placeholder for the parameter in the path.
Loop over the album structs in the slice, looking for one whose ID field value matches the
id parameter value. If it’s found, you serialize that album struct to JSON and return it as a
response with a 200 OK HTTP code.

As mentioned above, a real-world service would likely use a database query to perform this lookup.

Return an HTTP 404 error with http.StatusNotFound if the album isn’t found.*/

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

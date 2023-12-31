package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.POST("/albums", postAlbums)
	router.DELETE("/delete/:id", deleteAlbum)
	router.PUT("/albums/:id", updateAlbum)

	router.Run("localhost:8080")
}

// album represents data about a record
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// slice of album to data initialize the in-memory database everytime it is started
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// respoonds with a list of all the albums as JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// find an album based on the ID and return it
func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	//loop over the list of albums looking for a matching ID
	for _, a := range albums {

		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found."})
}

// postAlbums adds an album from the JSON received in response body
func postAlbums(c *gin.Context) {
	var newAlbum album
	//call BindJSON to bind the request from the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	// add the new album to the slice
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// deleteAlbum deletes an album based on ID
func deleteAlbum(c *gin.Context) {
	id := c.Param("id")

	//similar to getAlbumById, look for a matching ID and it's index
	index := -1
	for i, a := range albums {
		if a.ID == id {
			index = i
			break
		}
	}

	if index != -1 {
		//delete the album from the data store
		albums = append(albums[:index], albums[index+1:]...)
		//send back a message signifying success
		c.JSON(http.StatusOK, gin.H{"message": "Album deleted succesfully."})
		return
	}
	//send back this message if album not found
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found."})
}

// updateAlbum takes the payload and updates an existing album based on it
func updateAlbum(c *gin.Context) {
	id := c.Param("id")

	//create a struct to represent the updated album data
	var updatedAlbum struct {
		Title  string  `json:"title"`
		Artist string  `json:"artist"`
		Price  float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//look for if the album already exists
	index := -1
	for i, a := range albums {
		if a.ID == id {
			index = i
			break
		}
	}

	if index != -1 {
		albums[index].Title = updatedAlbum.Title
		albums[index].Artist = updatedAlbum.Artist
		albums[index].Price = updatedAlbum.Price

		c.JSON(http.StatusOK, gin.H{"message": "Album updated succesfully."})
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found."})
}

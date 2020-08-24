//Controllers/User.go
package Controllers

import (
	"fileshare/Models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//GetUsers ... Get all users
func GetUsers(c *gin.Context) {
	var user []Models.User
	err := Models.GetAllUsers(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

//CreateUser ... Create User
func CreateUser(c *gin.Context) {
	var user Models.User
	c.BindJSON(&user)
	err := Models.CreateUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate username and/or email!"})
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		c.JSON(http.StatusOK, user)
	}
}

//GetUserByID ... Get the user by id
func GetUserByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var user Models.User
	err := Models.GetUserByID(&user, id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}

	} else {
		c.JSON(http.StatusOK, user)
	}
}

//GetUserByUsername ... Get the user by username
func GetUserByUsername(c *gin.Context) {
	username := c.Params.ByName("username")
	var user Models.User
	err := Models.GetUserByUsername(&user, username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}

	} else {
		c.JSON(http.StatusOK, user)
	}
}

// //UpdateUser ... Update the user information
// func UpdateUser(c *gin.Context) {
// 	var user Models.User
// 	id := c.Params.ByName("id")
// 	err := Models.GetUserByID(&user, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, user)
// 	}
// 	c.BindJSON(&user)
// 	err = Models.UpdateUser(&user, id)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 	} else {
// 		c.JSON(http.StatusOK, user)
// 	}
// }

//DeleteUser ... Delete the user
func DeleteUser(c *gin.Context) {
	var user Models.User
	id := c.Params.ByName("id")
	err := Models.DeleteUser(&user, id)
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "1451") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User has dependencies!"})
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id, "deleted": true})
	}
}

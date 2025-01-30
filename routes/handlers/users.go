package handlers

import (
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
    var user structs.User
    if !utils.BindJSON(c, &user) {
        return
    }
    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
	fmt.Fprintln(os.Stdout, "USER ADD!")
}

func UpdateUser(c *gin.Context) {
	_, err, user := utils.CheckUser(c)
    if !err {
        return
    }

	if !utils.BindJSON(c, &user) {
        return
    }
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	fmt.Fprintln(os.Stdout, "USER UPDATED!")
}

func ReadUser(c *gin.Context) {
	_, ok, user := utils.CheckUser(c)
    if !ok {
        return
    }

	c.JSON(http.StatusOK, user)
	fmt.Fprintln(os.Stdout, "USER READ!")
}
package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindJSON (c *gin.Context, input interface{}) bool {
	if err:=c.ShouldBindJSON(input);err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return false
	}
	return true
}

func CheckUser (c *gin.Context) (int, bool) {
	userID,err :=strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "Invalid user ID"})
		return 0, false
	}
	return userID, true
}
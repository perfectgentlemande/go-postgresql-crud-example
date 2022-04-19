package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (ct *Controller) GetUser(c *gin.Context, params GetUserParams) {
	fmt.Println("Hello GetUser")
}
func (ct *Controller) PostUser(c *gin.Context) {
	fmt.Println("Hello PostUser")
}
func (ct *Controller) DeleteUserId(c *gin.Context, id string) {
	fmt.Println("Hello DeleteUserId")
}
func (ct *Controller) GetUserId(c *gin.Context, id string) {
	fmt.Println("Hello GetUserId")
}
func (ct *Controller) PutUserId(c *gin.Context, id string) {
	fmt.Println("Hello PutUserId")
}

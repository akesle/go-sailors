package service

import (
  "github.com/akesle/sailors/controllers"
  "github.com/gin-gonic/gin"
)

func Run() error {
  r := gin.Default()
  r.POST("/sailors", controllers.AddSailor)
  r.GET("/sailors", controllers.FindSailor)
  r.DELETE("/sailors", controllers.RemoveSailor)
  r.PUT("/sailors", controllers.ModifySailor)
  return r.Run()
}

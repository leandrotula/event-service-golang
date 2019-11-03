package main

import (
  "./service"
  "github.com/gin-gonic/gin"
  "log"
)

func SetupRouter() *gin.Engine {
  router := gin.Default()

  v1 := router.Group("api/v1") 
  {
    v1.GET("/events", service.GetEvents)
    v1.GET("/events/:id", service.GetEvent)
    v1.POST("/events", service.PostEvent)
    v1.PUT("/events/:id", service.UpdateEvent)
    v1.DELETE("/events/:id", service.DeleteEvent)
  }

  return router
}

func main() {
  router := SetupRouter()
  err := router.Run(":8080")

  if err != nil {
    log.Fatal("FATAL: Could not start the application")
  }
}

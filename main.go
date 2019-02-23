package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/zodahu/bitmark-map/server"
)

func setupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(nodesMiddleware())
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	r.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache")
	})

	r.POST("/register", server.Register)
	r.GET("/nodes", server.GetNodes)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}

// nodesMiddleware adds the nodes to the context
func nodesMiddleware() gin.HandlerFunc {
	ns := server.GetNodesInstance()
	return func(c *gin.Context) {
		c.Set("nodes", ns)
		c.Next()
	}
}

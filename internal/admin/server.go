package admin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jin06/mercury/internal/admin/handlers"
)

type adminServer struct {
	ctx context.Context
}

func (s *adminServer) start() (err error) {
	err = s.startAPI(s.ctx)
	if err != nil {
		return
	}
	return
}

func (s *adminServer) stop() error {
	return nil
}

func (s *adminServer) startAPI(ctx context.Context) error {
	r := gin.Default()
	r.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	{
		user := handlers.User{}
		userGroup := r.Group("/admin/api/user")
		userGroup.GET("/login", user.Login)
		userGroup.GET("/info", user.Info)
	}

	return r.Run(":8080") // Start the server on port 8080
}

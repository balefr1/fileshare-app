//Routes/Routes.go
package Routes

import (
	"fileshare/Controllers"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

//SetupRouter ... Configure routes
func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(CORSMiddleware(),
		gin.LoggerWithWriter(gin.DefaultWriter, "/health"),
		gin.Recovery(),
	)

	// --- STATIC CONTENT (REACT JS APP)(THESE SHOULD BE SERVED BY STANDALONE WEB SERVER NGINX APACHE (or even S3))
	r.Use(static.Serve("/", static.LocalFile("./Client/build", true)))
	r.StaticFile("/favicon.ico", "Client/public/favicon.ico")
	r.LoadHTMLGlob("Client/build/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/user/*action", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// --- END STATIC

	grp1 := r.Group("/user-api")
	{
		grp1.GET("users", Controllers.GetUsers)
		grp1.GET("username/:username", Controllers.GetUserByUsername)
		grp1.POST("user", Controllers.CreateUser)
		grp1.GET("user/:id", Controllers.GetUserByID)
		//grp1.PUT("user/:id", Controllers.UpdateUser)
		grp1.DELETE("user/:id", Controllers.DeleteUser)
	}
	grp2 := r.Group("/user-files")
	{
		grp2.GET("attachments/:username", Controllers.GetAttachments)
		grp2.GET("attachment/:id", Controllers.GetAttachmentByID)
		grp2.GET("attachment/:id/get", Controllers.DwAttachmentByID)
		grp2.POST("attachment/:username", Controllers.CreateAttachment)
		grp2.POST("attachment/:username/:place", Controllers.CreateAttachment)
		grp2.DELETE("/attachment/:id", Controllers.DeleteAttachment)
	}

	//r.LoadHTMLGlob("Templates/*")

	//health-check path
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	//r.POST("/upload", Controllers.Upload)
	//r.StaticFS("/static", http.Dir("Pages/static"))

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

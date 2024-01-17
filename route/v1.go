package route

import (
	"crypto/ecdsa"
	"os"

	"github.com/IceWhaleTech/CasaOS-Common/middleware"
	"github.com/IceWhaleTech/CasaOS-Common/utils/jwt"
	v1 "github.com/IceWhaleTech/CasaOS-UserService/route/v1"
	"github.com/IceWhaleTech/CasaOS-UserService/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	// r.Use(middleware.WriteLog())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// check if environment variable is set
	if ginMode, success := os.LookupEnv("GIN_MODE"); success {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.POST("/v1/users/register", v1.PostUserRegister)    // register
	r.POST("/v1/users/login", v1.PostUserLogin)          // login
	r.GET("/v1/users/name", v1.GetUserAllUsername)       // all/name
	r.POST("/v1/users/refresh", v1.PostUserRefreshToken) // refresh token
	// No short-term modifications
	r.GET("/v1/users/image", v1.GetUserImage) // image

	r.GET("/v1/users/status", v1.GetUserStatus) // init/check

	v1Group := r.Group("/v1")

	v1Group.Use(jwt.JWT(
		func() (*ecdsa.PublicKey, error) {
			_, publicKey := service.MyService.User().GetKeyPair()
			return publicKey, nil
		},
	))
	{
		v1UsersGroup := v1Group.Group("/users")
		v1UsersGroup.Use()
		{
			v1UsersGroup.GET("/current", v1.GetUserInfo)              // 获取当前用户信息
			v1UsersGroup.PUT("/current", v1.PutUserInfo)              // 修改当前用户信息
			v1UsersGroup.PUT("/current/password", v1.PutUserPassword) // 修改当前用户密码

			v1UsersGroup.GET("/current/custom/:key", v1.GetUserCustomConf)       // 获取当前用户自定义配置,比如app的安装列表(有序)
			v1UsersGroup.POST("/current/custom/:key", v1.PostUserCustomConf)     // 修改当前用户自定义配置，比如切换背景墙
			v1UsersGroup.DELETE("/current/custom/:key", v1.DeleteUserCustomConf) // 删除当前用户自定义配置

			v1UsersGroup.POST("/current/image/:key", v1.PostUserUploadImage) // 上传壁纸
			v1UsersGroup.PUT("/current/image/:key", v1.PutUserImage)
			// v1UserGroup.POST("/file/image/:key", v1.PostUserFileImage)
			v1UsersGroup.DELETE("/current/image", v1.DeleteUserImage) // 删除用户头像

			v1UsersGroup.PUT("/avatar", v1.PutUserAvatar) // 修改用户头像,没看到使用
			v1UsersGroup.GET("/avatar", v1.GetUserAvatar) // 获取用户头像

			v1UsersGroup.DELETE("/:id", v1.DeleteUser)               // 删除用户
			v1UsersGroup.GET("/:username", v1.GetUserInfoByUsername) // 获取用户信息
			v1UsersGroup.DELETE("", v1.DeleteUserAll)                // 删除所有用户
		}
	}

	return r
}

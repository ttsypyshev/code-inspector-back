package backend

import (
	"context"
	"net/http"
	"rip/lib/auth"
	"rip/lib/database"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (app *App) SetupRoutes(r *gin.Engine) {
	// Услуги (домен: `/info`)
	r.GET("/info", app.ChekcUser(), app.GetServiceList)                                               // Получить список услуг с фильтрацией, включая ID заявки-черновика пользователя и количество услуг в этой заявке.
	r.GET("/info/:id", app.GetServiceByID)                                                            // Получить данные конкретной услуги по ее ID.
	r.POST("/info", app.AuthMiddleware(), RoleMiddleware(database.Admin), app.CreateService)          // Добавить новую услугу (без изображения).
	r.PUT("/info/:id", app.AuthMiddleware(), RoleMiddleware(database.Admin), app.UpdateService)       // Изменить данные услуги по ее ID.
	r.POST("/info/:id", app.AuthMiddleware(), RoleMiddleware(database.Admin), app.UpdateServiceImage) // Добавить или заменить изображение для услуги с указанным ID. Если изображение уже существует, оно заменяется.
	r.DELETE("/info/:id", app.AuthMiddleware(), RoleMiddleware(database.Admin), app.DeleteService)    // Удалить услугу вместе с изображением.
	r.POST("/info/add-service", app.AuthMiddleware(), app.AddServiceToDraft)                          // Добавить услугу в заявку-черновик. Если черновик отсутствует, создается новый с указанными значениями.

	// Заявки (домен: `/project`)
	r.GET("/project", app.AuthMiddleware(), app.GetProjectList)                                               // Получить список заявок с фильтрацией по диапазону даты формирования и статусу (исключая удаленные и черновики, поля модератора и создателя отображаются через логины).
	r.GET("/project/:id", app.AuthMiddleware(), app.GetProjectByID)                                           // Получить данные конкретной заявки по ее ID, включая список услуг и изображения.
	r.PUT("/project/:id", app.AuthMiddleware(), app.UpdateProject)                                            // Изменить поля заявки по теме.
	r.PUT("/project/:id/submit", app.AuthMiddleware(), app.SubmitProject)                                     // Сформировать заявку создателем с установкой даты формирования. Происходит проверка обязательных полей.
	r.PUT("/project/:id/complete", app.AuthMiddleware(), RoleMiddleware(database.Admin), app.CompleteProject) // Завершить или отклонить заявку модератором, указываются модератор и дата завершения. При завершении производится расчет дополнительных полей.
	r.DELETE("/project/:id", app.AuthMiddleware(), app.DeleteProject)                                         // Удалить заявку (удаляется только при отсутствии даты формирования).

	// ММ (домен: `/file`)
	r.DELETE("/file/delete", app.AuthMiddleware(), app.DeleteFileFromProject) // Удалить файл из заявки (без первичного ключа файла).
	r.PUT("/file/update", app.AuthMiddleware(), app.UpdateFileInProject)      // Изменить количество, порядок или значение файла в заявке (без первичного ключа файла).

	// Пользователи (домен: `/user`)
	r.POST("/user/register", app.RegisterUser)                         // Регистрация нового пользователя.
	r.PUT("/user/update", app.AuthMiddleware(), app.UpdateUserProfile) // Изменить данные пользователя (личный кабинет).
	r.POST("/user/login", app.UserLogin)                               // Аутентификация пользователя.
	r.POST("/user/logout", app.AuthMiddleware(), app.UserLogout)       // Деавторизация пользователя.
}

func (app *App) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := auth.ValidateJWT(token, app.secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Проверка сессии в Redis
		ctx := context.Background()
		userID, role := claims.UserID, claims.Role
		isValid, err := CheckSessionExists(ctx, app.redisClient, userID)
		if err != nil || !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}

func (app *App) ChekcUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			c.Set("userID", uuid.Nil.String())
			c.Next()
			return
		}

		claims, err := auth.ValidateJWT(token, app.secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Проверка сессии в Redis
		ctx := context.Background()
		userID, role := claims.UserID, claims.Role
		isValid, err := CheckSessionExists(ctx, app.redisClient, userID)
		if err != nil || !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}

func RoleMiddleware(role database.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != role.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

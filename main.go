package main

import (
	"context"
	"log"
	"os"
	"time"

	"member_API/controllers"
	_ "member_API/docs" // 導入 swagger 文檔
	"member_API/graphql"
	"member_API/models"
	"member_API/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title Member API
// @version 1.0
// @description 這是一個使用 Go、Gin 框架和 PostgreSQL 構建的 RESTful 和 GraphQL API 服務，提供會員管理功能和 JWT 認證
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT 認證，格式：Bearer {token}

var db *gorm.DB

func initPostgreSQL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/member_api?sslmode=disable"
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	if err := gormDB.WithContext(ctx).AutoMigrate(&models.Member{}); err != nil {
		return err
	}

	db = gormDB
	controllers.SetupUserController(db)

	// 初始化 GraphQL
	if err := graphql.SetupGraphQL(db); err != nil {
		log.Printf("Warning: GraphQL setup failed: %v\n", err)
	} else {
		log.Println("GraphQL schema initialized!")
	}

	log.Println("Connected to PostgreSQL!")
	return nil
}

// HealthCheck 健康檢查端點
// @Summary 健康檢查
// @Description 檢查服務器狀態和數據庫連接狀態
// @Tags 系統
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "服務正常"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	status := "OK"
	dbStatus := "Connected"
	if db == nil {
		dbStatus = "Disconnected"
	} else if sqlDB, err := db.DB(); err != nil {
		dbStatus = "Error: " + err.Error()
	} else if err := sqlDB.PingContext(context.Background()); err != nil {
		dbStatus = "Error: " + err.Error()
	}
	c.JSON(200, gin.H{
		"status":          status,
		"postgres_status": dbStatus,
	})
}

func main() {
	// 初始化 PostgreSQL 連接
	err := initPostgreSQL()
	if err != nil {
		log.Printf("Warning: PostgreSQL connection failed: %v\n", err)
		log.Println("Starting server without PostgreSQL connection...")
	} else {
		defer func() {
			if sqlDB, err := db.DB(); err == nil {
				if err := sqlDB.Close(); err != nil {
					log.Printf("Error closing PostgreSQL connection: %v\n", err)
				}
			} else {
				log.Printf("Error retrieving SQL DB handle: %v\n", err)
			}
		}()
	}

	// 創建 Gin 路由器
	r := gin.Default()

	// 設置路由
	routes.SetupRouter(r)

	// Swagger 文檔路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 添加一個簡單的健康檢查端點
	r.GET("/health", HealthCheck)

	// 啟動服務器
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

package router

import (
	"github.com/Rakhulsr/foodcourt/internal/delivery/http"
	adminHandler "github.com/Rakhulsr/foodcourt/internal/delivery/http/admin"
	"github.com/Rakhulsr/foodcourt/internal/delivery/http/client"
	"github.com/Rakhulsr/foodcourt/internal/middleware"
	"github.com/Rakhulsr/foodcourt/internal/repository"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/Rakhulsr/foodcourt/pkg/engine"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())
	engine.SetupViewEngine(r)

	r.Use(middleware.FlashMessage())
	r.Use(middleware.CSRFProtection())

	boothRepo := repository.NewBoothRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	logRepo := repository.NewWhatsAppLogRepository(db)
	waUC := usecase.NewWhattsAppUsecase()

	boothUC := usecase.NewBoothUseCase(boothRepo)
	menuUC := usecase.NewMenuUseCase(menuRepo, boothRepo)
	authUC := usecase.NewAuthUseCase(adminRepo)
	paymentUC := usecase.NewPaymentService()

	orderUC := usecase.NewOrderUsecase(orderRepo, menuRepo, paymentUC, *waUC)
	logUC := usecase.NewLogUseCase(logRepo)

	adminMenuHandler := adminHandler.NewMenuHandler(menuUC, boothUC)
	adminBoothHandler := adminHandler.NewBoothHandler(boothUC)
	adminOrderHandler := adminHandler.NewOrderHandler(orderUC)
	dashboardHandler := adminHandler.NewDashboardHandler(orderRepo)
	adminLogHandler := adminHandler.NewLogHandler(logUC)

	menuHandler := client.NewMenuHandler(menuUC, boothUC)
	cartHandler := client.NewCartHandler(menuUC)
	orderHandler := client.NewOrderHandler(orderUC)

	authHandler := http.NewAuthHandler(authUC)

	r.GET("/", menuHandler.ClientHome)
	r.GET("/home", menuHandler.ClientHome)

	r.GET("/cart", cartHandler.ShowCart)
	r.POST("/cart/add", cartHandler.AddToCart)
	r.POST("/cart/update", cartHandler.UpdateCartItem)
	r.POST("/cart/proceed", cartHandler.ProceedCheckout)
	r.GET("/checkout", cartHandler.ShowCheckoutPage)

	r.GET("/order/success/:code", orderHandler.ShowSuccessPage)

	api := r.Group("/api")
	{
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "running", "api": "FoodCourt v1"})
		})

		api.POST("/orders", orderHandler.Create)

		api.POST("/webhooks/xendit", orderHandler.HandleXenditWebhook)

		api.GET("/menus", menuHandler.ListActive)
		api.GET("/menus/search", menuHandler.Search)

		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.JWTAuth())
		{
			adminRoutes.GET("/dashboard", dashboardHandler.Dashboard)

			adminRoutes.GET("/booths", adminBoothHandler.AdminList)
			adminRoutes.GET("/booths/create", adminBoothHandler.ShowCreateForm)
			adminRoutes.POST("/booths", adminBoothHandler.Create)
			adminRoutes.GET("/booths/edit/:id", adminBoothHandler.ShowEditForm)
			adminRoutes.PUT("/booths/:id", adminBoothHandler.Update)
			adminRoutes.DELETE("/booths/:id", adminBoothHandler.Delete)

			adminRoutes.GET("/menus", adminMenuHandler.ListAll)
			adminRoutes.GET("/menus/create", adminMenuHandler.ShowCreateForm)
			adminRoutes.POST("/menus", adminMenuHandler.Create)
			adminRoutes.GET("/menus/edit/:id", adminMenuHandler.ShowEditForm)
			adminRoutes.PUT("/menus/:id", adminMenuHandler.Update)
			adminRoutes.DELETE("/menus/:id", adminMenuHandler.Delete)

			adminRoutes.GET("/orders", adminOrderHandler.AdminList)
			adminRoutes.PATCH("/orders/:code/status", adminOrderHandler.AdminUpdateStatus)
			adminRoutes.POST("/orders/:code/notify", adminOrderHandler.SendNotification)

			adminRoutes.GET("/logs", adminLogHandler.List)

			adminRoutes.GET("/logs/track", adminLogHandler.TrackAndRedirect)
		}
	}

	auth := r.Group("/auth")
	{
		auth.GET("/login", authHandler.ShowLoginForm)
		auth.POST("/login", authHandler.Login)
		auth.GET("/logout", authHandler.Logout)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	return r
}

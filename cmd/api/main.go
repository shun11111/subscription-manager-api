package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"subscription-manager-api/internal/interface/http/handler"
	authmw "subscription-manager-api/internal/interface/http/middleware"
	"subscription-manager-api/internal/infrastructure/persistence"
	"subscription-manager-api/internal/usecase"
	"subscription-manager-api/internal/infrastructure/config"
)

func main() {
	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリ初期化
	userRepo := persistence.NewUserRepository(db)
	subscriptionRepo := persistence.NewSubscriptionRepository(db)
	planRepo := persistence.NewPlanRepository(db)

	// ユースケース初期化
	authService := usecase.NewAuthService(userRepo, cfg.JWTSecret)
	subscriptionService := usecase.NewSubscriptionService(subscriptionRepo)
	planService := usecase.NewPlanService(planRepo)

	// ハンドラ初期化
	authHandler := handler.NewAuthHandler(authService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	planHandler := handler.NewPlanHandler(planService)

	// Echo初期化
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルーティング
	api := e.Group("/api")
	{
		// 認証（認証不要）
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/login", authHandler.Login)
		}

		// サブスクリプション（認証必要）
		subscriptions := api.Group("/subscriptions")
		subscriptions.Use(authmw.AuthMiddleware(cfg.JWTSecret))
		{
			subscriptions.GET("", subscriptionHandler.ListSubscriptions)
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
		}

		// プラン・価格マスター（認証必要）
		plans := api.Group("/plans")
		plans.Use(authmw.AuthMiddleware(cfg.JWTSecret))
		{
			plans.GET("", planHandler.ListPlans)
			plans.POST("", planHandler.CreatePlan)
			plans.GET("/:id", planHandler.GetPlan)
			plans.PUT("/:id", planHandler.UpdatePlan)
			plans.DELETE("/:id", planHandler.DeletePlan)
		}
	}

	// サーバー起動
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		if err := e.Start(addr); err != nil {
			log.Printf("Server shutdown: %v", err)
		}
	}()

	// グレースフルシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}


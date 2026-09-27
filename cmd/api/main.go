package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"

	"github.com/eduardeddy/content-ms/internal/article_categories"
	"github.com/eduardeddy/content-ms/internal/articles"
	"github.com/eduardeddy/content-ms/internal/config"
	"github.com/eduardeddy/content-ms/internal/publisher"
	"github.com/eduardeddy/content-ms/internal/router"
	"github.com/eduardeddy/content-ms/internal/stories"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.GetDSN())
	if err != nil {
		log.Fatalf("❌ Error configurando conexión a DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ No se pudo conectar a PostgreSQL: %v", err)
	}
	log.Println("✅ Conexión exitosa a PostgreSQL mediante pgx/v5")

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Printf("⚠️ No se pudo conectar a NATS en %s: %v", cfg.NatsURL, err)
	} else {
		defer nc.Close()
		log.Printf("✅ Conexión exitosa a NATS (%s)", cfg.NatsURL)
	}

	publisherRepo := publisher.NewRepository(db)
	publisherService := publisher.NewService(publisherRepo)
	publisherHandler := publisher.NewHandler(publisherService)

	storiesRepo := stories.NewRepository(db)
	storiesService := stories.NewService(storiesRepo, publisherRepo)
	storiesHandler := stories.NewHandler(storiesService)

	articlesRepo := articles.NewRepository(db)
	articlesService := articles.NewService(articlesRepo, publisherRepo)
	articlesHandler := articles.NewHandler(articlesService)

	articleCategoriesRepo := article_categories.NewRepository(db)
	articleCategoriesService := article_categories.NewService(articleCategoriesRepo)
	articleCategoriesHandler := article_categories.NewHandler(articleCategoriesService)

	if nc != nil {
		publisherSubscriber := publisher.NewSubscriber(publisherService)
		publisherSubscriber.RegisterSubscribers(nc)

		storiesSubscriber := stories.NewSubscriber(storiesService)
		storiesSubscriber.RegisterSubscribers(nc)

		articlesSubscriber := articles.NewSubscriber(articlesService)
		articlesSubscriber.RegisterSubscribers(nc)

		articleCategoriesSubscriber := article_categories.NewSubscriber(articleCategoriesService)
		articleCategoriesSubscriber.RegisterSubscribers(nc)
	}

	r := router.Setup(publisherHandler, storiesHandler, articlesHandler, articleCategoriesHandler)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Servicio content-ms corriendo en el puerto %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error iniciando el servidor HTTP: %v", err)
	}
}

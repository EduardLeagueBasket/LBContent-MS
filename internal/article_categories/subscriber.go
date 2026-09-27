package article_categories

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	service Service
}

func NewSubscriber(service Service) *Subscriber {
	return &Subscriber{service: service}
}

type nestPacket[T any] struct {
	Pattern string `json:"pattern"`
	Data    T      `json:"data"`
	ID      string `json:"id"`
}

func (s *Subscriber) RegisterSubscribers(nc *nats.Conn) {
	subscribe := func(subject string, handler nats.MsgHandler) {
		if _, err := nc.Subscribe(subject, handler); err != nil {
			log.Printf("❌ Error al registrar suscriptor NATS para %s: %v", subject, err)
		}
	}

	subscribe("article-category.create", s.handleCreate)
	subscribe("article-category.list", s.handleList)
}

func (s *Subscriber) handleCreate(m *nats.Msg) {
	log.Printf("[NATS] article-category.create recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[ArticleCategoryDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.Name == "" {
		dto, err = parseFlat[ArticleCategoryDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := s.service.CreateArticleCategory(ctx, dto)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al crear categoría de artículo: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	log.Printf("✅ Categoría de artículo creada: %s", created.Name)
	replyNestOK(m, packet.ID, created)
}

func (s *Subscriber) handleList(m *nats.Msg) {
	log.Printf("[NATS] article-category.list recibido: %s", m.Data)

	packet, _, err := parseNestPacket[ArticleCategoryDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := s.service.ListArticleCategories(ctx)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar categorías de artículo: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	log.Printf("✅ Categorías de artículo listadas: %d", len(items))
	replyNestOK(m, packet.ID, items)
}

func parseNestPacket[T any](data []byte) (nestPacket[T], T, error) {
	var packet nestPacket[T]
	if err := json.Unmarshal(data, &packet); err != nil {
		var zero T
		return packet, zero, err
	}
	return packet, packet.Data, nil
}

func parseFlat[T any](data []byte) (T, error) {
	var dto T
	err := json.Unmarshal(data, &dto)
	return dto, err
}

func replyNestOK(m *nats.Msg, id string, response any) {
	if m.Reply == "" {
		return
	}
	payload, err := json.Marshal(map[string]any{
		"id":         id,
		"response":   response,
		"isDisposed": true,
	})
	if err != nil {
		log.Printf("[NATS ERROR] No se pudo serializar respuesta: %v", err)
		return
	}
	if err := m.Respond(payload); err != nil {
		log.Printf("[NATS ERROR] No se pudo responder: %v", err)
	}
}

func replyNestError(m *nats.Msg, id string, err error) {
	if m.Reply == "" {
		return
	}
	payload, marshalErr := json.Marshal(map[string]any{
		"id":         id,
		"err":        err.Error(),
		"isDisposed": true,
	})
	if marshalErr != nil {
		log.Printf("[NATS ERROR] No se pudo serializar error: %v", marshalErr)
		return
	}
	if respondErr := m.Respond(payload); respondErr != nil {
		log.Printf("[NATS ERROR] No se pudo responder error: %v", respondErr)
	}
}

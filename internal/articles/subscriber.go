package articles

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

	subscribe("articles.create", s.handleCreate)
	subscribe("articles.list", s.handleList)
	subscribe("articles.list-by-entity", s.handleListByEntity)
	subscribe("articles.feed-by-country", s.handleFeedByCountry)
	subscribe("articles.get", s.handleGet)
	subscribe("articles.update", s.handleUpdate)
	subscribe("articles.delete", s.handleDelete)
}

func (s *Subscriber) handleCreate(m *nats.Msg) {
	log.Printf("[NATS] articles.create recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[CreateArticleDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.Title == "" {
		dto, err = parseFlat[CreateArticleDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.CreateArticle(ctx, dto); err != nil {
		log.Printf("[NATS ERROR] Fallo al crear artículo: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	log.Printf("✅ Artículo creado: %s", dto.Title)
	replyNestOK(m, packet.ID, map[string]any{
		"message":     "created",
		"title":       dto.Title,
		"entity_type": dto.EntityType,
		"entity_id":   dto.EntityID,
	})
}

func (s *Subscriber) handleList(m *nats.Msg) {
	log.Printf("[NATS] articles.list recibido: %s", m.Data)

	packet, _, err := parseNestPacket[ArticleDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := s.service.ListArticles(ctx)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar artículos: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, items)
}

func (s *Subscriber) handleListByEntity(m *nats.Msg) {
	log.Printf("[NATS] articles.list-by-entity recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[ListByEntityDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.EntityID == "" {
		dto, err = parseFlat[ListByEntityDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := s.service.ListByEntity(ctx, dto)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar artículos por entidad: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, items)
}

func (s *Subscriber) handleFeedByCountry(m *nats.Msg) {
	log.Printf("[NATS] articles.feed-by-country recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[FeedByCountryDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.CountryCode == "" {
		dto, err = parseFlat[FeedByCountryDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := s.service.ListByCountry(ctx, dto.CountryCode)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar artículos por país: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, items)
}

func (s *Subscriber) handleGet(m *nats.Msg) {
	log.Printf("[NATS] articles.get recibido: %s", m.Data)

	packet, dto, err := parseArticleIDPayload(m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	article, err := s.service.GetArticle(ctx, dto.ID)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al obtener artículo: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, article)
}

func (s *Subscriber) handleUpdate(m *nats.Msg) {
	log.Printf("[NATS] articles.update recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[ArticleDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.ID == "" && dto.Title == "" {
		dto, err = parseFlat[ArticleDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.UpdateArticle(ctx, dto.ID, dto); err != nil {
		log.Printf("[NATS ERROR] Fallo al actualizar artículo: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, map[string]any{"message": "updated", "id": dto.ID})
}

func (s *Subscriber) handleDelete(m *nats.Msg) {
	log.Printf("[NATS] articles.delete recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[DeleteArticleDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.ID == "" {
		dto, err = parseFlat[DeleteArticleDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.DeleteArticle(ctx, dto); err != nil {
		log.Printf("[NATS ERROR] Fallo al eliminar artículo: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	log.Printf("✅ Artículo eliminado: %s", dto.ID)
	replyNestOK(m, packet.ID, map[string]any{"message": "deleted", "id": dto.ID})
}

func parseArticleIDPayload(data []byte) (nestPacket[ArticleIDDTO], ArticleIDDTO, error) {
	packet, dto, err := parseNestPacket[ArticleIDDTO](data)
	if err != nil {
		return packet, dto, err
	}
	if dto.ID != "" {
		return packet, dto, nil
	}
	flat, err := parseFlat[ArticleIDDTO](data)
	return packet, flat, err
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

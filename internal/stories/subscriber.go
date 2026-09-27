package stories

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

	subscribe("stories.create", s.handleCreate)
	subscribe("stories.list-by-entity", s.handleListByEntity)
	subscribe("stories.feed-by-country", s.handleFeedByCountry)
	subscribe("stories.delete", s.handleDelete)
}

func (s *Subscriber) handleCreate(m *nats.Msg) {
	log.Printf("[NATS] stories.create recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[CreateStoryDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.EntityID == "" {
		dto, err = parseFlat[CreateStoryDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.CreateStory(ctx, dto); err != nil {
		log.Printf("[NATS ERROR] Fallo al crear historia: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	log.Printf("✅ Historia creada exitosamente para publisher %s/%s", dto.EntityType, dto.EntityID)
	replyNestOK(m, packet.ID, map[string]any{
		"message":     "created",
		"entity_type": dto.EntityType,
		"entity_id":   dto.EntityID,
	})
}

func (s *Subscriber) handleListByEntity(m *nats.Msg) {
	log.Printf("[NATS] stories.list-by-entity recibido: %s", m.Data)

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

	feed, err := s.service.ListByEntity(ctx, dto)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar historias: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, feed)
}

func (s *Subscriber) handleFeedByCountry(m *nats.Msg) {
	log.Printf("[NATS] stories.feed-by-country recibido: %s", m.Data)

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

	feed, err := s.service.GetFeedByCountry(ctx, dto.CountryCode)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar feed por país: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, feed)
}

func (s *Subscriber) handleDelete(m *nats.Msg) {
	log.Printf("[NATS] stories.delete recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[DeleteStoryDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.ID == "" {
		dto, err = parseFlat[DeleteStoryDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.DeleteStory(ctx, dto); err != nil {
		log.Printf("[NATS ERROR] Fallo al eliminar historia: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	log.Printf("✅ Historia eliminada: %s", dto.ID)
	replyNestOK(m, packet.ID, map[string]any{"message": "deleted", "id": dto.ID})
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

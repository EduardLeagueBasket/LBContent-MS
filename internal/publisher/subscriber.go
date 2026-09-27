package publisher

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

	subscribe("publisher.sync", s.handleSync)
	subscribe("publisher.list", s.handleList)
	subscribe("publisher.create", s.handleCreate)
}

func (s *Subscriber) handleSync(m *nats.Msg) {
	log.Printf("[NATS] publisher.sync recibido: %s", m.Data)

	dto, err := parsePublisherPayload(m.Data)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al parsear JSON de publisher: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.UpsertPublisher(ctx, dto); err != nil {
		log.Printf("[NATS ERROR] Fallo al guardar/actualizar publisher: %v", err)
		return
	}

	log.Printf("✅ Publisher sincronizado exitosamente: %s (%s)", dto.DisplayName, dto.EntityType)
}

func (s *Subscriber) handleList(m *nats.Msg) {
	log.Printf("[NATS] publisher.list recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[ListPublishersDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := s.service.ListPublishers(ctx, dto)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al listar publishers: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, items)
}

func (s *Subscriber) handleCreate(m *nats.Msg) {
	log.Printf("[NATS] publisher.create recibido: %s", m.Data)

	packet, dto, err := parseNestPacket[CreatePlatformPublisherDTO](m.Data)
	if err != nil {
		replyNestError(m, "", err)
		return
	}
	if dto.DisplayName == "" {
		dto, err = parseFlat[CreatePlatformPublisherDTO](m.Data)
		if err != nil {
			replyNestError(m, packet.ID, err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := s.service.CreatePlatformPublisher(ctx, dto)
	if err != nil {
		log.Printf("[NATS ERROR] Fallo al crear publisher: %v", err)
		replyNestError(m, packet.ID, err)
		return
	}

	replyNestOK(m, packet.ID, created)
}

func parsePublisherPayload(data []byte) (UpdatePublisherDTO, error) {
	packet, dto, err := parseNestPacket[UpdatePublisherDTO](data)
	if err != nil {
		return UpdatePublisherDTO{}, err
	}
	if dto.EntityID != "" {
		return dto, nil
	}
	_ = packet
	return parseFlat[UpdatePublisherDTO](data)
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

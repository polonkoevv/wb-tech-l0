package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/stan.go"
	"github.com/polonkoevv/wb-tech/internal/models"
)

type ServiceInterface interface {
	Listen(context.Context, string) error
	LoadCache(context.Context) error
	GetAllFromCache() []models.Order
	GetFromCache(string) (models.Order, error)
}

type Service struct {
	conn    stan.Conn
	storage Repo
	Cache   map[string]models.Order
	logger  *slog.Logger
}

func New(conn stan.Conn, storage Repo, logger *slog.Logger) *Service {
	return &Service{
		conn:    conn,
		storage: storage,
		Cache:   make(map[string]models.Order),
		logger:  logger,
	}
}

func (s *Service) Listen(ctx context.Context, listenChannel string) error {
	//op := "service.listen"
	var err error

	sub, err := s.conn.Subscribe(listenChannel, func(m *stan.Msg) {
		data := models.Order{}
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			s.logger.Error("NATS: unmarshal:", err.Error())
			return
		}

		err = data.Validate()
		if err != nil {
			s.logger.Error("Data validation:", err.Error())
			return
		}

		for i := range data.Items {
			data.Items[i].OrderUid = data.OrderUid
		}

		err = s.storage.Save(ctx, data)
		if err != nil {
			s.logger.Error("Data saving:", err.Error())
			return
		}

		s.Cache[data.OrderUid] = data
	})
	if err != nil {
		s.logger.Error("NATS establishing connection:", err.Error())
		return err
	}

	defer func() {
		err := sub.Unsubscribe()
		if err != nil {
			s.logger.Error("NATS unsubscribing:", err.Error())
		}
		// Закрываем соединение
		err = s.conn.Close()
		if err != nil {
			s.logger.Error("NATS connection closing:", err.Error())
		}

	}()

	<-ctx.Done()

	return nil
}

func (s *Service) LoadCache(ctx context.Context) error {
	cache, err := s.storage.LoadCache(ctx)
	if err != nil {
		s.logger.Error("Loading cache", err.Error())
		return err
	}

	s.Cache = cache

	return nil
}

func (s *Service) GetAllFromCache() []models.Order {
	res := make([]models.Order, 0)
	for _, v := range s.Cache {
		res = append(res, v)
	}
	return res
}

func (s *Service) GetFromCache(uid string) (models.Order, error) {
	order, ok := s.Cache[uid]
	if !ok {
		return models.Order{}, fmt.Errorf("")
	}
	return order, nil
}

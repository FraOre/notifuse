package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Notifuse/notifuse/internal/domain"
	"github.com/Notifuse/notifuse/pkg/logger"
)

type EventDispatcherService struct {
	logger logger.Logger
}

func NewEventDispatcherService(
	logger logger.Logger,
) *EventDispatcherService {
	return &EventDispatcherService{
		logger: logger,
	}
}

func (s *EventDispatcherService) DispatchEvent(ctx context.Context, event *domain.WebhookEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		s.logger.WithField("error", err.Error()).Error("Failed to marshal webhook event for dispatch")
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://webhook.site/03d56604-eacf-41ca-acab-f1d6df9c583e", bytes.NewBuffer(payload))
	if err != nil {
		s.logger.WithField("error", err.Error()).Error("Failed to create HTTP request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		s.logger.WithField("error", err.Error()).Error("Failed to send event dispatcher")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.WithField("status_code", resp.StatusCode).Error("Event dispatch failed with non-2xx status")
	}
}
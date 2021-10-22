package mq

import (
	"context"
	"fmt"
	"github.com/igridnet/mproxy/logger"
	"github.com/igridnet/mproxy/pkg/session"
	"github.com/igridnet/users/api"
	"github.com/techcraftlabs/base/io"
	stdio "io"
	"strings"
	"time"
)


var _ session.Handler = (*Handler)(nil)

// Handler implements mqtt.Handler interface
type Handler struct {
	logger logger.Logger
	users *api.Client
	writer stdio.Writer
}

// New creates new Event entity
func New(logger logger.Logger,client *api.Client) *Handler {
	return &Handler{
		logger: logger,
		users: client,
		writer: io.Stderr,

	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (h *Handler) AuthConnect(c *session.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(),time.Minute)
	defer cancel()
	msg := fmt.Sprintf("AuthConnect() request- clientID: %s, username: %s, password: %s, client_CN: %s", c.ID, c.Username, string(c.Password), c.Cert.Subject.CommonName)
	_,_ = h.writer.Write([]byte(msg))
	node, err := h.users.GetNode(ctx,c.Username)
	if err != nil {
		msg := fmt.Sprintf("could not authenticate the node with id %s due to error: %v",c.Username,err)
		_,_ = h.writer.Write([]byte(msg))
		return err
	}

	if node.Key != string(c.Password){
		msg := fmt.Sprintf("password mismatch, not allowed")
		_,_ = h.writer.Write([]byte(msg))
		return err
	}
	return nil
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (h *Handler) AuthPublish(c *session.Client, topic *string, payload *[]byte) error {
	h.logger.Info(fmt.Sprintf("AuthPublish() - clientID: %s, topic: %s, payload: %s", c.ID, *topic, string(*payload)))
	return nil
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (h *Handler) AuthSubscribe(c *session.Client, topics *[]string) error {
	h.logger.Info(fmt.Sprintf("AuthSubscribe() - clientID: %s, topics: %s", c.ID, strings.Join(*topics, ",")))
	return nil
}

// Connect - after client successfully connected
func (h *Handler) Connect(c *session.Client) {
	h.logger.Info(fmt.Sprintf("Connect() - username: %s, clientID: %s", c.Username, c.ID))
}

// Publish - after client successfully published
func (h *Handler) Publish(c *session.Client, topic *string, payload *[]byte) {
	h.logger.Info(fmt.Sprintf("Publish() - username: %s, clientID: %s, topic: %s, payload: %s", c.Username, c.ID, *topic, string(*payload)))
}

// Subscribe - after client successfully subscribed
func (h *Handler) Subscribe(c *session.Client, topics *[]string) {
	h.logger.Info(fmt.Sprintf("Subscribe() - username: %s, clientID: %s, topics: %s", c.Username, c.ID, strings.Join(*topics, ",")))
}

// Unsubscribe - after client unsubscribed
func (h *Handler) Unsubscribe(c *session.Client, topics *[]string) {
	h.logger.Info(fmt.Sprintf("Unsubscribe() - username: %s, clientID: %s, topics: %s", c.Username, c.ID, strings.Join(*topics, ",")))
}

// Disconnect on conection lost
func (h *Handler) Disconnect(c *session.Client) {
	h.logger.Info(fmt.Sprintf("Disconnect() - client with username: %s and ID: %s disconenected", c.Username, c.ID))
}


package handler

import (
	"context"
)

type Handler interface {
	// Handle the claude code hook event and send the message through the telegram bot
	// The hook data is passed as a json string and you need to parse it first
	Handle(ctx context.Context, hookEvent string) error
}

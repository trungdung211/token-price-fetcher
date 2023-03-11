package external

import (
	"fmt"
	"sync"
	"time"

	"github.com/gtuk/discordwebhook"
	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
	"github.com/trungdung211/token-price-fetcher/internal/usecases/external"
	"go.uber.org/zap"
)

type discordAlert struct {
	l           *zap.Logger
	lastSend    map[string]time.Time
	delaySend   time.Duration
	botUserName string
	mu          sync.Mutex
}

func NewDiscordAlert(l *zap.Logger, botUserName string, delaySend time.Duration) external.ConditionAlert {
	return &discordAlert{
		l:           l,
		botUserName: botUserName,
		lastSend:    make(map[string]time.Time, 0),
		delaySend:   delaySend,
	}
}

func (d *discordAlert) Alert(config *model.UserConfig, message string) error {
	d.l.Debug("Alert", zap.Any("message", message))
	if config.DiscordWebhook != nil {
		d.mu.Lock()
		key := *config.DiscordWebhook
		lastDelay, found := d.lastSend[key]
		if found {
			lastDelay = lastDelay.Add(d.delaySend)
		}
		d.lastSend[key] = lastDelay
		d.mu.Unlock()

		go d.doSendAlert(*config.DiscordWebhook, message, lastDelay)
	}
	return nil
}

func (d *discordAlert) doSendAlert(discordWebhookUrl string, message string, scheduledAt time.Time) error {
	waitTime := time.Until(scheduledAt)
	if waitTime > 0 {
		time.Sleep(waitTime)
	}

	var username = d.botUserName
	var content = fmt.Sprintf(
		"\n ⚡Token Price Alert ⚡ \n\n%s\n",
		message,
	)

	dmessage := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(discordWebhookUrl, dmessage)
	if err != nil {
		d.l.Error("discordwebhook.SendMessage err", zap.Any("err", err))
	}
	return err
}

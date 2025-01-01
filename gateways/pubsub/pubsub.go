package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"registry-backend/config"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
)

type PubSubService interface {
	PublishNodePack(ctx context.Context, storageURL string) error
}

var _ PubSubService = (*pubsubimpl)(nil)

type pubsubimpl struct {
	client *pubsub.Client
	config *config.Config
	topic  *pubsub.Topic
}

func NewPubSubService(c *config.Config) (PubSubService, error) {
	if c == nil || c.PubSubTopic == "" {
		// Return a noop implementation if config is nil or storage is not enabled
		log.Info().Msg("No pub sub configuration found, using noop implementation")
		return &pubsubNoop{}, nil
	}

	// Initialize GCP storage client
	client, err := pubsub.NewClient(context.Background(), c.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("NewPubSubService: %v", err)
	}
	return &pubsubimpl{
		client: client,
		config: c,
		topic:  client.Topic(c.PubSubTopic),
	}, nil
}

// PublishNodePack implements PubSubService.
func (p *pubsubimpl) PublishNodePack(ctx context.Context, storageURL string) (err error) {
	u, err := url.Parse(storageURL)
	if err != nil {
		return fmt.Errorf("invalid storage URL: %w", err)
	}

	segments := strings.Split(u.Path, "/")
	if len(segments) < 2 {
		return fmt.Errorf("invalid storage URL: %w", err)
	}
	bucket := segments[1]
	object := strings.Join(segments[2:], "/")
	now := time.Now()
	messagePayload := map[string]interface{}{
		"kind":           "storage#object",
		"id":             fmt.Sprintf("%s/%s/%d", bucket, object, now.Unix()),
		"selfLink":       fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o/%s", object, bucket),
		"name":           object,
		"bucket":         bucket,
		"generation":     strconv.FormatInt(time.Now().Unix(), 10),
		"metageneration": "1",
		"mediaLink":      fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object),
	}
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}

	result := p.topic.Publish(ctx, &pubsub.Message{
		Data: jsonData,
		Attributes: map[string]string{
			"eventType": "OBJECT_FINALIZE", // Optional attribute for event type
		},
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Failed to publish message: %v", err)
	}

	return
}

var _ PubSubService = (*pubsubNoop)(nil)

type pubsubNoop struct{}

// PublishNodePack implements PubSubService.
func (p *pubsubNoop) PublishNodePack(ctx context.Context, storageURL string) error {
	return nil
}

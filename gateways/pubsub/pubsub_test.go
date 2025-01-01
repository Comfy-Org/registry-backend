package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"registry-backend/config"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublish(t *testing.T) {
	projectID, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		t.Skip("PROJECT_ID is not set")
	}

	client, err := pubsub.NewClient(context.Background(), projectID)
	require.NoError(t, err)

	topic := client.Topic(fmt.Sprintf("pubsub-topic-test-%d", time.Now().Unix()))
	t.Cleanup(func() {
		t.Logf("Deleting topic %s", topic.ID())
		topic.Delete(context.Background())
	})
	client.CreateTopic(context.Background(), topic.ID())

	pubsubsvc, err := NewPubSubService(&config.Config{ProjectID: projectID, PubSubTopic: topic.ID()})
	require.NoError(t, err)

	subscriptionID := fmt.Sprintf("sub-%d", time.Now().Unix())
	sub, err := client.CreateSubscription(context.Background(), subscriptionID, pubsub.SubscriptionConfig{
		Topic:               topic,
		AckDeadline:         10 * time.Second,
		RetainAckedMessages: false,
	})
	require.NoError(t, err)

	err = pubsubsvc.PublishNodePack(context.Background(), "https://storage.cloud.google.com/testbucket/path1/path2/file.tar.gz")
	require.NoError(t, err)
	pubsubsvc.(*pubsubimpl).topic.Flush()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data := map[string]string{}
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		t.Log(msg)
		assert.NoError(t, json.Unmarshal(msg.Data, &data))
		t.Log(data)
		msg.Ack()
		cancel()
	})
	<-ctx.Done()
	require.NoError(t, err)
	assert.Equal(t, "testbucket", data["bucket"])
	assert.Equal(t, "path1/path2/file.tar.gz", data["name"])
}

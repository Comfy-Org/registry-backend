package test

import (
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"cloud.google.com/go/storage"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	terraformDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "infrastructure/example/simple-node-pack-extractcd")
	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,
		Vars: map[string]interface{}{
			"prefix": "t" + strings.ToLower(random.UniqueId()),
		},
	}

	t.Cleanup(func() {
		terraform.Destroy(t, terraformOptions)
	})
	terraform.InitAndApply(t, terraformOptions)

	topicID := terraform.Output(t, terraformOptions, "topic_id")
	triggerID := terraform.Output(t, terraformOptions, "trigger_id")
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	bucketNotificationID := terraform.Output(t, terraformOptions, "bucket_notification_id")

	t.Run("CheckCloudBuildTrigger", func(t *testing.T) {
		gcb := gcp.NewCloudBuildService(t)
		res, err := gcb.GetBuildTrigger(context.Background(), &cloudbuildpb.GetBuildTriggerRequest{
			Name: triggerID,
		})
		require.NoError(t, err)
		assert.Equal(t, topicID, res.PubsubConfig.Topic)
	})

	t.Run("CheckBucketNotification", func(t *testing.T) {
		gcs, err := storage.NewClient(context.Background())
		require.NoError(t, err)

		notifications, err := gcs.Bucket(bucketName).Notifications(context.Background())
		require.NoError(t, err)
		found := false
		for id, notification := range notifications {
			notifyIDParts := strings.Split(bucketNotificationID, "/")
			topicIDParts := strings.Split(topicID, "/")
			t.Log(id, notification)
			if notification.TopicID != topicIDParts[len(topicIDParts)-1] || id != notifyIDParts[len(notifyIDParts)-1] {
				continue
			}
			found = true
		}
		assert.True(t, found)
	})
}

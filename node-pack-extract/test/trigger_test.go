package test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	cloudscheduler "cloud.google.com/go/scheduler/apiv1"
	"cloud.google.com/go/scheduler/apiv1/schedulerpb"
	"cloud.google.com/go/storage"

	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	terraformDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "infrastructure/examples/simple-node-pack-extract")
	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,
		Vars: map[string]interface{}{
			"prefix":     "t" + strings.ToLower(random.UniqueId()),
			"project_id": environment.GetFirstNonEmptyEnvVarOrFatal(t, []string{"GOOGLE_PROJECT", "GOOGLE_CLOUD_PROJECT", "GCLOUD_PROJECT", "CLOUDSDK_CORE_PROJECT"}),
		},
	}

	t.Cleanup(func() {
		terraform.Destroy(t, terraformOptions)
	})
	terraform.InitAndApply(t, terraformOptions)

	topicID := terraform.Output(t, terraformOptions, "topic_id")
	triggerID := terraform.Output(t, terraformOptions, "trigger_id")
	bucketNotificationID := terraform.Output(t, terraformOptions, "bucket_notification_id")
	schedulerID := terraform.Output(t, terraformOptions, "backfill_scheduler_id")
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	serviceAccount := terraform.Output(t, terraformOptions, "service_account")

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

	t.Run("CheckScheduler", func(t *testing.T) {
		client, err := cloudscheduler.NewCloudSchedulerClient(context.Background())
		require.NoError(t, err)
		defer client.Close()

		j, err := client.GetJob(context.Background(), &schedulerpb.GetJobRequest{
			Name: schedulerID,
		})
		require.NoError(t, err)
		h := j.GetHttpTarget()
		require.NotNil(t, h)
		assert.Contains(t, h.GetUri(), "/comfy-nodes/backfill")
		assert.Equal(t, http.MethodPost, h.GetHttpMethod().String())
		assert.Equal(t, "https://stagingapi.comfy.org", h.GetOidcToken().GetAudience())
		assert.Equal(t, serviceAccount, h.GetOidcToken().GetServiceAccountEmail())
	})
}

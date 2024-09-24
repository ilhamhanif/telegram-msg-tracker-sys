package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// https://cloud.google.com/storage/docs/json_api/v1/objects#resource-representations
type PubsubData struct {
	Kind                    string `json:"kind"`
	ID                      string `json:"id"`
	SelfLink                string `json:"selfLink"`
	MediaLink               string `json:"mediaLink"`
	Name                    string `json:"name"`
	Bucket                  string `json:"bucket"`
	Generation              string `json:"generation"`
	MetaGeneration          string `json:"metageneration"`
	ContentType             string `json:"contentType"`
	StorageClass            string `json:"storageClass"`
	Size                    string `json:"size"`
	SoftDeleteTime          string `json:"softDeleteTime"`
	HardDeleteTime          string `json:"hardDeleteTime"`
	Md5Hash                 string `json:"md5Hash"`
	ContentEncoding         string `json:"contentEncoding"`
	ContentDisposition      string `json:"contentDisposition"`
	ContentLanguage         string `json:"contentLanguage"`
	CacheControl            string `json:"cacheControl"`
	Crc32c                  string `json:"crc32c"`
	ComponentCount          string `json:"componentCount"`
	Etag                    string `json:"etag"`
	KmsKeyName              string `json:"kmsKeyName"`
	TemporaryHold           string `json:"temporaryHold"`
	EventBasedHold          string `json:"eventBasedHold"`
	RetentionExpirationTime string `json:"retentionExpirationTime"`
	Retention               struct {
		RetainUntilTime string `json:"retainUntilTime"`
		Mode            string `json:"mode"`
	} `json:"retention"`
	TimeCreated             string            `json:"timeCreated"`
	Updated                 string            `json:"updated"`
	TimeDeleted             string            `json:"timeDeleted"`
	TimeStorageClassUpdated string            `json:"timeStorageClassUpdated"`
	CustomTime              string            `json:"customTime"`
	Metadata                map[string]string `json:"metadata"`
	Acl                     []Acl             `json:"acl"`
	Owner                   struct {
		Entity   string `json:"entity"`
		EntityID string `json:"entityId"`
	} `json:"owner"`
	CustomerEncryption struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	} `json:"customerEncryption"`
}

type Acl struct {
	Kind        string `json:"kind"`
	Object      string `json:"object"`
	Generation  string `json:"generation"`
	ID          string `json:"id"`
	SelfLink    string `json:"selfLink"`
	Bucket      string `json:"bucket"`
	Entity      string `json:"entity"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	Domain      string `json:"domain"`
	EntityID    string `json:"entityId"`
	Etag        string `json:"etag"`
	ProjectTeam struct {
		ProjectNumber string `json:"projectNumber"`
		Team          string `json:"team"`
	} `json:"projectTeam"`
}

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC = "gcs_bucket_notif_log"

var attributes = map[string]string{
	"bucketId":           "test",
	"eventTime":          "2022-08-12T23:22:36.901891Z",
	"eventType":          "OBJECT_FINALIZE",
	"notificationConfig": "test",
	"objectGeneration":   "16603467",
	"objectId":           "itm/v",
	"payloadFormat":      "JSON_API_V1",
}

var pubsubData = PubsubData{
	Kind: "Test",
	ID:   "Test",
}

func (pd *PubsubData) publishToPubSub() error {

	// Setup PubSub client.
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error creating NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC)
	jsonData, err := json.Marshal(pd)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error encoding JSON: %w", err)
	}
	result := t.Publish(ctx, &pubsub.Message{
		Data:       jsonData,
		Attributes: attributes,
	})

	// Block until the result is returned
	// and a server-generated ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("publishToPubSub: Error publishing to PubSub: %w", err)
	}

	return nil

}

func main() {

	// Send 100 messages through Pub/Sub.
	for i := 0; i <= 100; i++ {

		if err := pubsubData.publishToPubSub(); err != nil {
			fmt.Printf("Error: %s", err)
		}
		fmt.Printf("%d\n", i)
	}

}

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type PubsubMessage struct {
	Attributes        map[string]string `json:"attributes"`
	Data              string            `json:"data"`
	MessageIDPascal   string            `json:"messageID"`
	MessageID         string            `json:"message_id"`
	PublishTimePascal string            `json:"publishTime"`
	PublishTime       string            `json:"publish_time"`
}

type PubsubSubscription struct {
	Message         *PubsubMessage `json:"message"`
	Subscription    string         `json:"subscription"`
	DeliveryAttempt int8           `json:"deliveryAttempt"`
}

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

const URL = "http://localhost:8080/GCSBucketNotifLogger"

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

var pubsubMessage = PubsubMessage{
	Attributes:  attributes,
	MessageID:   "5333919906745759",
	PublishTime: "2022-08-12T23:22:36.971Z",
}

var pubsubSubscription = PubsubSubscription{
	Message: &pubsubMessage,
}

func main() {

	// Setup message in JSON
	// mimic-ing real GCP Pub/Sub HTTP push message.
	messageJson, err := json.Marshal(pubsubData)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	pubsubSubscription.Message.Data = base64.StdEncoding.EncodeToString(messageJson)
	payloadJson, err := json.Marshal(pubsubSubscription)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// Sent the data to local endpoint
	// using HTTP POST.
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payloadJson))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	defer resp.Body.Close()

	// Print response and status code
	// given from the API.
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, string(body))

}

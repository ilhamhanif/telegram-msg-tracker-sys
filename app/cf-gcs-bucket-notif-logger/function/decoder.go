package function

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type Attributes struct {
	BucketID           string `json:"bucketId"`
	EventTime          string `json:"eventTime"`
	EventType          string `json:"eventType"`
	NotificationConfig string `json:"notificationConfig"`
	ObjectGeneration   string `json:"objectGeneration"`
	ObjectID           string `json:"objectId"`
	PayloadFormat      string `json:"payloadFormat"`
}

type PubsubMessage struct {
	Attributes        Attributes `json:"attributes"`
	Data              string     `json:"data"`
	MessageIDPascal   string     `json:"messageID"`
	MessageID         string     `json:"message_id"`
	PublishTimePascal string     `json:"publishTime"`
	PublishTime       string     `json:"publish_time"`
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

func (ps *PubsubSubscription) decodePubSubData(v *PubsubData) error {

	/*
		A method to Decode a PubSub data to v.
	*/

	// Convert (decode) string JSON.
	pubsubMessageDataDecoded, _ := base64.StdEncoding.DecodeString(ps.Message.Data)
	if err := json.Unmarshal(pubsubMessageDataDecoded, &v); err != nil {
		return fmt.Errorf("decodePubSubData: Error decoding PubSub Message: %w", err)
	}

	return nil

}

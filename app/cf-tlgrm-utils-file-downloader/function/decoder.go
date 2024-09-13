package function

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

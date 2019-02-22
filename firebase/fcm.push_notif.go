package firebase

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/TukTuk/common"

	fcm "github.com/appleboy/go-fcm"
)

func (fb *FireBase) SendPushNotification(ctx context.Context, data interface{}, deviceId, deviceType string) error {

	var err error

	client := fb.FBaseClient

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("[SendPushNotification][Error] Err in marshaling", err)
		return err
	}

	msg := &fcm.Message{
		To:         deviceId,
		TimeToLive: fb.Timeout,
	}

	if strings.EqualFold(deviceType, common.DEVICE_IOS) {
		msg.Notification = &fcm.Notification{
			Body: string(bytes),
		}
	}

	if strings.EqualFold(deviceType, common.DEVICE_ANDROID) || deviceType == "" {
		msg.Data = map[string]interface{}{
			"response": data,
		}
	}

	log.Printf("[SendPushNotification]Sending Push Notif:%+v data:%+v", data)

	res, err := client.Send(msg)
	if err != nil {
		log.Println("[SendPushNotification][Error] Err in sending notification", err)
		return err
	}

	log.Printf("Resp :%+v", res)
	return err
}

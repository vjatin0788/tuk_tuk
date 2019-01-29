package firebase

import (
	"context"
	"log"

	fcm "github.com/appleboy/go-fcm"
)

func (fb *FireBase) SendPushNotification(ctx context.Context, data interface{}, deviceId string) error {

	var err error

	client := fb.FBaseClient

	msg := &fcm.Message{
		To: deviceId,
		Data: map[string]interface{}{
			"response": data,
		},
		TimeToLive: fb.Timeout,
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

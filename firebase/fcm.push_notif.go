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

	// if res.StatusCode != http.StatusOK {
	// 	log.Println("[SendPushNotification][Error]Status code mismatch ", resp.StatusCode)
	// 	return errors.New("Statuscode mismatch")
	// }

	// if resp.Fail > 0 {
	// 	log.Println("[SendPushNotification][Error]Failed to send notification ", resp.Fail)
	// 	return errors.New("Fail TO send Notification")
	// }

	log.Printf("Resp :%+v", res)
	return err
}

package firebase

import (
	"context"
	"errors"
	"log"
	"net/http"
)

func (fb *FireBase) AddIds(ctx context.Context, ids []string) *FireBase {
	fb.Ids = append(fb.Ids, ids...)
	return fb
}

func (fb *FireBase) AddXds(ctx context.Context, xds []string) *FireBase {
	fb.Xds = append(fb.Xds, xds...)
	return fb
}

func (fb *FireBase) AddId(ctx context.Context, id string) *FireBase {
	fb.Ids = append(fb.Ids, id)
	return fb
}

func (fb *FireBase) AddXd(ctx context.Context, xd string) *FireBase {
	fb.Xds = append(fb.Xds, xd)
	return fb
}

func (fb *FireBase) SendPushNotification(ctx context.Context, data interface{}) error {

	var err error

	client := fb.FBaseClient

	//adding ids
	client.NewFcmRegIdsMsg(fb.Ids, data)

	//adding xds
	client.AppendDevices(fb.Xds)

	log.Println("[SendPushNotification]Sending Push Notif", data)

	resp, err := client.Send()
	if err != nil {
		log.Println("[SendPushNotification][Error] Err in sending notification", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("[SendPushNotification][Error]Status code mismatch ", resp.StatusCode)
		return errors.New("Statuscode mismatch")
	}

	// if resp.Fail > 0 {
	// 	log.Println("[SendPushNotification][Error]Failed to send notification ", resp.Fail)
	// 	return errors.New("Fail TO send Notification")
	// }

	log.Printf("Resp :%+v", resp)
	return err
}

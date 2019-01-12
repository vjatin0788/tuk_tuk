package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TukTuk/errs"

	"github.com/TukTuk/common"
)

func (pay *PaymentClient) InitiatePaymentRequest(ctx context.Context, rideId int64) (PaymentResponse, error) {
	var (
		paymentRes PaymentResponse
	)

	if rideId == 0 {
		log.Println("[InitiatePayment] Ride id empty")
		return paymentRes, errs.Err("Empty Error ID")
	}

	return pay.preparePayment(ctx, rideId)
}

func (pay *PaymentClient) preparePayment(ctx context.Context, rideId int64) (PaymentResponse, error) {
	var (
		err        error
		paymentRes PaymentResponse
	)

	log.Println("[preparePayment] Preparing request for getting payment")

	url := fmt.Sprintf("%s%s", pay.Cfg.Payment.Hostname, common.GET_PAYMENT)

	body, err := json.Marshal(&PaymentReq{
		RideId: rideId,
	})
	if err != nil {
		log.Println("[preparePayment][Error]Err in marshal:", err)
	}

	req, err := http.NewRequest(common.METHOD_POST, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("[preparePayment][Error]Err creating req ", err)
		return paymentRes, err
	}

	req.Header.Add("Content-Type", "application/json")
	log.Println("[preparePayment]req for payment ", req)

	resp, err := pay.Client.Do(req)
	if err != nil {
		log.Println("[preparePayment][Error]Err in resp ", err)
		return paymentRes, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("[preparePayment][Error]Status code mismatch ", resp.StatusCode)
		return paymentRes, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[preparePayment][Error]Err in bytes ", err)
		return paymentRes, err
	}

	err = json.Unmarshal(respBytes, &paymentRes)
	if err != nil {
		log.Println("[preparePayment][Error]Err in resp unmarshal ", err)
		return paymentRes, err
	}

	log.Printf("[preparePayment] Payment resp:%+v", paymentRes)

	return paymentRes, err
}

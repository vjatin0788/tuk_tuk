package maps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TukTuk/common"
)

func (mps *GMClient) GetDistance(ctx context.Context, destination, origin string) (DistanceMatrix, error) {

	var (
		defaultResp DistanceMatrix
		err         error
	)

	if destination == "" || origin == "" {
		log.Println("[GetDistance][Error] Empty destination or origin ")
		return defaultResp, errors.New("Empty destination or origin")
	}

	data, err := mps.prepareGetDistanceRequest(ctx, destination, origin, common.DRIVING_MODE)
	if err != nil {
		log.Println("[PrepareGetDistanceRequest][Error]Err creating req ", err)
		return defaultResp, err
	}

	return data, err
}

func (mps *GMClient) prepareGetDistanceRequest(ctx context.Context, destination, origins string, mode string) (DistanceMatrix, error) {

	log.Println("[PrepareGetDistanceRequest] Preparing request for getting distance")
	defaultRes := DistanceMatrix{}

	url := fmt.Sprintf("%s%s", mps.Cfg.Maps.Hostname, common.DISTANCE_MATRIX)

	req, err := http.NewRequest(common.METHOD_GET, url, nil)
	if err != nil {
		log.Println("[PrepareGetDistanceRequest][Error]Err creating req ", err)
		return defaultRes, err
	}

	log.Println("[PrepareGetDistanceRequest] Source, destin", origins, destination)

	query := req.URL.Query()
	query.Add("mode", mode)
	query.Add("origins", origins)
	query.Add("destinations", destination)
	query.Add("key", mps.Cfg.Maps.ApiKey)
	req.URL.RawQuery = query.Encode()

	log.Println("[PrepareGetDistanceRequest][Error]req for maps ", req)

	resp, err := mps.HttpClient.Do(req)
	if err != nil {
		log.Println("[PrepareGetDistanceRequest][Error]Err in resp ", err)
		return defaultRes, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("[PrepareGetDistanceRequest][Error]Status code mismatch ", resp.StatusCode)
		return defaultRes, err
	}

	//always use decoder in case of http req
	if err = json.NewDecoder(resp.Body).Decode(&defaultRes); err != nil && err != io.EOF {
		log.Println("[PrepareGetDistanceRequest][Error]Err unmarshaling resp", err)
		return defaultRes, err
	}
	return defaultRes, err
}

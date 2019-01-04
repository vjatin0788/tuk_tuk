package maps

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/TukTuk/common"
	gmaps "googlemaps.github.io/maps"
)

func TestGetDirection(t *testing.T) {
	ctx := context.Background()
	client, err := gmaps.NewClient(gmaps.WithAPIKey(common.API_KEY))
	if err != nil {
		log.Fatal("[InitMaps][Error] Err in creating new gmaps client", err)
	}

	mapsClient := &GMClient{
		Client: client,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	dist, err := mapsClient.GetDistance(ctx, " 30.210699,74.945281|30.223947,74.945037", "30.2178815,74.9443185")
	if err != nil {
		t.Errorf("[TestGetDirection] Test failed:%+v", err)
	}

	t.Log("Distance :", dist)
}

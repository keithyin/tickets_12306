package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestTrainoStation(t *testing.T) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	traino_stations := TrainoStations{}
	traino_stations.Init("G2557", ctx)

	fmt.Println(traino_stations.GetStationsBetweenFromTo("北京", "枣庄", true))

}

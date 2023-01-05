package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestShotgunMarriageStra(t *testing.T) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	from := "北京"
	to := "枣庄"
	date_str := "2023-01-19"

	// trainos := GetAllTrainoFromTo(ctx, from, to, date_str)
	trainos := []string{"K101", "G2565", "G2581"}

	traino_stations := make([]TrainoStations, 0)
	for _, traino := range trainos {
		traino_station := new(TrainoStations)
		traino_station.Init(traino, ctx)
		traino_stations = append(traino_stations, *traino_station)
	}

	shotgunMarriageStra := ShotgunMarriageStra{From: from, To: to, DateStr: date_str, TrainoStations: traino_stations}
	res := shotgunMarriageStra.Do(ctx)
	fmt.Println(res)

}

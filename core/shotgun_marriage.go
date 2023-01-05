package core

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

type ValidTainoInfo struct {
	From   string
	To     string
	Triano string
}

// 先上车后买票 策略
type ShotgunMarriageStra struct {
	From    string
	To      string
	DateStr string // 2022-10-10

	TrainoStations []TrainoStations
}

// 核心策略
func (receiver *ShotgunMarriageStra) Do(ctx context.Context) []ValidTainoInfo {

	valid_res := make([]ValidTainoInfo, 0)

	for n := 0; n < 10000; n++ {
		n_exceed_all_stations_len := true
		for _, traino_stations := range receiver.TrainoStations {
			stations := traino_stations.GetStationsBetweenFromTo(receiver.From, receiver.To, true)

			intermediate_station_num := len(stations)
			if n >= intermediate_station_num {
				continue
			}
			n_exceed_all_stations_len = false

			idx := intermediate_station_num - n - 1
			new_to := stations[idx]

			logrus.WithFields(logrus.Fields{
				"traino":          traino_stations.Traino,
				"lookup_stations": stations,
				"from":            receiver.From,
				"current_to":      new_to,
			}).Warn("check whether from -> to has ticket")

			valid := IsTicketValid(ctx, traino_stations.Traino, receiver.From, new_to, receiver.DateStr)
			if valid {
				valid_res = append(valid_res,
					ValidTainoInfo{From: receiver.From, To: new_to, Triano: traino_stations.Traino})
			}

		}

		if n_exceed_all_stations_len || len(valid_res) > 0 {
			break
		}
	}
	return valid_res

}

/*
*
from := "北京"

	to := "枣庄"
	date_str := "2023-01-19"
*/

func ShotgunMarriageStraPipeline(from, to, date_str string, trainos []string) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	if len(trainos) == 0 {
		trainos = GetAllTrainoFromTo(ctx, from, to, date_str)

	}

	traino_stations := make([]TrainoStations, 0)
	for _, traino := range trainos {
		traino_station := new(TrainoStations)
		traino_station.Init(traino, ctx)
		traino_stations = append(traino_stations, *traino_station)
	}
	// fmt.Println("traino_stations", traino_stations)

	shotgunMarriageStra := ShotgunMarriageStra{From: from, To: to, DateStr: date_str, TrainoStations: traino_stations}
	res := shotgunMarriageStra.Do(ctx)
	fmt.Println(res)
}

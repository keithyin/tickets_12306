package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)
}

type TrainoStations struct {
	Traino   string
	Stations []string
}

// https://trains.ctrip.com/TrainSchedule/G2557 通过该链接查找 车次的站点
func (receiver *TrainoStations) Init(traino string, ctx context.Context) error {
	receiver.Traino = traino
	receiver.Stations = make([]string, 0)
	target_url := fmt.Sprintf("https://trains.ctrip.com/TrainSchedule/%s", traino)

	body_selector := `#ctl00_MainContentPlaceHolder_pnlResult > div.s_bd > table.tb_result.tb_inquiry.tb_gray > tbody`
	stations_sel := `#ctl00_MainContentPlaceHolder_pnlResult > div.s_bd > table.tb_result.tb_inquiry.tb_gray > tbody > tr > td:nth-child(3)`

	js_stations := fmt.Sprintf(`[...document.querySelectorAll("%s")].map((e) => e.innerText)`, stations_sel)

	err := chromedp.Run(ctx,
		chromedp.Navigate(target_url),

		chromedp.WaitVisible(body_selector),
		chromedp.Evaluate(js_stations, &receiver.Stations),
	)

	if err != nil {
		logrus.Fatal(err)
	}
	return err
}

func (receiver *TrainoStations) GetStationsBetweenFromTo(from string, to string, right_inclusive bool) []string {
	if receiver.Traino == "" {
		logrus.Fatal("run .Init first")
	}

	results := make([]string, 0)
	do_record := false
	for _, station := range receiver.Stations {
		if strings.HasPrefix(station, to) {
			do_record = false
			break
		}

		if do_record {
			results = append(results, station)
		}

		if strings.HasPrefix(station, from) {
			do_record = true
		}
	}
	if right_inclusive {
		results = append(results, to)
	}
	return results
}

func (receiver *TrainoStations) GetStationsAfterTo(to string) []string {
	if receiver.Traino == "" {
		logrus.Fatal("run .Init first")
	}

	results := make([]string, 0)
	do_record := false
	for _, station := range receiver.Stations {
		if do_record {
			results = append(results, station)
		}

		if strings.HasPrefix(station, to) {
			do_record = true
		}
	}
	return results
}

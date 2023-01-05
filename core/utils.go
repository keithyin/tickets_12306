package core

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func select_station_func(target_station_name string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// #citem_2 > span:nth-child(1)
		founded_query_id := ""
		for i := 0; i < 10; i++ {
			query_id := fmt.Sprintf(`#citem_%d > span:nth-child(1)`, i)
			js_script := fmt.Sprintf(`document.querySelector("%s").innerText`, query_id)
			station_name := ""
			err := chromedp.Evaluate(js_script, &station_name).Do(ctx)
			if err != nil {
				fmt.Println(err)
				return err
			}
			// fmt.Println(station_name)
			if station_name == target_station_name {
				founded_query_id = query_id
				// fmt.Println("found!!!")
				break
			}
		}

		err := chromedp.Click(founded_query_id).Do(ctx)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	})
}

func IsTicketValid(ctx context.Context, traino string, from string, to string, date_str string) bool {
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.12306.cn/index/"),
		chromedp.WaitVisible(`#fromStationText`),
		chromedp.Click(`#fromStationText`),
		chromedp.SendKeys(`#fromStationText`, from),
		select_station_func(from),

		chromedp.SendKeys(`#toStationText`, to),
		select_station_func(to),
		chromedp.SetValue(`#train_date`, date_str),
		chromedp.Click(`#search_one`),
		chromedp.Sleep(100*time.Millisecond),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	// https://github.com/chromedp/chromedp/issues/656#issuecomment-836264205
	tabs, err := chromedp.Targets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(tabs) == 0 {
		log.Fatal("no tabs")
	}

	var let_ticket_tab *target.Info
	for _, tab := range tabs {
		if strings.Contains(tab.URL, "kyfw.12306.cn") {
			let_ticket_tab = tab
			break
		}
	}
	if let_ticket_tab == nil {
		return false
	}

	tabCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(let_ticket_tab.TargetID))
	defer cancel()

	valid := false

	err = chromedp.Run(tabCtx,
		chromedp.WaitVisible(`#trainum`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			train_num_str := ""
			err := chromedp.Evaluate(`document.querySelector("#trainum").innerText`, &train_num_str).Do(ctx)
			if err != nil {
				return err
			}
			train_num, err := strconv.Atoi(train_num_str)
			if err != nil {
				return err
			}

			for i := 0; i < train_num; i++ {
				query_id := fmt.Sprintf(`#train_num_%d > div > div > a`, i)
				train_no_str := ""
				err = chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").innerText`, query_id), &train_no_str).Do(ctx)
				if err != nil {
					return err
				}
				// fmt.Println("train_no -> ", train_no_str)

				ticket_id := ""
				chromedp.QueryAfter(fmt.Sprintf(`#train_num_%d`, i), func(ctx context.Context, eci runtime.ExecutionContextID, n ...*cdp.Node) error {
					node := n[0]
					parent_id, _ := node.Parent.Parent.Attribute("id")
					ticket_id = parent_id
					return nil
				}).Do(ctx)

				if train_no_str == traino {
					query_id = fmt.Sprintf(`#%s > td.no-br > a`, ticket_id)
					var nodes []*cdp.Node
					err = chromedp.Nodes(query_id, &nodes, chromedp.AtLeast(0)).Do(ctx)
					if err != nil {
						return err
					}
					if len(nodes) > 0 {
						valid = true
					}
					break
				}

			}
			return err
		}),
	)

	if err != nil {
		return false
	}

	return valid
}

func GetAllTrainoFromTo(ctx context.Context, from string, to string, date_str string) []string {

	result_trianos := make([]string, 0)

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.12306.cn/index/"),
		chromedp.WaitVisible(`#fromStationText`),
		chromedp.Click(`#fromStationText`),
		chromedp.SendKeys(`#fromStationText`, from),
		select_station_func(from),

		chromedp.SendKeys(`#toStationText`, to),
		select_station_func(to),
		chromedp.SetValue(`#train_date`, date_str),
		chromedp.Click(`#search_one`),
		chromedp.Sleep(100*time.Millisecond),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	// https://github.com/chromedp/chromedp/issues/656#issuecomment-836264205
	tabs, err := chromedp.Targets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(tabs) == 0 {
		log.Fatal("no tabs")
	}

	var let_ticket_tab *target.Info
	for _, tab := range tabs {
		if strings.Contains(tab.URL, "kyfw.12306.cn") {
			let_ticket_tab = tab
			break
		}
	}
	if let_ticket_tab == nil {
		return result_trianos
	}

	tabCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(let_ticket_tab.TargetID))
	defer cancel()

	err = chromedp.Run(tabCtx,
		chromedp.WaitVisible(`#trainum`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			train_num_str := ""
			err := chromedp.Evaluate(`document.querySelector("#trainum").innerText`, &train_num_str).Do(ctx)
			if err != nil {
				return err
			}
			train_num, err := strconv.Atoi(train_num_str)
			if err != nil {
				return err
			}

			for i := 0; i < train_num; i++ {
				query_id := fmt.Sprintf(`#train_num_%d > div > div > a`, i)
				train_no_str := ""
				err = chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").innerText`, query_id), &train_no_str).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Println("train_no -> ", train_no_str)
				result_trianos = append(result_trianos, train_no_str)

			}
			return err
		}),
	)

	if err != nil {
		return result_trianos
	}

	return result_trianos
}

package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
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
			fmt.Println(station_name)
			if station_name == target_station_name {
				founded_query_id = query_id
				fmt.Println("found!!!")
				break
			}
		}

		err := chromedp.Click(founded_query_id).Do(ctx)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		return nil
	})
}

func main() {

	from_station := "北京"
	to_station := "徐州"
	expected_train_no := "Z155"
	// expected_train_no := "G103"

	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.12306.cn/index/"),
		chromedp.WaitVisible(`#fromStationText`),
		chromedp.Click(`#fromStationText`),
		chromedp.SendKeys(`#fromStationText`, from_station),
		select_station_func(from_station),

		chromedp.SendKeys(`#toStationText`, "徐州"),
		select_station_func(to_station),
		chromedp.SetValue(`#train_date`, "2023-01-18"),
		chromedp.Click(`#search_one`),
		chromedp.Sleep(1*time.Second),
		// chromedp.WaitVisible(`#train_num_0`), // #\32 40000G10335_VNP_UUH
		// chromedp.WaitVisible(`#trainum`),

		// chromedp.ActionFunc(func(ctx context.Context) error {

		// 	train_num := ""
		// 	err := chromedp.QueryAfter(`#trainum`, func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		// 		node := nodes[0]
		// 		train_num = node.Value
		// 		return nil
		// 	}).Do(ctx)

		// 	fmt.Println(train_num)
		// 	return err
		// 	// query_id := fmt.Sprintf(`#train_num_%d > div > div > a`)
		// }),

		// chromedp.QueryAfter(`#fromStationText`, )
		// #ticket_0h0000Z15803_05_07 > td.no-br > a
		// #train_num_55

		// chromedp.QueryAfter(`#fromStationText`, func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		// 	node := nodes[0]
		// 	quards, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	quard := quards[0]
		// 	left_up_x := quard[0]
		// 	left_up_y := quard[1]

		// 	p := &input.DispatchMouseEventParams{
		// 		Type:   input.MousePressed,
		// 		X:      float64(left_up_x + 10),
		// 		Y:      float64(left_up_y + 5),
		// 		Button: input.Left,
		// 	}
		// 	err = p.Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	time.Sleep(4 * time.Second)
		// 	return nil
		// }),

		// #ul_list1 > li.ac_even.openLi.ac_over
		// #citem_0

	)

	if err != nil {
		log.Fatal(err)
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
			fmt.Println("tab -> ", tab.URL)
			break
		}
	}
	if let_ticket_tab == nil {
		return
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
			fmt.Println("train_num->", train_num_str)
			train_num, err := strconv.Atoi(train_num_str)

			for i := 0; i < train_num; i++ {
				query_id := fmt.Sprintf(`#train_num_%d > div > div > a`, i)
				train_no_str := ""
				err = chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").innerText`, query_id), &train_no_str).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Println("train_no -> ", train_no_str)

				ticket_id := ""
				chromedp.QueryAfter(fmt.Sprintf(`#train_num_%d`, i), func(ctx context.Context, eci runtime.ExecutionContextID, n ...*cdp.Node) error {
					node := n[0]
					parent_id, _ := node.Parent.Parent.Attribute("id")
					ticket_id = parent_id
					return nil
				}).Do(ctx)

				// 确定了之后，还需要滑动到相应位置，然后点击才会有用，否则是没有用的！！！！
				if train_no_str == expected_train_no {
					query_id = fmt.Sprintf(`#%s > td.no-br > a`, ticket_id)
					fmt.Println(query_id)

					// 似乎这个滑动是没有必要的！！！！
					for i := 0; i < 0; i++ {
						nodes := make([]*cdp.Node, 0)
						err = chromedp.Nodes(query_id, &nodes, chromedp.AtLeast(0), chromedp.NodeVisible).Do(ctx)
						if err != nil {
							return err
						}
						if len(nodes) > 0 {
							break
						}

						p := &input.DispatchMouseEventParams{
							Type:   input.MouseWheel,
							X:      0,
							Y:      0,
							DeltaX: 10,
							DeltaY: 10,
						}
						err = p.Do(ctx)
						if err != nil {
							return err
						}
					}

					err = chromedp.QueryAfter(query_id, func(ctx context.Context, eci runtime.ExecutionContextID, n ...*cdp.Node) error {
						node := n[0]
						quads, err := dom.GetContentQuads().WithNodeID(node.NodeID).Do(ctx)
						if err != nil {
							return err
						}
						quad := quads[0]
						left_up_x := quad[0]
						left_up_y := quad[1]
						right_down_x := quad[2]
						right_down_y := quad[3]
						fmt.Printf("left_up_x:%f, left_up_y:%f, right_down_x:%f, right_down_y:%f\n", left_up_x, left_up_y, right_down_x, right_down_y)
						p := &input.DispatchMouseEventParams{
							Type:   input.MouseWheel,
							X:      0,
							Y:      0,
							DeltaX: (left_up_x + right_down_x) / 2,
							DeltaY: (left_up_y + right_down_y) / 2,
						}
						err = p.Do(ctx)
						if err != nil {
							return err
						}
						return nil

					}).Do(ctx)
					if err != nil {
						return err
					}

					// js_script := fmt.Sprintf(`document.querySelector("%s").click()`, query_id)
					// fmt.Println(js_script)
					// chromedp.Evaluate(js_script, &res).Do(ctx)

					err = chromedp.Click(query_id).Do(ctx)
					if err != nil {
						return err
					}

				}

			}

			// err = chromedp.QueryAfter(`#trainum`, func(ctx context.Context, id runtime.ExecutionContextID, nodes ...*cdp.Node) error {
			// 	node := nodes[0]
			// 	fmt.Println("nodeVal -> ", node.Children[0].Value)
			// 	return nil
			// }).Do(ctx)

			return err
			// query_id := fmt.Sprintf(`#train_num_%d > div > div > a`)
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	fmt.Println("hello")
}

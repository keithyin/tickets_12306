package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func TestIsTicketValid(t *testing.T) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	logrus.Warn("valid=", IsTicketValid(ctx, "G197", "北京", "徐州", "2023-01-19"))
}

func TestGetAllTrainoFromTo(t *testing.T) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	fmt.Println(GetAllTrainoFromTo(ctx, "北京", "枣庄", "2023-01-19"))

}

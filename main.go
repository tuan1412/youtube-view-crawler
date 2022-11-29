package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("start-fullscreen", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("remote-debugging-port", "9222"),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	expandBtn := "//*[@id=\"description-inner\"]"

	task := chromedp.Tasks{
		chromedp.Navigate("https://www.youtube.com/watch?v=Rr5bP7uLnfk"),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(expandBtn),
		chromedp.Sleep(1 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return err
			}

			doc.Find("#info > span:nth-child(1)").Each(func(index int, info *goquery.Selection) {
				text := info.Text()
				fmt.Println(text)
			})

			return nil
		}),
	}

	_ = chromedp.Run(ctx, task)

}

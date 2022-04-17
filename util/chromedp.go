package util

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func GetUserDataDir() string {
	// Get Folder User Data Directory
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%v/.config/google-chrome/", dir)
}

func taskgetLinks(url string, titleNode *[]*cdp.Node, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(url),
		chromedp.Click("div#sponsor-button > ytd-button-renderer > a", chromedp.NodeVisible),
		chromedp.Nodes("yt-formatted-string.channel-title", titleNode, chromedp.ByQuery),
		chromedp.Nodes("yt-img-shadow.ytd-sponsorships-perk-renderer > img", nodes, chromedp.ByQueryAll, chromedp.NodeVisible),
	}
}

func RunTaskGetLinks(ctx context.Context, url string) {
	var titleNode []*cdp.Node
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, taskgetLinks(url, &titleNode, &nodes)); err != nil {
		log.Fatal(err)
	}
	title := titleNode[0].Children[0].NodeValue
	fmt.Printf("get links from channel: %s\n", title)
	CreateFileLinks(title, nodes)
}

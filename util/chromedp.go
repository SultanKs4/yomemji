package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

func taskgetLinks(url string, titleNode *[]*cdp.Node, linkNodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(url),
		chromedp.Nodes("ytd-channel-name > #container > #text-container > #text", titleNode, chromedp.ByQuery, chromedp.NodeVisible),
		chromedp.Click("div#sponsor-button > ytd-button-renderer > a", chromedp.NodeVisible),
		chromedp.Nodes("yt-img-shadow.ytd-sponsorships-perk-renderer > img", linkNodes, chromedp.ByQueryAll, chromedp.NodeVisible),
	}
}

func RunTaskGetLinks(ctx context.Context, url string) (string, []string) {
	var titleNode []*cdp.Node
	var linkNodes []*cdp.Node
	if err := chromedp.Run(ctx, taskgetLinks(url, &titleNode, &linkNodes)); err != nil {
		log.Fatal(err)
	}
	title := titleNode[0].Children[0].NodeValue
	fmt.Printf("get links from channel: %s\n", title)
	title = strings.Replace(title, "/", "", -1)

	links := []string{}
	for _, v := range linkNodes {
		src := v.AttributeValue("src")
		if src == "" {
			continue
		}
		links = append(links, src)
	}

	return title, links
}

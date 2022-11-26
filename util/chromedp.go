package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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

func navigateTask(url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(url),
	}
}

func getTitleChannel(titleNode *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Nodes("ytd-channel-name > #container > #text-container > #text", titleNode, chromedp.ByQuery, chromedp.NodeVisible),
	}
}

func getLinks(linkNodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Click("[aria-label='Join this channel']", chromedp.NodeVisible),
		chromedp.Nodes("yt-img-shadow.ytd-sponsorships-perk-renderer > img", linkNodes, chromedp.ByQueryAll, chromedp.NodeVisible),
	}
}

func RunTaskGetLinks(ctx context.Context, wg *sync.WaitGroup, url string) {
	newTabCtx, cancelTab := chromedp.NewContext(ctx)
	title := ""
	//Cancel to release resources once the function is complete
	defer func() {
		wg.Done()
		log.Printf("close tab channel: %s\n", title)
		cancelTab()
	}()
	var titleNode []*cdp.Node
	var linkNodes []*cdp.Node
	ctxTo, cancelTo := context.WithTimeout(newTabCtx, 50*time.Second)
	defer cancelTo()
	if err := chromedp.Run(ctxTo, navigateTask(url), getTitleChannel(&titleNode)); err != nil {
		log.Fatalf("failed load youtube channel: %s: %s", url, err)
	}
	// sanitize folder name linux tested linux
	// TODO: Add Windows + MacOs
	title = strings.Replace(titleNode[0].Children[0].NodeValue, "/", "", -1)

	ctxTo2, cancelTo2 := context.WithTimeout(newTabCtx, 10*time.Second)
	defer cancelTo2()
	if err := chromedp.Run(ctxTo2, getLinks(&linkNodes)); err != nil {
		log.Printf("skip channel %s", title)
	}
	links := []string{}
	for _, v := range linkNodes {
		src := v.AttributeValue("src")
		if src == "" {
			continue
		}
		links = append(links, src)
	}
	CreateFileLinks(title, links)
}

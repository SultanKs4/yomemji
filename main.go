package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/SultanKs4/yomemji/util"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	getLinksFromChannelUrl()
	getUrlArray(util.GetListFile())
}

func getUrlArray(pathArray []string) {
	for _, v := range pathArray {
		data, err := os.ReadFile(fmt.Sprintf("links/%s", v))
		if err != nil {
			log.Fatal(err)
		}
		stringData := strings.Split(string(data), "\n")
		stringData = stringData[:len(stringData)-1]

		folderName := strings.Split(v, ".txt")[0]
		path := fmt.Sprintf("images/%s", folderName)
		os.MkdirAll(path, os.ModePerm)
		fmt.Printf("download image for folder: %s\n", folderName)
		util.DownloadImgUrl(path, stringData)
	}
}

func getLinksFromChannelUrl() {
	f, err := os.Open("channel.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	urlChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text := scanner.Text()
			matched, err := regexp.MatchString("^https", text)
			if err != nil {
				log.Fatal(err)
			}
			if matched {
				urlChan <- text
			}
			continue
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		close(urlChan)
	}()

	allocCtx, cancelRoot := chromedp.NewExecAllocator(context.Background(), chromedp.UserDataDir(util.GetUserDataDir()), chromedp.Flag("headless", true))
	defer cancelRoot()

	taskCtx, cancelFirst := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancelFirst()
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for v := range urlChan {
		wg.Add(1)
		newTabCtx, ccl := chromedp.NewContext(taskCtx)
		defer ccl()
		go util.RunTaskGetLinks(newTabCtx, &wg, v)
	}

	// close first tab
	if err := page.Close().Do(cdp.WithExecutor(taskCtx, chromedp.FromContext(taskCtx).Target)); err != nil {
		panic(err)
	}

	wg.Wait()
}

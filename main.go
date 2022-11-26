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
	f, err := os.Open("channel.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	urlChan := make(chan string)
	go func(ch chan string) {
		scanner := bufio.NewScanner(f)
		re := regexp.MustCompile(`^\| (\[.*\])\((.*)\)`)
		for scanner.Scan() {
			text := scanner.Text()
			if matches := re.FindStringSubmatch(text); len(matches) != 0 {
				ch <- matches[2]
			}
			continue
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		close(ch)
	}(urlChan)

	options := []chromedp.ExecAllocatorOption{
		chromedp.UserDataDir(util.GetUserDataDir()),
		chromedp.Flag("headless", true),
		chromedp.Flag("mute-audio", true),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	allocCtx, cancelRoot := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancelRoot()

	taskCtx, cancelFirst := chromedp.NewContext(allocCtx, chromedp.WithErrorf(log.Printf))
	defer cancelFirst()
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for v := range urlChan {
		wg.Add(1)
		go util.RunTaskGetLinks(taskCtx, &wg, v)
	}

	// close first tab
	if err := page.Close().Do(cdp.WithExecutor(taskCtx, chromedp.FromContext(taskCtx).Target)); err != nil {
		panic(err)
	}

	wg.Wait()
}

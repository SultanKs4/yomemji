package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SultanKs4/yomemji/util"
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
		os.Mkdir(path, os.ModePerm)
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

	scanner := bufio.NewScanner(f)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.UserDataDir(util.GetUserDataDir()), chromedp.Flag("headless", false))
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	for scanner.Scan() {
		util.RunTaskGetLinks(taskCtx, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

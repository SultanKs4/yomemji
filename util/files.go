package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func CreateFileLinks(title string, links []string) {
	if len(links) == 0 {
		return
	}
	if err := os.MkdirAll("links", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(fmt.Sprintf("links/%s.txt", title), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var wg sync.WaitGroup
	for _, v := range links {
		wg.Add(1)

		urlChan := make(chan string)
		go func(url string) {
			defer wg.Done()
			// change reso badges
			url = strings.Replace(url, "s64", "s512", -1)
			// change reso emoji
			url = strings.Replace(url, "w48-h48", "w448-h448", -1)
			urlChan <- url
		}(v)
		url := <-urlChan

		if _, err := f.Write([]byte(fmt.Sprintf("%s\n", url))); err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
}

func DownloadImgUrl(path string, urls []string) {
	var wg sync.WaitGroup

	for i, v := range urls {
		wg.Add(1)
		go func(path string, filename string, url string) {
			defer wg.Done()

			output, err := os.Create(fmt.Sprintf("%s/%s", path, filename))
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			_, err = io.Copy(output, res.Body)
			if err != nil {
				log.Fatal(err)
			}
		}(path, fmt.Sprintf("%s.png", strconv.Itoa(i+1)), v)
	}
	wg.Wait()
}

func GetListFile() []string {
	f, err := os.Open("links")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	var fileArr []string
	for _, v := range files {
		if v.IsDir() {
			fmt.Printf("skip %s because directory", v.Name())
			continue
		}
		fileArr = append(fileArr, v.Name())
	}
	return fileArr
}

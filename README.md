# YOMEMJI

Get links Youtube Channel Membership Badges and Emoji then save it to local disk as images

## Why :question:

The idea is to retrieve image badges and custom emojis from channel youtube automated and make it for personal uses e.g. *meme material*, *emoji discord*

## How it works :question:

I am using [chromedp](https://github.com/chromedp/chromedp) for get all links and http package for download image and save it to local disk.

Step:

1. Read all url channel youtube at `channel.txt`.
2. Open url channel via chrome with user data browser that already login to youtube.
3. click `Join` button.
4. Get all images link.
5. Save it to `txt` files.
6. Download images based on images links saved from before.

## How to use :question:

1. List url channel url in `channel.txt`
2. Change folder User Data Directory if needed
3. run file using command `go run main.go`

## Limitation :construction:

1. Need Chrome Browser because using `chromedp`
2. Only tested in `WSL2 Ubuntu` so the code only load Chrome Data Directory Ubuntu.
3. For exclude login google account etc process i using **Chrome User Data** from browser that already logged in manually.
4. When running this program it will open chrome as `non-headless` because idk when using `headless` mode only load badges links not emoji links.

## Disclaimer :warning:

This project only for personal uses and education purpose.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/antchfx/htmlquery"
)

func main() {
	startUrl := ""
	rootFolder := ""

	imageUrls, err := getImageUrls(startUrl)
	if err != nil {
		panic(err)
	}

	for _, link := range imageUrls {
		downloadFile(link, rootFolder)
	}
}

func getImageUrls(startUrl string) ([]string, error) {
	var imageUrls []string
	for true {
		doc, err := htmlquery.LoadURL(startUrl)
		if err != nil {
			return nil, fmt.Errorf("load html %s: %w", startUrl, err)
		}

		node, err := htmlquery.Query(doc, "//div[@id='pages']/span/text()")
		if err != nil {
			return nil, fmt.Errorf("find current page: %w", err)
		}
		currentPageIndex := node.Data
		log.Println("current page:", currentPageIndex)

		nodes, err := htmlquery.QueryAll(doc, "//img[@class='tupian_img']")
		if err != nil {
			return nil, fmt.Errorf("find images %w", err)
		}
		for _, n := range nodes {
			for _, a := range n.Attr {
				if a.Key == "src" {
					imageUrls = append(imageUrls, a.Val)
					log.Println(a.Val)
				}
			}
		}

		node, err = htmlquery.Query(doc, "//div[@id='pages']//a[text()='下一页']")
		if err != nil {
			panic(err)
		}
		var href string
		for _, a := range node.Attr {
			if a.Key == "href" {
				href = a.Val
			}
		}
		ss := strings.Split(href, "/")
		nextPageIndex := strings.Split(ss[len(ss)-1], ".")[0]
		startUrl = href

		log.Println("next page url:", href)
		log.Println("next page:", nextPageIndex)

		if currentPageIndex == nextPageIndex {
			break
		}

	}
	return imageUrls, nil
}

func downloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ss := strings.Split(url, "/")
	if len(ss) == 0 {
		return fmt.Errorf("no file to download")
	}

	fileName := ss[len(ss)-1]

	fullPath := path.Join(filePath, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	log.Println("downloaded", fullPath)
	return nil
}

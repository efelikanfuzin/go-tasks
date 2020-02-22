package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
)

func startWith(str string, char rune) bool {
	return str[0] == byte(char)
}

func downlaodImage(link string, index int, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(link)
	if err != nil {
		log.Fatal("Error when download image", err)
	}
	defer resp.Body.Close()

	file, err := os.Create("./images/image_" + strconv.Itoa(index) + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}

func findAllImageLinks(content []byte) []string {
	var image_link = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
	images_links_match := image_link.FindAllStringSubmatch(string(content), -1)
	raw_links := make([]string, len(images_links_match))

	for i := range raw_links {
		if startWith(images_links_match[i][1], '/') {
			raw_links[i] = "https://tehnoholod.ru" + images_links_match[i][1]
		} else {
			raw_links[i] = images_links_match[i][1]
		}
	}

	return raw_links
}

func main() {
	url := "https://tehnoholod.ru"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error when open site", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	links := findAllImageLinks(body)

	var wg sync.WaitGroup
	wg.Add(len(links))

	for i, link := range links {
		fmt.Println(link)
		go downlaodImage(link, i, &wg)
	}
	wg.Wait()
	fmt.Printf("Download %d images", len(links))
}

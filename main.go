package main

import (
	"bufio"
	"fmt"
	"github.com/opesun/goquery"
	"io/ioutil"
	"log"
	"os"
)

type Title struct {
	title string
	link  string
}

func (t Title) String() string {
	return fmt.Sprintf("title: %s\tlink: %s", t.title, t.link)
}

func init() {

}

func main() {
	file, err := os.OpenFile("links.txt", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		requestUrl := scanner.Text()
		fmt.Println(requestUrl)

		//data, err := ioutil.ReadFile("out.html")
		//if err != nil {
		//	log.Fatal(err)
		//}

		//nodes, err := goquery.ParseString(string(data))
		nodes, err := goquery.ParseUrl(requestUrl)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile("out.html", []byte(nodes.Html()), 0644)

		foundNodes := nodes.Find(".flex-dow-txt")
		fmt.Println(foundNodes[0].Data)
		foundNodes = foundNodes.Find(".title")

		titles := make([]Title, 0, len(foundNodes))
		for _, node := range foundNodes {
			title := getTitle(requestUrl, node)
			if title.title != "" {
				titles = append(titles, title)
			}
		}
		for i, title := range titles {
			fmt.Printf("%d: %s\n", i, title)
		}

	}
}

func getTitle(requestUrl string, node *goquery.Node) Title {
	for _, attr := range node.Attr {
		if attr.Val == "title" {
			hrefNode := node.Child[0]
			title := hrefNode.Attr[1].Val
			link := hrefNode.Attr[0].Val
			linkRunesSlice := []rune(link)
			link = string(linkRunesSlice[1:])
			return Title{title: title, link: requestUrl + link}
		}
	}
	return Title{}
}

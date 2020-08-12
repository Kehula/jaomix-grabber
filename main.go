package main

import (
	"bufio"
	"fmt"
	"github.com/opesun/goquery"
	"io/ioutil"
	"log"
	"os"
	"sync"
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

		nodes, err := goquery.ParseUrl(requestUrl)
		if err != nil {
			log.Fatal(err)
		}
		writeNodes2File("out.html", nodes)

		foundNodes := nodes.Find(".flex-dow-txt")
		fmt.Println(foundNodes[0].Data)
		foundNodes = foundNodes.Find(".title")

		titles := make([]Title, 0, len(foundNodes))
		for _, node := range foundNodes {
			title := getTitle("https://jaomix.ru", node)
			if title.title != "" {
				titles = append(titles, title)
			}
		}
		parseChapters(titles)
	}
}

func writeNodes2File(filename string, nodes goquery.Nodes) {
	err := ioutil.WriteFile(filename, []byte(nodes.Html()), 0664)
	if err != nil {
		log.Fatal(err)
	}
}

func getTitle(requestUrl string, node *goquery.Node) Title {
	for _, attr := range node.Attr {
		if attr.Val == "title" {
			hrefNode := node.Child[0]
			title := hrefNode.Attr[1].Val
			link := hrefNode.Attr[0].Val
			return Title{title: title, link: requestUrl + link}
		}
	}
	return Title{}
}

func parseChapters(titles []Title) {
	createDir("chapters/")
	titles = lookup4NewChapters(titles)

	var wg sync.WaitGroup
	for _, title := range titles {
		wg.Add(1)
		go func(sync *sync.WaitGroup, title Title) {
			defer sync.Done()
			fmt.Println(title.link)
			nodes, err := goquery.ParseUrl(title.link)
			if err != nil {
				log.Fatal(err)
			}
			removeNodesList(&nodes, ".adsbygoogle", "header-sticky",
				".block-sidebar-rtb", ".adblock-service",
				"script", "noscript")

			writeNodes2File(fmt.Sprintf("chapters/%s.html", title.title), nodes)
		}(&wg, title)
	}
	wg.Wait()
}

func lookup4NewChapters(titles []Title) []Title {
	fileInfos, err := ioutil.ReadDir("chapters/")
	if err != nil {
		log.Fatal(err)
	}
	if len(fileInfos) > 0 {
		fileInfo := fileInfos[len(fileInfos)-1]
		for i := 0; i < len(titles); i++ {
			if fileInfo.Name() == titles[i].title+".html" {
				titles = titles[:i]
				break
			}
		}
	}
	return titles
}

func createDir(dirname string) {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		createDirError := os.Mkdir(dirname, 0755)
		if createDirError != nil {
			log.Fatal(createDirError)
		}
	}
}

func removeNodesList(nodes *goquery.Nodes, selectorsList ...string) {
	for _, selector := range selectorsList {
		removeNodes(selector, nodes)
	}
}

func removeNodes(selector string, nodes *goquery.Nodes) {
	foundNodes := nodes.Find(selector)
	for _, foundNode := range foundNodes {
		parentNode := foundNode.Parent
		parentNode.Remove(foundNode.Node)
	}
}

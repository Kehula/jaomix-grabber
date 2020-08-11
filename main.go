package main

import (
	"bufio"
	"fmt"
	"github.com/opesun/goquery"
	"log"
	"os"
)

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
		//*[@id="open-0"]
		foundNodes := nodes.Find(".flex-dow-txt")
		fmt.Println(foundNodes[0].Data)
		foundNodes = foundNodes.Find(".title")
		fmt.Println(foundNodes.Text())
		//resp, err := http.Get(requestUrl)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//defer resp.Body.Close()
		//ioutil.WriteFile("out.html", func (readerCloser io.ReadCloser) []byte{
		//	result, err := ioutil.ReadAll(readerCloser)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	return result
		//}(resp.Body), 0644)

	}


}

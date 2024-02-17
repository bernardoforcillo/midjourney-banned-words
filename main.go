package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type OutputFileStruct struct {
	Words map[string][]string `json:"words"`
}

func main() {
	url := "https://blog.easyprompt.xyz/the-complete-list-of-banned-words-in-midjourney-you-need-to-know-12111a5bbf87"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching URL:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Unexpected status code:", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	class := "md me fr mf b mg mh mi mj mk ml mm mn mo mp mq mr ms mt mu mv mw mx my mz na ow ox oy bj"

	elements := findElementsByClass(doc, class)

	cleanedElements := []string{}
	for _, element := range elements {
		re := regexp.MustCompile(`\([^)]*\)`)
		result := re.ReplaceAllString(element, "")
		result = strings.ToLower(result)
		cleanedElements = append(cleanedElements, result)
	}
	outputContent := OutputFileStruct{
		Words: map[string][]string{
			"en": cleanedElements,
		},
	}
	outputFile, err := os.Create("output.json")
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outputFile.Close()
	encoderOutput := json.NewEncoder(outputFile)
	encoderOutput.Encode(outputContent)
}

func findElementsByClass(n *html.Node, class string) []string {
	var elements []string
	if n.Type == html.ElementNode && n.Data == "li" && n.Parent.Data == "ul" {
		elements = append(elements, n.FirstChild.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		elements = append(elements, findElementsByClass(c, class)...)
	}
	return elements
}

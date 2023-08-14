package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// NOTE: Future perf update - cache component files so we don't do I/O every time.

var standardHtmlTags = map[string]bool{"a": true, "abbr": true, "address": true, "area": true, "article": true, "aside": true, "audio": true, "b": true, "base": true, "bdi": true, "bdo": true, "blockquote": true, "body": true, "br": true, "button": true, "canvas": true, "caption": true, "cite": true, "code": true, "col": true, "colgroup": true, "data": true, "datalist": true, "dd": true, "del": true, "details": true, "dfn": true, "dialog": true, "div": true, "dl": true, "dt": true, "em": true, "embed": true, "fieldset": true, "figcaption": true, "figure": true, "footer": true, "form": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true, "head": true, "header": true, "hr": true, "html": true, "i": true, "iframe": true, "img": true, "input": true, "ins": true, "kbd": true, "label": true, "legend": true, "li": true, "link": true, "main": true, "map": true, "mark": true, "meta": true, "meter": true, "nav": true, "noscript": true, "object": true, "ol": true, "optgroup": true, "option": true, "output": true, "p": true, "param": true, "picture": true, "pre": true, "progress": true, "q": true, "rp": true, "rt": true, "ruby": true, "s": true, "samp": true, "script": true, "section": true, "select": true, "small": true, "source": true, "span": true, "strong": true, "style": true, "sub": true, "summary": true, "sup": true, "svg": true, "table": true, "tbody": true, "td": true, "template": true, "textarea": true, "tfoot": true, "th": true, "thead": true, "time": true, "title": true, "tr": true, "track": true, "u": true, "ul": true, "var": true, "video": true, "wbr": true}

func compileHtml(inputHtml string) string {
	// Modify data
	r := regexp.MustCompile("<([\\w-]+)([^>]*)/>")
	modifiedData := r.ReplaceAllString(inputHtml, "<$1$2></$1>")

	// Parse the HTML
	reader := bytes.NewReader([]byte(modifiedData))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Println("Error parsing HTML:", true, err)
		return ""
	}

	// Iterate through elements and modify as needed
	doc.Find("*").Each(func(index int, element *goquery.Selection) {
		tagName := goquery.NodeName(element)
		if !standardHtmlTags[tagName] {
			fmt.Println("found custom tag", tagName)

			// Read the custom component's HTML file
			componentPath := fmt.Sprintf("./components/%s.html", tagName)
			componentHtml, err := os.ReadFile(componentPath)
			if err != nil {
				fmt.Println("Error reading component file:", err)
				return
			}

			// Get the inner HTML, including any text content
			innerHTML, _ := element.Html()
			componentHtmlString := strings.Replace(string(componentHtml), "{children}", innerHTML, -1)

			// Perform replacements for each attribute
			element.Each(func(index int, item *goquery.Selection) {
				for _, attr := range item.Nodes[0].Attr {
					attributeName := attr.Key
					attributeValue := attr.Val
					componentHtmlString = strings.Replace(componentHtmlString, "{"+attributeName+"}", attributeValue, -1)
				}
			})

			// Parse the modified component's HTML as a fragment
			componentFragment, err := html.ParseFragment(strings.NewReader(componentHtmlString), element.Get(0).Parent)
			if err != nil {
				fmt.Println("Error parsing component HTML:", err)
				return
			}

			// Replace the custom tag with the component's HTML content
			for _, node := range componentFragment {
				element.Get(0).Parent.InsertBefore(node, element.Get(0))
			}
			element.Remove()
		}
	})

	// Serialize the modified HTML
	htmlStr, err := doc.Find("body").Html()
	if err != nil {
		fmt.Println("Error serializing HTML:", true, err)
		return ""
	}

	return htmlStr
}

func compileHtmlFromFile(filePath string) string {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", true, filePath, err)
		return ""
	}

	compiledHtml := compileHtml(string(fileData))
	return compiledHtml
}

package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/html"

	"unicode"
)

type Link struct {
	Href string
	Text string
}

func camelCase(input string) int {
	count := 0
	if len(input) > 0 {
		count = 1
		for _, character := range input {
			if unicode.IsUpper(character) {
				count++
			}
		}
	}
	return count
}

func checkIfErrorOccurred(fileError error) {
	if fileError != nil {
		panic(fileError)
	}
}

func fileReader() {
	name := []byte("This is a go training test program where we will learn file handling.")
	fileName := "sample.txt"
	writeError := ioutil.WriteFile(fileName, name, fs.FileMode(fs.ModePerm))
	checkIfErrorOccurred(writeError)
	data, readError := ioutil.ReadFile(fileName)
	checkIfErrorOccurred(readError)
	fmt.Println("File Data for first sample is:\n", string(data))

	file, fileError := os.OpenFile("sample2.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	checkIfErrorOccurred(fileError)
	fileWriteError := os.WriteFile(file.Name(), name, fs.ModePerm.Perm())
	checkIfErrorOccurred(fileWriteError)
	fileData, fileReadError := os.ReadFile(file.Name())
	checkIfErrorOccurred(fileReadError)
	fmt.Println("File Data for second sample is:\n", string(fileData))
	file.Close()
}

func htmlLinkTokenizer(fileName string) []Link {
	links := []Link{}

	file, fileError := os.OpenFile(fileName, os.O_RDONLY, 0644)
	checkIfErrorOccurred(fileError)
	defer file.Close()
	fileData, fileReadError := os.ReadFile(file.Name())
	checkIfErrorOccurred(fileReadError)

	reader := strings.NewReader(string(fileData))
	tokenizer := html.NewTokenizer(reader)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		} else if tokenType == html.StartTagToken {
			currentTag, _ := tokenizer.TagName()
			key, val, _ := tokenizer.TagAttr()
			if string(key) == "href" {
				data := ""
				tokenizer.Next()
				tag, _ := tokenizer.TagName()
				for string(tag) != string(currentTag) && tokenType != html.CommentToken {
					data += tokenizer.Token().Data
					tokenType = tokenizer.Next()
					tag, _ = tokenizer.TagName()
				}
				data = strings.TrimSpace(data)
				links = append(links, Link{Href: string(val), Text: data})
			}
		}
	}
	return links
}

func htmlLinkTokenizerDriver() {
	fileNames := []string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}
	for _, name := range fileNames {
		fmt.Printf("Links extracted from %s are:\n", name)
		fmt.Println(htmlLinkTokenizer(name))
	}
}

func main() {
	wordCountInput := "saveChangesInTheEditor"
	fmt.Printf("Word count for %s is %d \n", wordCountInput, camelCase(wordCountInput))
	fileReader()
	htmlLinkTokenizerDriver()
}

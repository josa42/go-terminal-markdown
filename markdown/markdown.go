package markdown

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	link "github.com/josa42/go-terminal-hyperlink"
	image "github.com/josa42/go-terminal-image"
)

var (
	// MaxImageWidth :
	MaxImageWidth = 500
	bold          = color.New(color.Bold).Add(color.FgHiWhite).SprintFunc()
	undeline      = color.New(color.Underline).SprintFunc()
	gray          = color.New(color.FgHiBlack).SprintFunc()
	codeLine      = color.New(color.BgHiBlack).Add(color.FgHiWhite).SprintFunc()
	headlineExp   = regexp.MustCompile(`(^|\n)(#{1,6})([^#][^\n]+)`)
	headlineH1Exp = regexp.MustCompile(`(?m)^(.+)\n=+$`)
	headlineH2Exp = regexp.MustCompile(`(?m)^(.+)\n-+$`)
	linksExp      = regexp.MustCompile(`\[([^\[]+)\]\(([^\)]+)\)`)
	imgExp        = regexp.MustCompile(`!\[([^\[]*)\]\(([^\)]+)\)`)
	boldExp       = regexp.MustCompile(`(\*\*|__)([^*]+)(\*\*|__)`)
	blockquoteExp = regexp.MustCompile(`\n\>(.*)?`)
	lineExp       = regexp.MustCompile(`\n(-{5,})`)
	emptyLineExp  = regexp.MustCompile(`\n\n\n+`)
	codeBlockExp  = regexp.MustCompile("```[a-z]*\\n([\\s\\S]*?)\\n```")
	tmpFiles      = []string{}
)

func parse(md string) string {

	for _, v := range headlineExp.FindAllStringSubmatch(md, -1) {
		level := len(v[2])
		text := strings.TrimSpace(v[3])

		search := strings.TrimSpace(v[0])
		replace := formatHeadline(level, text)

		md = strings.Replace(md, search, replace, -1)
	}

	for _, v := range headlineH1Exp.FindAllStringSubmatch(md, -1) {
		text := strings.TrimSpace(v[1])
		search := strings.TrimSpace(v[0])
		replace := formatHeadline(1, text)

		md = strings.Replace(md, search, replace, -1)
	}

	for _, v := range headlineH2Exp.FindAllStringSubmatch(md, -1) {
		text := strings.TrimSpace(v[1])
		search := strings.TrimSpace(v[0])
		replace := formatHeadline(2, text)

		md = strings.Replace(md, search, replace, -1)
	}

	for _, v := range imgExp.FindAllStringSubmatch(md, -1) {
		href := strings.TrimSpace(v[2])

		if strings.HasPrefix(href, "http") {
			href = downloadFile(href)
			tmpFiles = append(tmpFiles, href)
		}

		search := strings.TrimSpace(v[0])
		replace := image.CreateWithSize(href, image.Size{MaxWidth: MaxImageWidth})

		md = strings.Replace(md, search, replace, -1)
	}

	for _, v := range linksExp.FindAllStringSubmatch(md, -1) {
		text := v[1]
		href := strings.TrimSpace(v[2])

		search := strings.TrimSpace(v[0])
		replace := link.Create(text, href)

		md = strings.Replace(md, search, replace, -1)
	}

	for _, v := range boldExp.FindAllStringSubmatch(md, -1) {
		text := v[2]

		search := v[0]
		replace := bold(text)

		md = strings.Replace(md, search, replace, -1)
	}

	for _, v := range blockquoteExp.FindAllStringSubmatch(md, -1) {
		text := v[1]

		search := strings.TrimSpace(v[0])
		replace := fmt.Sprintf(" %s%s", gray("|"), text)

		md = strings.Replace(md, search, replace, 1)
	}

	for _, v := range lineExp.FindAllStringSubmatch(md, -1) {
		search := strings.TrimSpace(v[0])
		md = strings.Replace(md, search, gray(strings.Repeat("-", 80)), 1)
	}

	for _, v := range emptyLineExp.FindAllStringSubmatch(md, -1) {
		md = strings.Replace(md, v[0], "\n\n", -1)
	}

	for _, v := range codeBlockExp.FindAllStringSubmatch(md, -1) {

		code := strings.Split(v[1], "\n")

		// prefix := strings.Repeat(" ", 4)

		for idx, line := range code {
			line = strings.Replace(line, "\t", "  ", -1)
			if len(line) > 72 {
				line = string([]rune(line)[0:72])
			}
			line = "  " + line + strings.Repeat(" ", (73)-len(line))
			line = codeLine(line)

			no := strconv.Itoa(idx + 1)
			code[idx] = fmt.Sprintf("%4s %s", no, line)
		}

		emptyLine := strings.Repeat(" ", 5) + codeLine(strings.Repeat(" ", 75))

		md = strings.Replace(md, v[0], emptyLine+"\n"+strings.Join(code, "\n")+"\n"+emptyLine, -1)
	}

	return strings.TrimSpace(md)
}

func formatHeadline(level int, text string) string {
	switch level {
	case 1:
		text = bold(undeline(text))
	case 2:
		text = bold(text)
	case 3:
		text = bold(text)
	case 4:
		text = bold(text)
	case 5:
		text = bold(text)
	case 6:
		text = bold(text)
	}

	return fmt.Sprintf("\n%s %s\n", strings.Repeat("#", level), strings.TrimSpace(text))
}

func downloadFile(url string) string {
	output, err := ioutil.TempFile(os.TempDir(), "prefix")
	// defer os.Remove(output.Name())

	// output, err := os.Create(fileName)
	if err != nil {
		// fmt.Println("Error while creating", fileName, "-", err)
		return ""
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		// fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()

	_, errCopy := io.Copy(output, response.Body)
	if errCopy != nil {
		// fmt.Println("Error while downloading", url, "-", err)
		return ""
	}

	return output.Name()
}

// Print :
func Print(md string) {
	out := parse(md)
	fmt.Println(out)

	for _, file := range tmpFiles {
		os.Remove(file)
	}
}

// func main() {
// 	out := parse(`
// # Hallo Welt
// ##Foobar
// ## Foobar2
//
// Some text with **bold** parts.
// Some text with __bold__ parts.
//
// ####### Foobar3
//
// [Google](http://google.de)
//
// ![](/Users/josa/Desktop/download.png)
// ![](https://www.huement.com/web/wp-content/uploads/2013/10/logo-1.jpg)
// 	`)
//
// 	fmt.Println(out)
//
// 	for _, file := range tmpFiles {
// 		os.Remove(file)
// 	}
// }

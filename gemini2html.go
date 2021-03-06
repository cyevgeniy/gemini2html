package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defTemplate = `<html><head>
    <meta charset='utf-8'>
	<style>body{font-family:Open Sans,Arial;color:#454545;font-size:16px;margin:2em auto;max-width:800px;padding:1em;line-height:1.4;text-align:justify}</style>
    <title>Personal blog</title><meta name='viewport' content='width=device-width, initial-scale=1'>
    </head><body>{%content%}</body></html>`

	POST_TEMPLATE = "templates/post.txt"
	DEFAULT_TEMPLATE = "templates/default.txt"
)

func trim(str string) string {
	return strings.TrimSpace(str)
}

func isList(str string) bool {
	return len(trim(str)) > 2 && trim(str)[0:2] == "* "
}

func isLink(str string) bool {
	return len(trim(str)) > 3 && trim(str)[0:3] == "=> "
}

func isHeader(str string) bool {
	strCut := trim(str)
	idx := strings.Index(strCut, " ")

	res := true

	if idx != -1 {
		hd := strCut[0:idx]
		for i := 0; i < len(hd); i++ {
			if hd[i] != '#' {
				res = false
				break
			}
		}
	} else {
		res = false
	}

	return res
}

func toHeader(str string) string {
	strCut := trim(str)
	idx := strings.Index(strCut, " ")

	lvl := 1

	hd := strCut[0:idx]
	lvl = len(hd)

	return "<h" + strconv.Itoa(lvl) + ">" + strCut[idx+1:] + "</h" + strconv.Itoa(lvl) + ">"

}

func toHref(str string) string {
	cutStr := trim(str[2:])
	idx := strings.Index(cutStr, " ")

	var href, cpt string
	if idx == -1 {
		href = cutStr
		cpt = href
	} else {
		href = cutStr[:idx]
		cpt = cutStr[idx+1:]
	}

	dotIdx := strings.LastIndex(href, ".gmi")
	if dotIdx != -1 {
		href = href[:dotIdx+1] + "html"
	}

	res := "<p><a href=\"" + href + "\">" + cpt + "</a></p>"

	return res
}

func toLi(str string) string {
	return "<li>" + trim(str)[2:] + "</li>"
}

func isVerb(str string) bool {
	return len(trim(str)) >= 3 && trim(str)[0:3] == "```"
}

// Did not exist in gemini specs, just for my personal use
func isImage(str string) bool {
	return len(trim(str)) > 3 && trim(str)[0:3] == "=] "
}

func toImg(str string) string {
	cutStr := trim(str[2:])
	return "<p><img src=\"" + cutStr + "\"/></p>"
}

func parseFile(filename string, writer *bufio.Writer) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cant open input file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	verb, list := false, false

	for scanner.Scan() {
		r := scanner.Text()

		if isVerb(r) {
			verb = !verb

			if verb {
				writer.WriteString("<pre>")
			} else {
				writer.WriteString("</pre>")
			}
			continue
		}

		if verb {
			writer.WriteString(r + "\n")
			continue
		}

		if !list && isList(r) {
			writer.WriteString("<ul>" + toLi(r))
			list = true
			continue
		}

		if !isList(r) && list {
			writer.WriteString("</ul>")
			list = false
			continue
		}

		if isList(r) {
			writer.WriteString(toLi(r))
			continue
		}

		if isLink(r) {
			writer.WriteString(toHref(r))
			continue
		}

		if isImage(r) {
			writer.WriteString(toImg(r))
			continue
		}

		if isHeader(r) {
			writer.WriteString(toHeader(r))
			continue
		}

		if r != "" {
			writer.WriteString("<p>" + trim(r) + "</p>")
			continue
		}
	}

	if list {
		writer.WriteString("</ul>")
	} else {
		if verb {
			writer.WriteString("</pre>")
		}
	}
}

func isFileExists(fileName string) bool {

	_, err := os.Stat(fileName)

	return !os.IsNotExist(err)

}

func validateFileExistance(fileName string) {
	if !isFileExists(fileName) {
		log.Fatal("Required file does not exists: " + fileName)
	}
}

func getTemplateContent(template string) string {

	var tmplStr string

	bts, err := ioutil.ReadFile(template)
	if err != nil {
		fmt.Println("Can't read content of template file. Default template will be used")
		tmplStr = defTemplate
	} else {
		tmplStr = string(bts)
	}

	return tmplStr
}

func file2html(filenameIn, filenameOut, template string) {

	fout, err := os.Create(filenameOut)
	if err != nil {
		log.Fatal("Can't create output file " + filenameOut)
	}

	validateFileExistance(filenameIn)

	tmplStr := getTemplateContent(template)

	headerFooter := strings.Split(tmplStr, "{%content%}")

	if len(headerFooter) != 2 {
		log.Fatal("Fuck, can't parse template")
	}

	writer := bufio.NewWriter(fout)
	writer.WriteString(headerFooter[0])
	parseFile(filenameIn, writer)
	writer.WriteString(headerFooter[1])
	writer.Flush()

	err = fout.Close()
	if err != nil {
		log.Fatal("Can't close output file")
	}
}

func getFileName(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return filename
	}

	return filename[0:idx]
}

func getFileNames(dir string) []string {

	fileList := make([]string, 0)

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return fileList
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			fileList = append(fileList, dir+file.Name())
		}
	}

	return fileList
}

func copyFile(fSrc, fDst string) int64 {

	src, err := os.Open(fSrc)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	// Create new file
	dst, err := os.Create(fDst)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	bytesCnt, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	return bytesCnt
}

func prepareDirs() {
	if _, err := os.Stat("_site"); os.IsNotExist(err) {
		os.Mkdir("_site", 0777)
	} else {
		err := os.RemoveAll("_site")

		if err != nil {
			log.Panic("Can't remove _site directory")
		}

		os.Mkdir("_site", 0777)
	}

	os.Mkdir("_site/posts", 0777)
	os.Mkdir("_site/assets", 0777)
	assets := getFileNames("assets/")
	for i := range assets {
		base := filepath.Base(assets[i])
		fmt.Printf("Copying file %s into "+"_site/assets/"+base, assets[i])
		bc := copyFile(assets[i], "_site/assets/"+base)
		fmt.Printf(" Done. Copied %d bytes\n", bc)
	}
}

func convertFiles(filenames []string, template string) {

	for i := range filenames {
		baseName := filepath.Base(filenames[i])
		ext := filepath.Ext(baseName)
		if ext == ".gmi" || ext == ".gemini" {
			fullPath := "_site/" + getFileName(filenames[i]) + ".html"
			file2html(filenames[i], fullPath, template)
		}
	}
}

func generateSite() {

	prepareDirs()

	postsFiles := getFileNames("./posts/")
	convertFiles(postsFiles, POST_TEMPLATE)

	otherFiles := getFileNames("./")
	convertFiles(otherFiles, DEFAULT_TEMPLATE)

}

func main() {
	generateSite()
}


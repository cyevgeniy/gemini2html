package main

import( "bufio"
	"os"
	"log"
	"strings"
)


func trim(str string) string {
	return strings.TrimSpace(str)
}

func isList(str string) bool {
	return len(trim(str)) > 2 && trim(str)[0:2] == "* "
}

func isLink(str string) bool {
	return len(trim(str)) > 3 &&  trim(str)[0:3] == "=> "
}

func toHref(str string) string {
	cutStr := trim(str[2:])
	idx := strings.Index(cutStr, " ")

	var href, cpt string
	if idx == -1 {
		href = cutStr
	} else {
		href = cutStr[:idx]
		cpt = cutStr[idx +1:]
	}

	res := "<a href=\"" + href  + "\">" + cpt + "</a>"

	return res
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

		if r == "```" {
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
		} else {
			if !list &&  isList(r) {
				writer.WriteString("<ul>")
				list = true
			} else {
				if !isList(r)  &&  list {
					writer.WriteString("</ul>")
					list = false
				}
			}

			if isList(r) {
				writer.WriteString("<li>" + trim(r)[2:] + "</li>")
			} else {
				if isLink(r) {
					writer.WriteString(toHref(r))
				} else {
					if r != "" {
						writer.WriteString("<p>" + trim(r) + "</p>")
					}
				}
			}
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

func main() {
	writer := bufio.NewWriter(os.Stdout)
	parseFile("test.gemini", writer)
	writer.Flush()
}

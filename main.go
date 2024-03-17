package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	var replaceExpr string
	flag.StringVar(&replaceExpr, "replace", "", `Replacement expression, e.g., "@search@replace@"`)
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: toolname --replace \"@search@replace@\" filename")
		os.Exit(1)
	}
	filePath := flag.Arg(0)

	if len(replaceExpr) < 3 {
		fmt.Println("Invalid replace expression format. Use \"@search@replace@\"")
		os.Exit(1)
	}

	delimiter := string(replaceExpr[0])
	parts := regexp.MustCompile(regexp.QuoteMeta(delimiter)).Split(replaceExpr[1:], -1)
	if len(parts) < 2 {
		fmt.Println("Invalid replace expression format. Use \"@search@replace@\"")
		os.Exit(1)
	}
	searchPattern, replacement := parts[0], parts[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re, err := regexp.Compile(searchPattern)
	if err != nil {
		fmt.Printf("Invalid regex pattern: %s\n", err)
		os.Exit(1)
	}

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := re.ReplaceAllString(line, replacement)
		fmt.Println(modifiedLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
}

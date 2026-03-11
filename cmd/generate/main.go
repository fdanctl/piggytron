package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	folderPath := "web/static/"

	jsFile, err := os.Create(folderPath + "app.js")
	if err != nil {
		panic(err)
	}
	defer jsFile.Close()

	cssFile, err := os.Create(folderPath + "styles.css")
	if err != nil {
		panic(err)
	}
	defer cssFile.Close()

	fmt.Println("Generating app.js...")
	pipefiles(folderPath+"js/", jsFile)
	fmt.Println("Generating styles.css...")
	pipefiles(folderPath+"css/", cssFile)
}

func pipefiles(dirPath string, dst *os.File) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, v := range dir {
		if v.IsDir() {
			pipefiles(dirPath+v.Name()+"/", dst)
		} else {
			src, err := os.Open(dirPath + v.Name())
			if err != nil {
				panic(err)
			}
			defer src.Close()

			scanner := bufio.NewScanner(src)
			for scanner.Scan() {
				line := scanner.Bytes()
				line = append(line, '\n')
				io.Writer.Write(dst, line)
			}
			if err := scanner.Err(); err != nil {
				panic(err)
			}

		}
	}
}

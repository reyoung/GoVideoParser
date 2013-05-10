package main

import (
	"flag"
	"fmt"
	"github.com/reyoung/GoVideoParser"
	"log"
)

func main() {
	Definition := 0
	flag.IntVar(&Definition, "definition", 0, "0 = Normal, 1=High, 2=Super")
	flag.Parse()
	youku := GoVideoParser.YoukuParser{}
	for _, durl := range flag.Args() {
		result, err := youku.Parse(durl, GoVideoParser.DefinitionType(Definition))
		if err != nil {
			log.Fatal(err)
		}
		for i, rurl := range result.GetURLS() {
			saveFileName := fmt.Sprintf("%s%d.%s", result.GetTitle(), i+1, result.GetFileType())
			fmt.Printf("wget '%s' -U 'Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)' -O '%s'\n", rurl, saveFileName)
		}
	}
}

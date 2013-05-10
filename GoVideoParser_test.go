package GoVideoParser

import (
	"log"
	"testing"
)

func TestYoukuGetType(*testing.T) {
	p := new(YoukuParser)
	if p.GetType() != PT_YOUKU {
		log.Fatal("Failed in Youku GetType ", p.GetType())
	}
}

func TestYoukuGetVideoID(*testing.T) {
	p := new(YoukuParser)
	url := "http://v.youku.com/v_show/id_XNTUzMTAxNzE2.html"
	id, err := p.getVideoID(url)
	if id != "XNTUzMTAxNzE2" || err != nil {
		log.Fatal("Failed in Youku GetVideoID where url:", url, " id:", id)
	}
}

func TestYoukuParse(t *testing.T) {
	p := YoukuParser{}
	result, err := p.Parse("http://v.youku.com/v_show/id_XNTUzMTAxNzE2.html", DT_NORMAL)
	if err != nil {
		log.Fatal(err)
	}
	if result.GetTitle() != "2013.5.4 相声《金刚腿》 刘春山、许健、王玥波" {
		log.Fatal("Title Parse Error")
	}

	if result.GetFileType() != "flv" {
		log.Fatal("File Type Parse Error")
	}

	for _, url := range result.GetURLS() {
		t.Log(url)
	}

	result, err = p.Parse("http://v.youku.com/v_show/id_XNTUzMTAxNzE2.html", DT_HIGH)

	if err != nil {
		log.Fatal(err)
	}

	if result.GetTitle() != "2013.5.4 相声《金刚腿》 刘春山、许健、王玥波" {
		log.Fatal("Title Parse Error of High Definition")
	}

	for _, url := range result.GetURLS() {
		t.Log(url)
	}
}

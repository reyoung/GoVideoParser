package GoVideoParser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type DefinitionType uint8

const (
	DT_NORMAL = (DefinitionType)(0)
	DT_HIGH   = (DefinitionType)(1)
	DT_SUPER  = (DefinitionType)(2)
)

type ParserType uint8

const (
	PT_YOUKU = (ParserType)(0)
)

type IParserError interface {
	Error() string
}

type VideoParseResult interface {
	GetURLS() []string
	GetTitle() string
	GetFileType() string
}

type YoukuVideoParseResult struct {
	URLS     []string
	Title    string
	FileType string
}

func (r YoukuVideoParseResult) GetURLS() []string {
	return r.URLS
}

func (r YoukuVideoParseResult) GetTitle() string {
	return r.Title
}

func (r YoukuVideoParseResult) GetFileType() string {
	return r.FileType
}

type IParser interface {
	Parse(url string, defi DefinitionType) (VideoParseResult, error)
	GetType() ParserType
}

type YoukuParser struct {
}

type YoukuError struct {
	What string
}

func (e YoukuError) Error() string {
	return e.What
}

type GeneralVideoParserError struct {
	What string
}

func (e GeneralVideoParserError) Error() string {
	return e.What
}
func (YoukuParser) GetType() ParserType {
	return PT_YOUKU
}

func (p YoukuParser) Parse(url string, defi DefinitionType) (VideoParseResult, error) {
	video_id, err := p.getVideoID(url)
	if err != nil {
		return nil, err
	}
	video_url := fmt.Sprintf("http://v.youku.com/v_show/id_%s.html", video_id)
	//ofile, _ := os.Create("D:\\2013\\GoVideo\\Output.log")
	doc, err := goquery.NewDocument(video_url)
	if err != nil {
		return nil, err
	}
	title := p.getTitleByGoQuery(doc)

	json_res, err := http.Get(fmt.Sprint("http://v.youku.com/player/getPlayList/VideoIDS/", video_id))
	if err != nil {
		return nil, err
	}

	json_body := json_res.Body

	defer json_body.Close()

	dec := json.NewDecoder(json_body)

	var m interface{}
	if err = dec.Decode(&m); err == io.EOF {
		// Do nothing
	} else if err != nil {
		return nil, err
	}
	data := m.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})
	segs := data["segs"].(map[string]interface{})
	// fmt.Fprint(ofile, "OK to convert segs", segs)

	var key string = ""
	if defi == DT_NORMAL {
		key = "flv"
	} else if defi == DT_HIGH {
		_, ok := segs["mp4"]
		if ok {
			key = `mp4`
		} else {
			key = `flv`
		}
	} else {
		_, hd2flag := segs["hd2"]
		if hd2flag {
			key = `hd2`
		} else {
			_, mp4flag := segs["mp4"]
			if mp4flag {
				key = `mp4`
			} else {
				key = `flv`
			}
		}
	}

	seed := int(data["seed"].(float64))

	mixed := bytes.Buffer{}
	source := []byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/\:._-1234567890`)
	for len(source) != 0 {
		seed = (seed*211 + 30031) & 0xFFFF
		index := seed * len(source) >> 16
		c := source[index]
		source = append(source[:index], source[index+1:]...)
		mixed.WriteByte(c)
	}

	mixed_array := mixed.Bytes()
	ids := strings.Split(data[`streamfileids`].(map[string]interface{})[key].(string), `*`)

	vidSb := bytes.Buffer{}

	for _, i := range ids {
		if len(i) != 0 {
			real_i, _ := strconv.Atoi(i)
			vidSb.WriteByte(mixed_array[real_i])
		}
	}
	vidSbArray := vidSb.Bytes()
	vidlow := vidSbArray[0:8]
	vidhigh := vidSbArray[10:]
	var file_type string = ""

	switch key {
	case "flv":
		file_type = "flv"
	case "hd2":
		file_type = "flv"
	case "mp4":
		file_type = "mp4"
	}

	segArray := segs[key].([]interface{})

	return_list := make([]string, 0, 10)

	for _, s := range segArray {
		sm := s.(map[string]interface{})
		nonum, _ := strconv.Atoi(sm["no"].(string))
		no := fmt.Sprintf("%02X", nonum)
		url := fmt.Sprintf("http://f.youku.com/player/getFlvPath/sid/00_%s/st/%s/fileid/%s%s%s?K=%s", no, file_type, vidlow, no, vidhigh, sm["k"].(string))
		return_list = append(return_list, url)
	}
	result := YoukuVideoParseResult{
		return_list, title, file_type}
	retv := VideoParseResult(result)
	return retv, nil
}

func (YoukuParser) getTitleByGoQuery(doc *goquery.Document) string {
	title := doc.Find("head title").Text()
	return strings.Replace(title, "—在线播放—优酷网，视频高清在线观看", "", -1)
}

func (YoukuParser) getVideoID(url string) (string, error) {
	reg, _ := regexp.Compile(`http://(v|www)\.youku\.com/v_show/id_(?P<id>[A-Za-z0-9]+)(|_rss)\.html`)
	if submatches := reg.FindStringSubmatch(url); submatches != nil {
		return submatches[2], nil
	}
	reg, _ = regexp.Compile(`http://(v|www)\.youku\.com/v_show/id_([A-Za-z0-9]+)(|_rss)\.html`)

	if submatches := reg.FindStringSubmatch(url); submatches != nil {
		return submatches[2], nil
	}
	reg, _ = regexp.Compile(`^loader\.swf\?VideoIDS=([A-Za-z0-9]+)`)

	if submatches := reg.FindStringSubmatch(url); submatches != nil {
		return submatches[1], nil
	}

	reg, _ = regexp.Compile(`^([A-Za-z0-9]+)$`)
	if submatches := reg.FindStringSubmatch(url); submatches != nil {
		return submatches[1], nil
	}
	return "", YoukuError{"Youku Parser don't Support such kind url"}
}

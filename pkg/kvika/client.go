package kvika

import (
	"fmt"
	"strings"

	curl "github.com/andelf/go-curl"
)

type Kvika struct {
}

type Response struct {
	StatusCode int
	Events     Events
}

func New() *Kvika {
	k := &Kvika{}
	return k
}

func (k *Kvika) Perform(req *Request, callback func(r *Recorder, buf []byte)) (*Response, error) {
	var err error
	easy := curl.EasyInit()
	defer easy.Cleanup()

	err = easy.Setopt(curl.OPT_URL, req.URL.String())

	{ // set headers
		allHeaders := make([]string, 0)
		for name, headers := range req.Headers {
			for _, header := range headers {
				allHeaders = append(allHeaders, fmt.Sprintf("%vs: %s", name, header))
			}
		}
		err = easy.Setopt(curl.OPT_HTTPHEADER, allHeaders)
		if err != nil {
			return nil, err
		}
	}
	switch method := strings.ToUpper(req.Method); method {
	case "GET":
		break
	case "POST":
		err = easy.Setopt(curl.OPT_POST, true)
		if err != nil {
			return nil, err
		}
	case "PUT":
		err = easy.Setopt(curl.OPT_PUT, true)
		if err != nil {
			return nil, err
		}
	default:
		err = easy.Setopt(curl.OPT_CUSTOMREQUEST, method)
		if err != nil {
			return nil, err
		}
	}
	// set payload
	if len(req.Payload) > 0 {
		err = easy.Setopt(curl.OPT_POSTFIELDSIZE, len(req.Payload))
		if err != nil {
			return nil, err
		}
		err = easy.Setopt(curl.OPT_READDATA, req.Payload)
		if err != nil {
			return nil, err
		}
	}
	r := newRecorder()
	// make a callback function
	err = easy.Setopt(curl.OPT_WRITEFUNCTION, func(buf []byte, userdata interface{}) bool {
		callback(r, buf)
		return true
	})
	if err != nil {
		return nil, err
	}
	r.Start()
	err = easy.Perform()
	if err != nil {
		return nil, err
	}
	err = recordCurlInfo(r, easy)
	if err != nil {
		return nil, err
	}
	respCode, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		return nil, err
	}
	resp := &Response{
		StatusCode: respCode.(int),
		Events:     r.sortedEvents(),
	}
	return resp, nil
}

const (
	NameLookupTime    = "NAMELOOKUP_TIME"
	ConnectTime       = "CONNECT_TIME"
	AppConnectTime    = "APPCONNECT_TIME"
	PreTransferTime   = "PRETRANSFER_TIME"
	StartTransferTime = "STARTTRANSFER_TIME"
	TotalTime         = "TOTAL_TIME"
)

func recordCurlInfo(r *Recorder, easy *curl.CURL) error {
	type Entry struct {
		Info curl.CurlInfo
		Name string
	}
	// https://curl.haxx.se/libcurl/c/curl_easy_getinfo.html
	entries := []Entry{
		{
			Info: curl.INFO_NAMELOOKUP_TIME,
			Name: NameLookupTime,
		},
		{
			Info: curl.INFO_CONNECT_TIME,
			Name: ConnectTime,
		},
		{
			Info: curl.INFO_APPCONNECT_TIME,
			Name: AppConnectTime,
		},
		{
			Info: curl.INFO_PRETRANSFER_TIME,
			Name: PreTransferTime,
		},
		{
			Info: curl.INFO_STARTTRANSFER_TIME,
			Name: StartTransferTime,
		},
		{
			Info: curl.INFO_TOTAL_TIME,
			Name: TotalTime,
		},
	}
	for _, entry := range entries {
		t, err := easy.Getinfo(entry.Info)
		if err != nil {
			return err
		}
		r.recordRaw(t.(float64)*1000.0, entry.Name, nil)
	}
	return nil
}

package network

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/mohamidsaiid/uniclipboard/internal/jsonParser"
)

type Servers struct {
	Wg  *sync.WaitGroup
	BaseURL string
	Port string
	ValidServer chan url.URL
	FinishedReqs chan bool
}

func (s *Servers)RequestServer(url url.URL) (*http.Response, error){
	s.Wg.Add(1)
	defer s.Wg.Done()
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Servers) ValidateServer(resp *http.Response) (bool, error) {
	var dst struct {
		Ok bool `json:"ok"`
	}	
	err := jsonParser.ReadJSON(resp, &dst)
	if err != nil {
		return false, err
	}
	return dst.Ok, nil
}
package discovery

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/mohamidsaiid/uniclipboard/internal/discovery/network"
)

func ValidServer(baseURL string, port string, path string, start, end int) (url.URL, error) {
	srvrs := &network.Servers{
		Wg:           &sync.WaitGroup{},
		BaseURL:      baseURL,
		Port:         port,
		ValidServer:  make(chan url.URL),
		FinishedReqs: make(chan bool),
	}

	loopThroughServers(path, start, end, srvrs)

	go func() {
		srvrs.Wg.Wait()
		srvrs.FinishedReqs <- true
	}()

	select {
	case serverURL := <-srvrs.ValidServer:
		return serverURL, nil
	case <-srvrs.FinishedReqs:
		return url.URL{}, errors.New("no working server was found after finishing reqs")
	case <-time.After(time.Minute * 2):
		return url.URL{}, errors.New("no working server was found after 2min")
	}
}

func loopThroughServers(path string, i, end int, s *network.Servers) {
	for i <= end {
		ip := fmt.Sprintf("%s.%d%s", s.BaseURL, i, s.Port)
		URL := url.URL{Scheme: "http", Host: ip, Path: path}

		go func() {
			resp, err := s.RequestServer(URL)
			if err != nil {
				return
			}
			ok, err := s.ValidateServer(resp)
			if err != nil {
				return
			}
			if ok {
				s.ValidServer <- URL
			}
		}()
		i++
	}
}

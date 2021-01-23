package network

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
)

//Curl make a http request
func Curl(url string, addheader []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body []byte
	spaceClient := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", RandomAgent())
	if addheader != nil {
		request.Header.Set(addheader[0], addheader[1])
	}

	response, err := spaceClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"Status": response.StatusCode,
			"Reason": response.Status,
			"URL":    url,
		}).Error("Status code not daijobu")
		return nil, errors.New(response.Status)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err

	}

	return body, nil
}

//CoolerCurl make a cooler http request *with multitor*
func CoolerCurl(urls string, addheader []string) ([]byte, error) {
	counter := 0
	for {
		counter++
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		proxyURL, err := url.Parse(config.BotConf.MultiTOR)
		if err != nil || counter == 3 {
			return nil, err
		}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				DialContext: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).DialContext,
			},
		}

		request, err := http.NewRequest("GET", urls, nil)
		if err != nil || counter == 3 {
			return nil, err
		}
		request.Header.Set("cache-control", "no-cache")
		request.Header.Set("User-Agent", RandomAgent())
		if addheader != nil {
			request.Header.Set(addheader[0], addheader[1])
		}
		response, err := client.Do(request.WithContext(ctx))
		if err != nil || counter == 3 {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK && counter == 3 {
			return nil, errors.New("Multi Tor get Error")
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil || counter == 3 {
			return nil, err
		}

		if data != nil {
			return data, nil
		}
	}
}

//RandomAgent Create random useragent for bypass some IDS
func RandomAgent() string {
	Agent := []string{"Windows / Firefox 77 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:77.0) Gecko/20100101 Firefox/77.0",
		"Linux / Firefox 77 [Desktop]: Mozilla/5.0 (X11; Linux x86_64; rv:77.0) Gecko/20100101 Firefox/77.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"}
	return Agent[rand.Intn(len(Agent))]
}

//CurlPost Make post request
func CurlPost(url string, payload []byte) error {
	body := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", RandomAgent())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

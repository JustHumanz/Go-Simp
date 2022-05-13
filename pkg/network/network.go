package network

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"google.golang.org/grpc"
)

//Curl make a http request
func Curl(url string, addheader map[string]string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body []byte
	spaceClient := http.Client{Transport: &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
	}}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", RandomAgent())
	if len(addheader) > 0 {
		for k, v := range addheader {
			request.Header.Set(k, v)
		}
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
		return nil, errors.New(strconv.Itoa(response.StatusCode))
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err

	}

	return body, nil
}

//CoolerCurl make a cooler http request *with multitor*
func CoolerCurl(urls string, addheader map[string]string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dialSocksProxy, err := proxy.SOCKS5("tcp", config.GoSimpConf.MultiTOR, nil, proxy.Direct)
	if err != nil {
		log.Error(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialSocksProxy.Dial,
		},
	}

	request, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", RandomAgent())
	if len(addheader) > 0 {
		for k, v := range addheader {
			request.Header.Set(k, v)
		}
	}
	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("multi Tor get Error " + response.Status)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
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

func InitNet() net.Listener {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return listen
}

func InitgRPC(server string) *grpc.ClientConn {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	return conn
}

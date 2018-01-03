package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

// FromIP makes a best effort to compute the request client IP.
func FromIP(req *http.Request) string {
	if f := req.Header.Get("X-Forwarded-For"); f != "" {
		return f
	}
	f := req.RemoteAddr
	ip, _, err := net.SplitHostPort(f)
	if err != nil {
		return f
	}
	return ip
}

// IsFromLocalIP check req is from local or not
func IsFromLocalIP(req *http.Request) bool {
	localIP := "::1"

	fromIP := FromIP(req)

	if localIP == fromIP {
		return true
	}

	return false
}

// PostBuffer2URL do URL post
func PostBuffer2URL(buffer []byte, url string, contentFormat string) ([]byte, error) {
	var responseBody []byte
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buffer))
	if err != nil {
		return responseBody, err
	}
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", contentFormat)

	client := &http.Client{}
	timeout := time.Duration(210 * time.Second) // set HTTP post timeout = 60s
	client.Timeout = timeout
	resp, err := client.Do(req)

	if err != nil {
		return responseBody, err
	}
	defer resp.Body.Close()
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

// GetBuffer2URL do URL get
func GetBuffer2URL(buffer []byte, url string, contentFormat string) ([]byte, error) {
	var responseBody []byte
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(buffer))
	if err != nil {
		return responseBody, err
	}
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", contentFormat)

	client := &http.Client{}
	timeout := time.Duration(210 * time.Second) // set HTTP post timeout = 60s
	client.Timeout = timeout
	resp, err := client.Do(req)

	if err != nil {
		return responseBody, err
	}
	defer resp.Body.Close()
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

// SendURI2URL post URI as ?a=12&b=44 to URL
// return: result byte, uri combined result, connect or post error
func SendURI2URL(uris map[string]string, url, method string) ([]byte, string, error) {
	var responseBody []byte
	var uriStr string
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return responseBody, uriStr, err
	}

	query := req.URL.Query()

	for key, value := range uris {
		query.Add(key, value)
	}

	req.URL.RawQuery = query.Encode()
	uriStr = req.URL.String()

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	timeout := time.Duration(210 * time.Second) // set HTTP post timeout = 210s
	client.Timeout = timeout
	resp, err := client.Do(req)
	if err != nil {
		return responseBody, uriStr, err
	}
	defer resp.Body.Close()

	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseBody, uriStr, err
	}

	return responseBody, uriStr, nil
}

// IPIsLocal check ip address is local ip or not
func IPIsLocal(remoteIP, localIP string) bool {
	if remoteIP == localIP || remoteIP == "::1" || remoteIP == "localhost" {
		return true
	}
	return false
}

// HTTPDownload Download file via HTTP
func HTTPDownload(dlURL, localFileName string) (int64, error) {
	if PathExist(localFileName) {
		os.Remove(localFileName)
	}

	localOut, err := os.Create(localFileName)
	if err != nil {
		return 0, err
	}
	defer localOut.Close()

	resp, err := http.Get(dlURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(localOut, resp.Body)
	return n, err
}

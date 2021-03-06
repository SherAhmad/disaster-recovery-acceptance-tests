package runner

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"strings"

	"time"

	"fmt"
	"net/url"

	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetAppURL(appName string) string {
	appStats := string(RunCommandAndRetry("cf app "+appName, 5).Out.Contents())
	var appURL string
	for _, line := range strings.Split(appStats, "\n") {
		if strings.HasPrefix(line, "routes:") {
			s := strings.Split(line, " ")
			appURL = s[len(s)-1]
		}
	}

	Expect(appURL).NotTo(BeEmpty())
	return appURL
}

func Get(url string) *http.Response {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	response, err := client.Get("https://" + url)
	Expect(err).NotTo(HaveOccurred())
	return response
}

func Post(url string, contentType string, body io.Reader) *http.Response {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	response, err := client.Post("https://"+url, contentType, body)
	Expect(err).NotTo(HaveOccurred())
	return response
}

func StatusCode(rawURL string) func() (int, error) {
	parsedURL, err := url.Parse(rawURL)
	Expect(err).NotTo(HaveOccurred(), "error parsing api url")
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	return func() (int, error) {
		client := &http.Client{
			Timeout: time.Minute,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}}
		fmt.Fprintf(GinkgoWriter, "Trying to connect to api url: %s\n", parsedURL.String())
		resp, err := client.Get(parsedURL.String())
		Expect(err).NotTo(HaveOccurred(), "error connecting to api url")
		return resp.StatusCode, err
	}
}

func RandomStringNumber() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random string number")
	}

	return base64.RawURLEncoding.EncodeToString(b)
}

type InstanceStatusResponse struct {
	State string
}

type AppStatusResponse map[string]InstanceStatusResponse

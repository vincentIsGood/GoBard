package httpclient

import (
	"bardtest/utils"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const USER_AGENT string = "simple-client/1.0"

type SimpleHttpClient struct{
    client *http.Client
    session *cookiejar.Jar

    defaultHeaders map[string](string)
    urlToCookies map[string]([]*http.Cookie)
}

func New() *SimpleHttpClient{
    jar, err := cookiejar.New(nil)
    utils.PanicOnError(err)

    client := &http.Client{}
    client.Jar = jar

    defaultHeaders := make(map[string]string)
    defaultHeaders["User-Agent"] = USER_AGENT

    return &SimpleHttpClient{
        client: client,
        session: jar,
        defaultHeaders: defaultHeaders,
        urlToCookies: make(map[string][]*http.Cookie),
    }
}

func (httpClient *SimpleHttpClient) AddCookie(urlStr string, key string, value string){
    if httpClient.urlToCookies[urlStr] == nil{
        httpClient.urlToCookies[urlStr] = make([]*http.Cookie, 0)
    }
    httpClient.urlToCookies[urlStr] = append(httpClient.urlToCookies[urlStr], &http.Cookie{
        Name: key, Value: value,
    })
}

// add headers: https://www.golangprograms.com/how-do-you-set-headers-in-an-http-request-with-an-http-client-in-go.html
func (httpClient *SimpleHttpClient) prepareCookies(){
    for urlStr, cookies := range httpClient.urlToCookies{
        urlObj, err := url.Parse(urlStr)
        utils.PanicOnError(err)

        httpClient.session.SetCookies(urlObj, cookies)
    }
}

func (httpClient *SimpleHttpClient) addDefaultHeadersFor(req *http.Request){
    httpClient.addHeadersFor(req, httpClient.defaultHeaders)
}
func (httpClient *SimpleHttpClient) addHeadersFor(req *http.Request, headers map[string]string){
    for key, value := range headers{
        req.Header.Add(key, value)
    }
}

func (httpClient *SimpleHttpClient) configureRequest(req *http.Request){
    httpClient.prepareCookies()
    httpClient.addDefaultHeadersFor(req)
}
func (httpClient *SimpleHttpClient) configureRequestWithHeaders(req *http.Request, headers map[string]string){
    httpClient.prepareCookies()
    httpClient.addDefaultHeadersFor(req)
    if headers != nil{
        httpClient.addHeadersFor(req, headers)
    }
}

func (httpClient *SimpleHttpClient) SendRequest(method string, url string) (*http.Response, error){
    req, err := http.NewRequest(method, url, nil)
    utils.PanicOnError(err)
    httpClient.configureRequest(req)

    return httpClient.client.Do(req)
}

func (httpClient *SimpleHttpClient) SendRequestWithHeaders(method string, url string, headers map[string]string) (*http.Response, error){
    req, err := http.NewRequest(method, url, nil)
    utils.PanicOnError(err)
    httpClient.configureRequestWithHeaders(req, headers)

    return httpClient.client.Do(req)
}

func (httpClient *SimpleHttpClient) SendRequestWithBody(method string, url string, bodyContent string) (*http.Response, error){
    req, err := http.NewRequest(method, url, strings.NewReader(bodyContent))
    utils.PanicOnError(err)
    httpClient.configureRequest(req)

    return httpClient.client.Do(req)
}

func (httpClient *SimpleHttpClient) SendRequestWithHeadersAndBody(method string, url string, headers map[string]string, bodyContent string) (*http.Response, error){
    req, err := http.NewRequest(method, url, strings.NewReader(bodyContent))
    utils.PanicOnError(err)
    httpClient.configureRequestWithHeaders(req, headers)

    return httpClient.client.Do(req)
}

func (httpClient *SimpleHttpClient) String() string{
    return fmt.Sprintf("SimpleHttpClient{defaultHeaders: %s, urlToCookies: %s}", 
        httpClient.defaultHeaders, httpClient.urlToCookies)
}

// Public utils
func EncodePostForm(data map[string]string) string{
    resultStr := make([]string, 0)
    for key, value := range data{
        resultStr = append(resultStr, key + "=" + url.QueryEscape(value))
    }
    return strings.Join(resultStr, "&")
}
package vcs

import (
	"os"
	"bufio"
	"io"
	"strings"
	"net/url"
	"net/http"
	//"fmt"
)
var DefaultDomainPath string = "/root/.proxy_domain_file"
var DefaultProxyServer string = "127.0.0.1:80"
var ProxyDomainMap = make(map[string]bool)
var HttpProxyServer string = ""
func Init(proxyDomainFile string, httpProxyServer string)  error{
	HttpProxyServer = httpProxyServer
	fd, err := os.Open(proxyDomainFile)
	if err != nil{
		return err
	}
	rd := bufio.NewReader(fd)
	for{
		line, err := rd.ReadString('\n')
		if err != nil{
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		domain := strings.Trim(line, " \n")
		if domain != "" {
			ProxyDomainMap[domain] = true
		}
	}
	return nil
}

func UseProxy(cvsUrl string) bool {
	u, err := url.Parse(cvsUrl)
	if err != nil{
		return false
	}
	host := u.Host
	_, ok := ProxyDomainMap[host]
	//fmt.Println(cvsUrl, " use proxy: ", ok)
	return ok
}

func httpGetWithProxy(httpUrl string)(*http.Response,error){
	if UseProxy(httpUrl) {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(HttpProxyServer)
		}
		transport := &http.Transport{Proxy: proxy}
		client := &http.Client{Transport: transport}
		//fmt.Println(httpUrl, "use proxy http Get")
		return client.Get(httpUrl)
	} else {
		return http.Get(httpUrl)
	}
}

func ProxyEnv() []string{
	//fmt.Println("use proxy Env")
	return []string {"http_proxy=" + HttpProxyServer, "https_proxy=" + HttpProxyServer}
}

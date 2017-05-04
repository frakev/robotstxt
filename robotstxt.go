package robotstxt

import (
	"errors"
	"net/http"
	"io/ioutil"
	"net/url"
	"time"
	"fmt"
	"bufio"
	"strings"
	"regexp"
	"runtime"

	"github.com/alexcesaro/log/stdlog"
)

var	logger = stdlog.GetFromFlags()

type Directive struct {
	allow []string
	disallow []string
}

func getHost(site string) (finalUrl string) {

	u, err := url.Parse(site)
	if err != nil {
		logger.Error(funcName(), err)
	}
	finalUrl = u.Scheme + "://" + u.Host + "/robots.txt"
	logger.Debug("Robots.txt: ",finalUrl)
	return finalUrl
}

func getRobotstxt(url string) (err error, robots string) {

	robots = ""
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	rbtUrl := getHost(url)

	req, err := http.NewRequest("GET", rbtUrl, nil)
	if err != nil {
		return err, robots
	}

	req.Header.Set("User-Agent", "Robot-Test")

	resp, err := client.Do(req)
	if err != nil {
		return err, robots
	}

	status := resp.StatusCode
	switch {
	case 400 <= status && status < 500:
		err = errors.New(resp.Status + ":" + rbtUrl)
		return err, robots
	case 500 <= status && status < 600:
		err = errors.New(resp.Status + ":" + rbtUrl)
		return err, robots
	case 200 <= status && status < 300:
		logger.Info("Status:",resp.Status, ":", rbtUrl)
	}

	rbt, err := ioutil.ReadAll(resp.Body)
	robots = fmt.Sprintf("%s",rbt)
	resp.Body.Close()
	if err != nil {
		return err, robots
	}
	return nil, robots
}

func parseRobotstxt(url string) (res map[string]*Directive, err error) {

	err, rbts := getRobotstxt(url)
	if err != nil {
		logger.Error(funcName(), err)
		return res, err
	}
	scanner := bufio.NewScanner(strings.NewReader(rbts))
	scanner.Split(bufio.ScanWords)

	res = func() (r map[string]*Directive){
		var CurrentUserAgent string
		r = make(map[string]*Directive)
		for scanner.Scan() {
			switch scanner.Text() {
			case "User-agent:":
				scanner.Scan()
				CurrentUserAgent = scanner.Text()
				r[CurrentUserAgent] = &Directive{}
				logger.Debug("User-agent:",scanner.Text())
			case "Disallow:":
				scanner.Scan()
				logger.Debug("Disallow:",scanner.Text())
				r[CurrentUserAgent].disallow = append(r[CurrentUserAgent].disallow, scanner.Text())
			case "Allow:":
				scanner.Scan()
				logger.Debug("Allow:",scanner.Text())
				r[CurrentUserAgent].allow = append(r[CurrentUserAgent].allow, scanner.Text())
			default:
				break
			}
		}
		return r
	}()

	return res, err
}

func IsAllowed(url string) (ok bool, err error) {

	res, err := parseRobotstxt(url)

	if err != nil {
		return false, err
	}

	for user, path := range res {
		if exists, _ := regexp.MatchString("(?i)Robot-Test|\\*", user); !exists {
			continue
		}
		for _, pattern := range path.disallow {
			pattern := strings.Replace(pattern, "*", ".*", -1)
			if pattern == "" {
				continue
			}
	    	match, _ := regexp.MatchString(pattern, url)
	    	if match {
		        logger.Error(user ,url, pattern, match)
		        return false, err
		    }
	    }
	}
	return true, err
}

func funcName() string {
    pc, _, _, _ := runtime.Caller(1)
    return "[" + runtime.FuncForPC(pc).Name() + "]"
}
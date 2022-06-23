package url_args

import (
	"net/url"
	"strings"
)

func Compare(link string, link2 string, args []string) bool {
	if link == link2 {
		return true
	}

	url1, err1 := url.Parse(link)
	url2, err2 := url.Parse(link2)
	if err1 != nil || err2 != nil {
		return false
	}
	if url1.Host != url2.Host {
		return false
	}

	if url1.Path != url2.Path {
		return false
	}

	if args == nil || len(args) == 0 {
		return true
	}

	query1 := url1.Query()
	query2 := url2.Query()

	var urlArgs = make([]string, 0)
	for _, k := range args {
		if strings.HasPrefix(k, "url.") {
			urlArgs = append(urlArgs, strings.Replace(k, "url.", "", 1))
		} else {
			if query1.Get(k) != query2.Get(k) {
				return false
			}
		}
	}
	if len(urlArgs) > 0 {
		l1 := query1.Get("url")
		if l1 == "" {
			return false
		}
		l2 := query2.Get("url")
		if l2 == "" {
			return false
		}

		u1, err1 := url.QueryUnescape(l1)
		if err1 != nil {
			return false
		}
		u2, err1 := url.QueryUnescape(l2)
		if err1 != nil {
			return false
		}
		return Compare(u1, u2, urlArgs)
	}
	return true
}

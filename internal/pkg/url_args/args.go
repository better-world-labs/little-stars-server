package url_args

import "net/url"

func getQueryKeys(values url.Values) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		if k == "url" {
			link, err := url.QueryUnescape(values.Get(k))
			if err == nil {
				u, err := url.Parse(link)
				if nil == err {
					queryKeys := getQueryKeys(u.Query())
					for _, ke := range queryKeys {
						keys = append(keys, "url."+ke)
					}
				}
			}
		} else {
			keys = append(keys, k)
		}
	}
	return keys
}

func GetArgs(l string) []string {
	link, err := url.Parse(l)
	if err != nil {
		return []string{}
	}
	query := link.Query()
	return getQueryKeys(query)
}

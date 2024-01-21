package internal

import "net/url"

func MergeURLValues(dst url.Values, src url.Values, scope string) {
	for k, v := range src {
		if scope != "" {
			k = scope + "[" + k + "]"
		}

		for _, s := range v {
			dst.Add(k, s)
		}
	}
}

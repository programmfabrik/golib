package golib

import "net/url"

// MustBuildURL uses "u" as basis and adds given query values
// Query values must be given as key1, value1, [key2, value2[, ...]]
// Panics if the URL is not parsable
func MustBuildURL(u string, qv ...string) string {
	u1, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	u1.RawQuery = AddToQuery(u1.Query(), qv...).Encode()
	return u1.String()
}

// AddToQuery adds qv (key, value) to query.
// The map is returned for easy chaining
func AddToQuery(q url.Values, qv ...string) url.Values {
	for i := 0; i < len(qv)-1; i += 2 {
		k := qv[i]
		v := qv[i+1]
		if v == "" {
			continue
		}
		q.Set(k, v)
	}
	return q
}

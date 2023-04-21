package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if newUrl, found := pathsToUrls[r.URL.Path]; found {
			http.Redirect(w, r, newUrl, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

type redirectInfo struct {
	Path, Url string
}

func parseYAML(yml []byte) ([]redirectInfo, error) {
	var data []redirectInfo
	err := yaml.Unmarshal(yml, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func buildMap(redirectList []redirectInfo) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, info := range redirectList {
		pathsToUrls[info.Path] = info.Url
	}
	return pathsToUrls
}

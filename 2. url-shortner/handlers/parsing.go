package url_shortner

type pathURL struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func makePathsToURLMapFromSlice(paths []pathURL) map[string]string {
	pathsToUrl := make(map[string]string)
	for _, path := range paths {
		pathsToUrl[path.Path] = path.Url
	}
	return pathsToUrl
}

type unmarshaller = func(in []byte, out interface{}) (err error)
type parser = func([]byte) ([]pathURL, error)

func parserFromUnmarshaller(unmarshaller unmarshaller) parser {
	return func(bytes []byte) ([]pathURL, error) {
		var paths []pathURL
		err := unmarshaller(bytes, &paths)
		if err != nil {
			return nil, err
		}
		return paths, nil
	}
}

package url_shortner

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	"gopkg.in/yaml.v3"
)

const pathNotFound string = "Path not found!"

func handlerFromParser(
	parser parser,
	bytes []byte,
	fallback http.Handler,
) (http.HandlerFunc, error) {
	paths, err := parser(bytes)
	if err != nil {
		return nil, err
	}
	pathsToURL := makePathsToURLMapFromSlice(paths)
	return MapHandler(pathsToURL, fallback), nil
}

func handlerFromUnmarshaller(
	unmarshaller unmarshaller,
	bytes []byte,
	fallback http.Handler,
) (http.HandlerFunc, error) {
	return handlerFromParser(parserFromUnmarshaller(unmarshaller), bytes, fallback)
}

func redirector(
	getDestinationPath func(string) (string, error),
	fallback http.Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		dest, err := getDestinationPath(path)
		if err == nil {
			http.Redirect(w, r, dest, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func MapHandler(pathsToURL map[string]string, fallback http.Handler) http.HandlerFunc {
	return redirector(
		func(path string) (string, error) {
			if dest, ok := pathsToURL[path]; ok {
				return dest, nil
			} else {
				return "", errors.New(pathNotFound)
			}
		},
		fallback,
	)
}

func DBHandler(db *sql.DB, table string, pathCol string, destCol string, fallback http.Handler) http.HandlerFunc {
	preparedQuery := fmt.Sprintf("SELECT %s AS dest FROM %s WHERE %s = ?;", destCol, table, pathCol)
	return redirector(
		func(path string) (string, error) {
			stmt, err := db.Prepare(preparedQuery)
			if err != nil {
				return "", err
			}
			defer stmt.Close()
			var dest string
			queryRow := stmt.QueryRow(path)
			err = queryRow.Scan(&dest)
			if err != nil {
				return "", err
			}
			return dest, nil
		},
		fallback,
	)
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return handlerFromUnmarshaller(yaml.Unmarshal, yamlBytes, fallback)
}

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return handlerFromUnmarshaller(json.Unmarshal, jsonBytes, fallback)
}

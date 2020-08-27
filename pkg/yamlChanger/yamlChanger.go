package yamlChanger

import (
	"fmt"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
	"strings"
)

type PathError struct {
}

func (p *PathError) Error() string {
	return fmt.Sprintf("valuePath or part of it is empty")
}


func GetPathSplits(path string) (res []interface{}, err error) {
	splits := strings.Split(path, ".")
	if len(splits) == 0 {
		return nil, &PathError{}
	}
	for _, value := range splits {
		if value == "" {
			return nil, &PathError{}
		}
		res = append(res, value)
	}

	return res, nil

}

func ChangeYaml(body map[string]interface{}, newValue string, path []interface{}) (string, error) {

	err := dyno.Set(body, newValue, path...)
	if err != nil {
		return "", err
	}

	byteRes, err := yaml.Marshal(&body)
	return string(byteRes), err
}

package env

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	files  []string
	values map[string]string
}

func Load(files ...string) (*Env, error) {
	if len(files) == 0 {
		files = []string{".env"}
	}

	env := &Env{
		files:  files,
		values: make(map[string]string),
	}

	err := env.load()

	if err != nil {
		return nil, err
	}

	env.set()

	return env, nil
}

func (env *Env) load() error {

	for _, file := range env.files {
		err := env.readFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) readFile(file string) error {
	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	var lines []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			key, value, err := parseLine(fullLine)

			if err == nil {
				env.values[key] = value
			}
		}
	}
	return nil
}

func (env *Env) set() {
	for key, value := range env.values {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}

func parseLine(line string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	// now split key from value
	splitString := strings.SplitN(line, "=", 2)

	if len(splitString) != 2 {
		// try yaml mode!
		splitString = strings.SplitN(line, ":", 2)
	}

	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	// Parse the key
	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	// Parse the value
	value = splitString[1]
	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values
	if strings.Count(value, "\"") == 2 || strings.Count(value, "'") == 2 {
		// pull the quotes off the edges
		value = strings.Trim(value, "\"'")

		// expand quotes
		value = strings.Replace(value, "\\\"", "\"", -1)
		// expand newlines
		value = strings.Replace(value, "\\n", "\n", -1)
	}

	return
}

func (env *Env) Get(key string) string {
	vals := env.values

	return vals[key]
}

func (env *Env) Set(key, val string) {
	env.values[key] = val
}

func (env *Env) GetInt(key string) (int, error) {
	v := env.Get(key)
	return strconv.Atoi(v)
}

func (env *Env) GetDouble(key string) (float64, error) {
	v := env.Get(key)
	return strconv.ParseFloat(v, 64)
}

func (env *Env) GetBool(key string) (bool, error) {
	v := env.Get(key)
	return strconv.ParseBool(v)
}

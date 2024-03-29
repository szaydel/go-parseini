package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func NewDict() Dict {
	return Dict{}
}

func NewSection() Section {
	return Section{}
}

func MustLoadReader(reader *bufio.Reader) Dict {
	dict, err := LoadReader(reader)
	if err != nil {
		panic(err)
	}
	return dict
}

func MustLoad(filename string) Dict {
	dict, err := Load(filename)
	if err != nil {
		panic(err)
	}
	return dict
}

func LoadReader(reader *bufio.Reader) (dict Dict, err error) {
	dict = NewDict()
	lineno := 0
	section := EMPTY_STRING
	dict[section] = NewSection()

	for err == nil {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineno++
		if len(l) == 0 {
			continue
		}
		line := trim(string(l), unicode.IsSpace)

		for line[len(line)-1] == '\\' {
			line = line[:len(line)-1]
			l, _, err := reader.ReadLine()
			if err != nil {
				return nil, err
			}
			line += trim(string(l), unicode.IsSpace)
		}

		section, err = dict.parseLine(section, line)
		if err != nil {
			return nil, newError(
				err.Error() + fmt.Sprintf("':%d'.", lineno))
		}
	}

	return
}

func Load(filename string) (Dict, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return LoadReader(reader)
}

func Write(filename string, dict *Dict) error {
	buffer := dict.format()
	return ioutil.WriteFile(filename, buffer.Bytes(), PERMISSION)
}

func (e Error) Error() string {
	return string(e)
}
func (dict Dict) parseLine(section, line string) (string, error) {
	begin := line[0]
	end := line[len(line)-1]

	// comments should be ignored normally
	if begin == POUND || begin == SEMICOLON {
		return section, nil
	}

	// is section if ^[...somevalue...]$
	if begin == LEFT_BRKT && end == RIGHT_BRKT {
		section := trim(line[1:len(line)-1], unicode.IsSpace)
		section = tolower(section)
		dict[section] = NewSection()
		return section, nil
	}

	// key = value
	if m := regDoubleQuote.FindAllStringSubmatch(line, 1); m != nil {
		key, val := m[0][1], m[0][2]
		dict.add(section, key, val)
		return section, nil

	} else if m = regSingleQuote.FindAllStringSubmatch(line, 1); m != nil {
		key, val := m[0][1], m[0][2]
		dict.add(section, key, val)
		return section, nil

	} else if m = regNoQuote.FindAllStringSubmatch(line, 1); m != nil {
		key, val := m[0][1], m[0][2]
		dict.add(section, key, trim(val, unicode.IsSpace))
		return section, nil

	} else if m = regNoValue.FindAllStringSubmatch(line, 1); m != nil {
		key, val := m[0][1], EMPTY_STRING
		dict.add(section, key, val)
		return section, nil
	}

	return section, newError("iniparser: syntax error at ")
}

func (dict Dict) add(section, key, value string) {
	key = strings.ToLower(key)
	dict[section][key] = value
}

func (dict Dict) GetBool(section, key string) (bool, bool) {
	sec, ok := dict[section]
	if !ok {
		return false, false
	}
	value, ok := sec[key]
	if !ok {
		return false, false
	}
	v := value[0]

	switch {
	case v == 'y' || v == 'Y' || v == '1' || v == 't' || v == 'T':
		return true, true
	case v == 'n' || v == 'N' || v == '0' || v == 'f' || v == 'F':
		return false, true
	}

	return false, false
}

func (dict Dict) setValue(section, key, value string) {
	_, ok := dict[section]
	if !ok {
		dict[section] = NewSection()
	}
	dict[section][key] = value
}

func (dict Dict) SetBool(section, key string, value bool) {
	dict.setValue(section, key, fmtBool(value))
}

func (dict Dict) GetString(section, key string) (string, bool) {
	sec, ok := dict[section]
	if !ok {
		return EMPTY_STRING, false
	}
	value, ok := sec[key]
	if !ok {
		return EMPTY_STRING, false
	}
	return value, true
}

func (dict Dict) SetString(section, key, value string) {
	dict.setValue(section, key, value)
}

func (dict Dict) GetInt(section, key string) (int, bool) {
	sec, ok := dict[section]
	if !ok {
		return 0, false
	}
	value, ok := sec[key]
	if !ok {
		return 0, false
	}
	i, err := atoi(value)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (dict Dict) SetInt(section, key string, value int) {
	dict.SetString(section, key, fmtInt(int64(value), 10))
}

func (dict Dict) GetDouble(section, key string) (float64, bool) {
	sec, ok := dict[section]
	if !ok {
		return 0, false
	}
	value, ok := sec[key]
	if !ok {
		return 0, false
	}
	d, err := pFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return d, true
}

func (dict Dict) SetDouble(section, key string, value float64) {
	dict.SetString(section, key, fmtFloat(value, 'f', -1, 64))
}

func (dict Dict) Delete(section, key string) {
	_, ok := dict[section]
	if !ok {
		return
	}
	delete(dict[section], key)
	// If there are no items left in the section,
	// delete the section.
	if len(dict[section]) == 0 {
		delete(dict, section)
	}
}

func (dict Dict) GetSections() []string {
	size := len(dict)
	sections := make([]string, size)
	i := 0
	for section, _ := range dict {
		sections[i] = section
		i++
	}
	return sections
}

func (dict Dict) String() string {
	return (*dict.format()).String()
}

func (dict Dict) format() *bytes.Buffer {
	var buffer bytes.Buffer
	for section, vals := range dict {
		if len(section) > 0 {
			buffer.WriteString(fmt.Sprintf("[%s]\n", section))
		}
		for key, val := range vals {
			buffer.WriteString(fmt.Sprintf("%s = %s\n", key, val))
		}
		buffer.WriteString("\n")
	}
	return &buffer
}

func newError(message string) (e error) {
	return Error(message)
}

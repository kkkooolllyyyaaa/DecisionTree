package dataWorker

import (
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Students struct {
	Studs []Student
}

func (st Students) ColumnAsSlice(idx int) (result []string) {
	for _, value := range st.Studs {
		result = append(result, value.Attributes[idx])
	}
	return
}

func (st Students) Grades() (result []string) {
	for _, value := range st.Studs {
		result = append(result, value.Grade)
	}
	return
}

func (st Students) SelectWhereEq(idx int, target string) (result Students) {
	for _, value := range st.Studs {
		attrValue := value.Attributes[idx]
		if attrValue == target {
			result.Studs = append(result.Studs, value)
		}
	}
	return
}

var studentNotFound = errors.New("Student not found")

func (st Students) findById(id string) (found Student, err error) {
	for _, value := range st.Studs {
		if value.Id == id {
			return value, nil
		}
	}
	return Student{}, studentNotFound
}

type Student struct {
	Id         string
	Attributes []string
	Grade      string
}

type Header struct {
	AttributesNames []string
}

func Parse(filename, sep string) ([]Student, Header, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, Header{}, errors.Wrap(err, "reading file...")
	}
	data := string(rawData)
	lines := strings.Split(data, "\n")

	header, err := parseHeader(lines[0], sep)
	if err != nil {
		return nil, Header{}, errors.Wrap(err, "parsing header...")
	}

	studentsLines := lines[1:]
	students := make([]Student, 0, len(studentsLines))
	for _, line := range studentsLines {
		person, err := parseStudent(line, sep)
		if err != nil {
			return nil, Header{}, errors.Wrap(err, "parsing line...")
		}
		students = append(students, person)
	}
	return students, header, nil
}

func Write(filename, sep string, data [][]string, columns int) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Can't create file")
	}
	nl := "\n"
	toWrite := "actual" + sep + "predicted" + nl
	_, err = file.WriteString(toWrite)
	if err != nil {
		return errors.Wrap(err, "Can't write header")
	}

	for _, value := range data {
		if len(value) != columns {
			return errors.New("Length of data is not equal to columns size")
		}
		toWrite = value[0]
		for i := 1; i < len(value); i++ {
			toWrite += sep + value[i]
		}
		toWrite += nl
		_, err = file.WriteString(toWrite)
		if err != nil {
			return errors.Wrap(err, "Can't write line")
		}
	}
	return nil
}

func parseStudent(line string, sep string) (Student, error) {
	columns := strings.Split(line, sep)
	size := len(columns)
	return Student{Id: columns[0], Attributes: columns[1 : size-1], Grade: columns[size-1]}, nil
}

func parseHeader(line string, sep string) (Header, error) {
	headers := strings.Split(line, sep)
	return Header{AttributesNames: headers}, nil
}

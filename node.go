package main

import "bufio"
import "io"
import "math/rand"
import "os"
import "compress/gzip"
import "strings"
import "time"

import "github.com/pborman/uuid"

type Node struct {
	Id   string
	Name string
	Addr string
}

func NewRandomNode() (Node, error) {
	id := generateNewId()
	name, err := generateNewName()
	if err != nil {
		return Node{}, err
	}

	return Node{id, name, ""}, nil
}

func NewNode(id string, name string, addr string) Node {
	return Node{id, name, addr}
}

func generateNewId() string {
	newId := uuid.NewRandom().URN()
	return newId[strings.LastIndex(newId, ":")+1:]
}

func generateNewName() (string, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	nameIdx := r1.Intn(258000)
	surnameIdx := r1.Intn(88799)

	nameFile, err := os.Open("names.dat.gz")
	if err != nil {
		return "", err
	}

	nameReader, err := gzip.NewReader(nameFile)
	if err != nil {
		return "", err
	}
	defer nameReader.Close()

	name, err := readLine(nameReader, nameIdx)

	surnameFile, err := os.Open("names.dat.gz")
	if err != nil {
		return "", err
	}

	surnameReader, err := gzip.NewReader(surnameFile)
	if err != nil {
		return "", err
	}
	defer surnameReader.Close()

	surname, err := readLine(surnameReader, surnameIdx)

	return strings.Join([]string{name, surname}, " "), nil

}

func readLine(r io.Reader, line int) (string, error) {
	sc := bufio.NewScanner(r)
	lastLine := 0
	for sc.Scan() {
		lastLine++
		if lastLine == line {
			return sc.Text(), nil
		}
	}
	return "", io.EOF
}

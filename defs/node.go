package defs

import "bufio"
import "io"
import "math/rand"
import "os"
import "strings"
import "time"

import "github.com/pborman/uuid"

type NodeService interface {
	AddNode(Node) error
	NodeExists(string) (bool, error)
	RemoveNode(string)
	GetNode(string) (Node, error)
	NodeInfo() Node
}

type Node struct {
	Id   string
	Name string
	Addr string
}

func GenerateNewNode() (Node, error) {
	id := generateNewId()
	name, err := generateNewName()
	if err != nil {
		return Node{}, err
	}

	return Node{id, name, ""}, nil
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

	nameFile, err := os.Open("names.dat")
	if err != nil {
		return "", err
	}
	defer nameFile.Close()

	name, err := readLine(nameFile, nameIdx)

	surnameFile, err := os.Open("surnames.dat")
	if err != nil {
		return "", err
	}
	defer surnameFile.Close()

	surname, err := readLine(surnameFile, surnameIdx)

	return name + " " + surname, nil

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

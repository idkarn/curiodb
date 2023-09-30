package common

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var ResponseStrings = map[string]string{
	"T1": "Table with this id not found",
	"C2": "Column with this name was not found",
	"R0": "Row with id %d has been found",
	"R1": "Row with this id was not found",
	"R2": "Row with id %d has been deleted",
	"R3": "New row has been created with id %d",
}

func DecodeJson[T IDecodedJson](r *http.Request) (*T, error) {
	var parsed T
	if err := json.NewDecoder(r.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

func EncodeJson[T Row[string] | TableColumn | []Row[string]](obj T) (string, error) {
	out, err := json.Marshal(obj)
	return string(out), err
}

func CreateFile(name string) *os.File {
	f, err := os.Create(name)
	if err != nil {
		log.Println("Couldn't open file")
	}
	return f
}

func WriteFile(file io.Writer, binData any) error {
	enc := gob.NewEncoder(file)
	if err := enc.Encode(binData); err != nil {
		return err
	}
	return nil
}

func ReadConfigFile[T []Table | []TableMetaData](name string) (T, error) {
	f, err := os.Open(name)
	if err != nil {
		return *new(T), err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var content *T
	if err := dec.Decode(&content); err != nil {
		return *content, err
	}

	return *content, nil
}

func SyncDBFiles() {

}

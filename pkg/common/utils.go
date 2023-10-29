package common

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const DATA_FILE_NAME = ".store/data.bin"
const METADATA_FILE_PATH = ".store/metadata.bin"

var ResponseStrings = map[string]string{
	"T1": "Table with this id not found",
	"C1": "This column type is not allowed",
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

func Dump() {
	df := CreateFile("data.bin")
	mf := CreateFile("metadata.bin")
	defer df.Close()
	defer mf.Close()

	err := WriteFile(df, Store.Tables)
	if err != nil {
		log.Println("Write data failed")
	}
	err = WriteFile(mf, Store.TablesMetaData)
	if err != nil {
		log.Println("Write metadata failed")
	}
}

func Load() (bool, DatabaseStore) {
	data, err := ReadConfigFile[[]Table]("data.bin")
	if err != nil {
		log.Println(err)
		return false, DatabaseStore{}
	}
	metadata, err := ReadConfigFile[[]TableMetaData]("metadata.bin")
	if err != nil {
		log.Println(err)
		return false, DatabaseStore{}
	}

	return true, DatabaseStore{
		Tables:         data,
		TablesMetaData: metadata,
	}
}

func Config(configData DatabaseStore) {
	Store = configData
}

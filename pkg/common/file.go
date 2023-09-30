package common

import (
	"fmt"
	"log"
)

func Dump() {
	df := CreateFile("data.bin")
	mf := CreateFile("metadata.bin")
	defer df.Close()
	defer mf.Close()

	err := WriteFile(df, Store.Tables)
	if err != nil {
		fmt.Println(err)
		log.Println("Write data failed")
	}
	err = WriteFile(mf, Store.TablesMetaData)
	if err != nil {
		fmt.Println(err)
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

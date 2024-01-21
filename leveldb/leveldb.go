package leveldb

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

// AddDataToDatabase stores key-value pairs in the specified database.
func AddDataToDatabase(databasename string, key, value []byte) {
	db, err := leveldb.OpenFile(databasename, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Put(key, value, nil); err != nil {
		log.Fatalf("Failed to put data: %v", err)
	}
}

// FetchFromDatabase retrieves data for a given key from the specified database.
func FetchFromDatabase(databasename string, key []byte) []byte {
	db, err := leveldb.OpenFile(databasename, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	data, err := db.Get(key, nil)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
	return data
}

// // DeleteFromDatabase removes a key and its value from the specified database.
// func DeleteFromDatabase(databasename string, key []byte) {
// 	db, err := leveldb.OpenFile(databasename, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to open database: %v", err)
// 	}
// 	defer db.Close()

// 	if err := db.Delete(key, nil); err != nil {
// 		log.Fatalf("Failed to delete data: %v", err)
// 	}

// 	fmt.Println("Data deleted successfully!")
// }

// // GetDataFromLevelDB retrieves and unmarshals JSON data for a given key from the "database".
// func GetDataFromLevelDB(from string) (float64, error) {
// 	db, err := leveldb.OpenFile("database", nil)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to open database: %v", err)
// 	}
// 	defer db.Close()

// 	fromBalanceBytes, err := db.Get([]byte(from), nil)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get data: %v", err)
// 	}

// 	var amount float64
// 	if err := json.Unmarshal(fromBalanceBytes, &amount); err != nil {
// 		return 0, fmt.Errorf("failed to unmarshal data: %v", err)
// 	}

// 	return amount, nil
// }

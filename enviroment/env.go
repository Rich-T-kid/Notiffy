package enviroment

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	keys           = []string{"test", "test2"}
	_              = loadenviromentVaribles()
	errMissingKeys = errors.New("enviroments var's are not set")
)

func init() {
	fmt.Println("It critical that the env.go Files is imported first into Main") //atleast the first non std package
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mising, err := testExistance()
	if err != nil {
		log.Fatalf("keys: %s are missing with the following error %e\n", mising, err)
	}
	fmt.Println("Enviroment variables are set up and good to go!")

}
func loadenviromentVaribles() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading env file resulted in an error ->", err)
	}
	return err
}

func testExistance() ([]string, error) {
	var missingKeys []string
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			missingKeys = append(missingKeys, key)
		}
	}
	if len(missingKeys) > 0 {
		fmt.Println(missingKeys)
		return missingKeys, errMissingKeys
	}
	return missingKeys, nil
}

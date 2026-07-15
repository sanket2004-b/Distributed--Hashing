package test

import (
	"distributed-hashing/client/methods"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

type User struct {
	Name string
	Age  int
}

func DecodeValue[T any](data []byte) (T, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return t, err
}

func UnitTesting() {

	var wg sync.WaitGroup
	fmt.Printf("\n********Adding Keys to HashMap*****\n")
	for i := 21; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := "Sanket" + strconv.Itoa(i)
			key := "user" + strconv.Itoa(i)
			age := i
			user := User{Name: name, Age: age}
			err := methods.SetKeyValue(key, user)
			fmt.Printf("Added key: %v value: %+v\n", key, user)
			if err != nil {
				fmt.Printf("Error while setting :key: %v: %v \n", key, err.Error())
			}
		}(i)

	}
	wg.Wait()
}

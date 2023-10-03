package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const fileName = "accounts.json"

func main() {

	//At the begining of the pragram we will try to read the file
	var accounts Accounts
	m := make(map[string]AccountDetails)
	dir, _ := os.Getwd()
	data, err := os.ReadFile(dir + "/" + fileName)
	if err == nil {
		json.Unmarshal(data, &accounts)
	}

	for _, el := range accounts.Accounts {
		m[el.Id] = el.Details
	}
	fmt.Println("Welcome to password volt.")
	help()
	var input string
	fmt.Scanln(&input)
loop:
	for {
		switch input {
		case "add":
			add(&accounts, m)
		case "update":
			update(&accounts, m)
		case "get":
			get(m)
		case "exit":
			break loop
		default:
			help()
		}
		fmt.Println("How can I help you next?")
		fmt.Scanln(&input)
	}

	// At the end of the program we will rewrite the file for next use
	b, _ := json.Marshal(accounts)
	os.WriteFile(dir+"/"+fileName, b, 0700)

}

// For those function we will start without encryption
// Those three function define the flow of the app
func add(a *Accounts, m map[string]AccountDetails) {
	fmt.Println("Enter a unique identifier for the element you want to add, eg. Gmail, personal")
	var id string
	fmt.Scanln(&id)
	// We need to check if the element is already present in the map
	_, ok := m[id]
	if ok {
		fmt.Println("The id already existed use the update keyword instead")
		return
	}
	var details AccountDetails
	fmt.Println("Enter the username used to login")
	fmt.Scanln(&details.Username)
	fmt.Println("Enter the password used to login")
	var password string
	fmt.Scanln(&password)
	details.Password = []byte(password)
	details.encrypt()
	m[id] = details
	var info LoginInfo
	info.Details = details
	info.Id = id
	a.Accounts = append(a.Accounts, info)

}

func update(a *Accounts, m map[string]AccountDetails) {
	fmt.Println("update")
}

func get(m map[string]AccountDetails) {
	fmt.Println("Enter the id of the element you want to get")
	var id string
	fmt.Scanln(&id)
	// We need to check if the element is already present in the map
	val, ok := m[id]
	if ok {
		decryptVal := val.decrypt()
		fmt.Printf("usertname: %v\nPassword: %v\n",
			decryptVal.Username, string(decryptVal.Password))
	} else {
		fmt.Println("This id doesn't exist")
	}
}

func help() {
	fmt.Println("Here are the functionnalities of the program:")
	fmt.Println("add: Adds a new element to the password volt")
	fmt.Println("get: gets a specific element from the password volt")
	fmt.Println("update: Updates a value in the password volt")
}

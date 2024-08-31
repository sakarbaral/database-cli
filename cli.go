package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sakarbaral/database/models"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A CLI for interacting with the users database",
}

var db *Driver

func init() {
	var dir string

	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "./", "Directory to store data")

	var err error
	db, err = New(dir, nil)
	if err != nil {
		log.Fatalf("Error initializing driver: %v", err)
	}

	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(readAllCmd)
	rootCmd.AddCommand(deleteCmd)
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a user to the collection",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 5 {
			fmt.Println("Usage: write <collection> <name> <age> <contact> <company> <address>")
			return
		}

		collection := args[0]
		name := args[1]
		age := args[2]
		contact := args[3]
		company := args[4]
		address := args[5]

		user := models.User{
			Name:    name,
			Age:     json.Number(age),
			Contact: contact,
			Company: company,
			Address: models.Address{City: address},
		}

		if err := db.Write(collection, name, user); err != nil {
			fmt.Printf("Error writing user: %v\n", err)
		} else {
			fmt.Println("User written successfully")
		}
	},
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a user from the collection",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: read <collection> <name>")
			return
		}

		collection := args[0]
		name := args[1]

		var user models.User
		if err := db.Read(collection, name, &user); err != nil {
			fmt.Printf("Error reading user: %v\n", err)
		} else {
			fmt.Printf("User: %+v\n", user)
		}
	},
}

var readAllCmd = &cobra.Command{
	Use:   "readall",
	Short: "Read all users from a collection",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: readall <collection>")
			return
		}

		collection := args[0]

		records, err := db.ReadAll(collection)
		if err != nil {
			fmt.Printf("Error reading all records: %v\n", err)
			return
		}

		for _, record := range records {
			fmt.Println(record)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user from the collection",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: delete <collection> <name>")
			return
		}

		collection := args[0]
		name := args[1]

		if err := db.Delete(collection, name); err != nil {
			fmt.Printf("Error deleting user: %v\n", err)
		} else {
			fmt.Println("User deleted successfully")
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

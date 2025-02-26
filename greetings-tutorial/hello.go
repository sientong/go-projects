package main

import (
	"fmt"
	greetings "greetings/greetings_functions"
	"log"
)

func greetingsWithError() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	message, err := greetings.Hello("")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}

func randomGreetings() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	message, err := greetings.GenerateRandomGreetings("Gladys")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
}

func greetingsMultiplePeople() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"Gladys", "Samantha", "Darrin"}

	messages, err := greetings.GreetingsMultiplePeople(names)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}

func main() {
	// greetingsWithError()
	// randomGreetings()
	greetingsMultiplePeople()
}

package main

import (
	"fmt"
	"unsafe"
)

type Creature struct {
	Name    string
	isAlive bool
}

func foo() *Creature {
	myCreature := Creature{Name: "dino", isAlive: true}
	fmt.Printf("%p\n", &myCreature)
	fmt.Println(myCreature)
	return &myCreature
}

func main() {
	b := []byte("ABC♥")
	fmt.Println(b) // [65 66 67 226 130 172]
	fmt.Println(len(b))
	fmt.Printf("a: %T, %d\n", b, unsafe.Sizeof(b))
	fmt.Println(string([]byte{65, 66, 67, 226, 130, 172}))

	fileName := "thisIsAreallyLargeNameLikeWTF.txt"
	fileNameInBytes := []byte(fileName)
	fileNameInBytesSize := len(fileNameInBytes)
	println(fileNameInBytesSize)
	fmt.Printf("%v", fileNameInBytes)
	fmt.Println(int([]byte{65, 66}))
}

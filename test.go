package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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
	/* b := []byte("ABC♥")
	fmt.Println(b) // [65 66 67 226 130 172]
	fmt.Println(len(b))
	fmt.Printf("a: %T, %d\n", b, unsafe.Sizeof(b))
	fmt.Println(string([]byte{65, 66, 67, 226, 130, 172}))

	fileName := "thisIsAreallyLargeNameLikeWTF.txt"
	fileNameInBytes := []byte(fileName)
	fileNameInBytesSize := len(fileNameInBytes)
	println(fileNameInBytesSize)
	fmt.Printf("%v", fileNameInBytes)
	fmt.Println(int([]byte{65, 66})) */
	/* buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(184549376))
	fmt.Printf("buffer %v", buffer)
	fmt.Printf("value %d", int32(binary.LittleEndian.Uint32([]byte{11, 0, 0, 0}))) */
	src := strings.NewReader("GeeksforGeeks\n")

	// Defining destination using Stdout
	dst := os.Stdout

	// Calling CopyN method with its parameters
	_, _ = io.CopyN(dst, src, 5)
	println("kjadkjdsa" + "jds")
}

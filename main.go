package main

import (
	"github.com/biorhitm/memfs"
	"fmt"
	"syscall"
	//"unsafe"
)

/*
BNF-определения
RVALUE = <EXPRESSION>
EXPRESSION = <FLOAT> {'+|-|*|/|%|^' <FLOAT>}
FLOAT = <INT_NUM>['.'<INT_NUM>]
INT_NUM = <DIGIT>{<DIGIT>}
DIGIT = '0'..'9'
*/
func main() {
	mapIntf, err := memfs.Mmap("test.l")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	
	p := mapIntf.BaseAddress()
	lexem, errorCode, errorIndex := BuildLexems(p, mapIntf.GetSize())

	L := lexem
	for L != nil {
		b := make([]uint16, L.Size)
		for i := 0; i < L.Size; i++ {
			b[i] = L.Text[i]
		}
		S := syscall.UTF16ToString(b)

		fmt.Printf("Лехема: type: %d size: %d %s", L.Type, L.Size, S)
		L = L.Next
		fmt.Println()
	}
	
	if errorCode != 0 {
		fmt.Printf("Ошибка %d в %d\n", errorCode, errorIndex)
	}

	mapIntf.Munmap()
	mapIntf = nil
}

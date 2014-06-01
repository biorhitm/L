package main

import (
	"fmt"
	"github.com/biorhitm/memfs"
	"github.com/biorhitm/lsa"
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
	lexem, errorCode, errorIndex := lsa.BuildLexems(p, mapIntf.GetSize())

	L := lexem

	for L != nil {
		if L.Size > 0 {
			fmt.Printf("Лехема: %d size: %d %s\n", L.Type, L.Size, (*L).LexemAsString())
		}
		L = L.Next
	}

	fmt.Printf("----------EOF-----------\n")

	if errorCode != 0 {
		fmt.Printf("Ошибка %d в %d\n", errorCode, errorIndex)
	}

	mapIntf.Munmap()
	mapIntf = nil
}

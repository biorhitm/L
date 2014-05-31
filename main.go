package main

import (
	"fmt"
	"github.com/biorhitm/memfs"
	"syscall"
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
	var prevLexemType TLexemType = ltUnknown

	for L != nil {
		var S string
		if L.Size > 0 {
			b := make([]uint16, L.Size)
			for i := 0; i < L.Size; i++ {
				b[i] = L.Text[i]
			}
			S = syscall.UTF16ToString(b)
			if prevLexemType != ltUnknown && prevLexemType != ltEOL {
				//fmt.Printf(" ")
			}
			fmt.Printf("Лехема: %d size: %d %s ", L.Type, L.Size, S)
		}

		if L.Type == ltEOL {
			fmt.Printf(" \\n\n")
		}

		prevLexemType = L.Type
		L = L.Next
	}

	fmt.Printf("----------EOF-----------\n")

	if errorCode != 0 {
		fmt.Printf("Ошибка %d в %d\n", errorCode, errorIndex)
	}

	mapIntf.Munmap()
	mapIntf = nil
}

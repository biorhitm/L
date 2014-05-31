package main

import (
	"github.com/biorhitm/memfs"
	"unsafe"
)

type TBigShortArray [0x1FFFFFFFFFFFF]uint16
type PBigShortArray *TBigShortArray

type TLexemType uint

// TLexemType
const (
	ltUnknown = iota
	ltEOF
	ltNumber       // 12
	ltString       // "test"
	ltChar         // 'a' 'x' '%'
	ltSymbol       // ! @ # $ % ^ & * () - + = [] {} и т.д.
	ltIdent        // имя функции, переменной или типа
	ltReservedWord //функция, конец, если и т.д.
)

type TLexem struct {
	Next PLexem
	Text PBigShortArray
	Size int
	Type TLexemType
}

type PLexem *TLexem

func isLetter(C uint16) bool {
	return (0x0410 <= C && C <= 0x042F) ||
		(0x0430 <= C && C <= 0x044F) ||
		C == 0x0401 || C == 0x0451
}

func isIdentLetter(C uint16) bool {
	return C == '_' || isLetter(C)
}

func isDigit(C uint16) bool {
	return '0' <= C && C <= '9'
}

func createNewLexem(parent PLexem, text uint64, _type TLexemType) PLexem {
	L := new(TLexem)
	L.Text = PBigShortArray(unsafe.Pointer(uintptr(text)))
	L.Type = _type
	L.Size = 0
	L.Next = nil
	if parent != nil {
		parent.Next = L
	}
	return L
}

func BuildLexems(text memfs.PBigByteArray, size uint64) PLexem {
	var idx uint64 = 1
	var C uint16
	var curLexem, firstLexem PLexem
	var startIdx uint64

	var T PBigShortArray = nil

	if text[0] == 0xFF && text[1] == 0xFE {
		// типа Unicode
		T = PBigShortArray(unsafe.Pointer(text))
		size = (size - 2) / 2
	}

	addrOfText := uint64(uintptr(unsafe.Pointer(T)))

	curLexem = new(TLexem)
	firstLexem = curLexem

	for idx < size {
		C = T[idx]
		switch {
		case isIdentLetter(C):
			{
				startIdx = idx
				curLexem = createNewLexem(curLexem, addrOfText+idx*2, ltIdent)
				idx++
				for idx < size && (isIdentLetter(T[idx]) || isDigit(T[idx])) {
					idx++
				}
				curLexem.Size = int(idx - startIdx)
			}

		case isDigit(C):
			{
				startIdx = idx
				curLexem = createNewLexem(curLexem, addrOfText+idx*2, ltNumber)
				idx++
				for idx < size && isDigit(T[idx]) {
					idx++
				}
				curLexem.Size = int(idx - startIdx)
			}
			
		default:
			idx++
		}
	}

	createNewLexem(curLexem, 0, ltEOF)

	return firstLexem
}

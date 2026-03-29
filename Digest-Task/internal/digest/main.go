package digest

import (
	"math/cmplx"
	"strings"
	"unsafe"
)

// GetCharByIndex returns the i-th character from the given string.
func GetCharByIndex(str string, idx int) rune {
	currentIndex := 0
	for _, r := range str {
		if currentIndex == idx {
			return r
		}
		currentIndex++
	}
	panic("Index out of range")
}

// GetStringBySliceOfIndexes returns a string formed by concatenating specific characters from the input string based
// on the provided indexes.
func GetStringBySliceOfIndexes(str string, indexes []int) string {
	var builder strings.Builder
	builder.Grow(len(indexes))

	for _, value := range indexes {
		builder.WriteRune(GetCharByIndex(str, value))
	}

	return builder.String()
}

// ShiftPointer shifts the given pointer by the specified number of bytes using unsafe.Add.
func ShiftPointer(pointer **int, shift int) {
	if *pointer == nil {
		panic("pointer is nil")
	}

	*pointer = (*int)(unsafe.Add(unsafe.Pointer(*pointer), shift))
}

// IsComplexEqual compares two complex numbers and determines if they are equal.
func IsComplexEqual(a, b complex128) bool {
	epsilon := 1e-5
	if cmplx.IsInf(a) && cmplx.IsInf(b) {
		return a == b
	}

	return cmplx.Abs(a-b) < epsilon
}

// GetRootsOfQuadraticEquation returns two roots of a quadratic equation ax^2 + bx + c = 0.
func GetRootsOfQuadraticEquation(a, b, c float64) (complex128, complex128) {
	d := complex(b*b-4*a*c, 0)
	x1 := (complex(-b, 0) + cmplx.Sqrt(d)) / complex(2*a, 0)
	x2 := (complex(-b, 0) - cmplx.Sqrt(d)) / complex(2*a, 0)

	return x1, x2
}

// Sort sorts in-place the given slice of integers in ascending order.
func Sort(source []int) {
	if len(source) < 2 {
		return
	}

	left, right := 0, len(source)-1
	randomIndex := left + (right-left)/4
	pivot := source[randomIndex]

	for left <= right {
		for source[left] < pivot {
			left++
		}
		for source[right] > pivot {
			right--
		}
		if left <= right {
			source[left], source[right] = source[right], source[left]
			left++
			right--
		}
	}

	if 0 < right {
		Sort(source[:right+1])
	}

	if left < len(source) {
		Sort(source[left:])
	}
}

// ReverseSliceOne in-place reverses the order of elements in the given slice.
func ReverseSliceOne(s []int) {
	sSize := len(s)
	for i := 0; i < sSize/2; i++ {
		s[i], s[sSize-i-1] = s[sSize-i-1], s[i]
	}
}

// ReverseSliceTwo returns a new slice of integers with elements in reverse order compared to the input slice.
// The original slice remains unmodified.
func ReverseSliceTwo(s []int) []int {
	newSlice := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		newSlice[len(s)-i-1] = s[i]
	}

	return newSlice
}

// SwapPointers swaps the values of two pointers.
func SwapPointers(a, b *int) {
	if a == nil || b == nil {
		panic("a or b is nil")
	}
	*a, *b = *b, *a
}

// IsSliceEqual compares two slices of integers and returns true if they contain the same elements in the same order.
func IsSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// DeleteByIndex deletes the element at the specified index from the slice and returns a new slice.
// The original slice remains unmodified.
func DeleteByIndex(s []int, idx int) []int {
	resSlice := make([]int, 0, len(s)-1)
	resSlice = append(resSlice, s[:idx]...)
	resSlice = append(resSlice, s[idx+1:]...)

	return resSlice
}

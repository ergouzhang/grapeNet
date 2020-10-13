package continer

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSList_Sort(t *testing.T) {
	slist := New()
	for i := 0; i < 100; i++ {
		slist.Push(rand.Int() % 65556)
	}

	slist.Range(func(i interface{}) bool {
		fmt.Println(i)
		return true
	})

	slist.Sort(func(a, b interface{}) bool {
		av, aok := a.(int)
		bv, bok := b.(int)
		if aok && bok {
			return av > bv
		}

		return false
	})

	fmt.Println("sorted:")
	slist.Range(func(i interface{}) bool {
		fmt.Println(i)
		return true
	})

}

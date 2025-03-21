package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestList_Remove(t *testing.T) {
	t.Run("remove from empty list", func(t *testing.T) {
		l := NewList()
		item := &ListItem{Value: 10}

		l.Remove(item)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := make([]int, 0)

		assert.Nil(t, l.Front(), "front list element should be nil but got %v", item, l.Front())
		assert.Nil(t, l.Back(), "back list element should be nil but got %v", item, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("remove from single element list", func(t *testing.T) {
		l := NewList()

		item := l.PushBack(10)

		l.Remove(item)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := make([]int, 0)

		assert.Nil(t, l.Front(), "front list element should be nil but got %v", item, l.Front())
		assert.Nil(t, l.Back(), "back list element should be nil but got %v", item, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("remove the last element", func(t *testing.T) {
		l := NewList()

		item1 := l.PushBack(10)
		item2 := l.PushBack(20)
		item3 := l.PushBack(30)

		l.Remove(item3)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{10, 20}

		assert.Equal(t, item1, l.Front(), "front list element should be %v but got %v", item1, l.Front())
		assert.Equal(t, item2, l.Back(), "back list element should be %v but got %v", item2, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("remove the middle element", func(t *testing.T) {
		l := NewList()

		item1 := l.PushBack(10)
		item2 := l.PushBack(20)
		item3 := l.PushBack(30)

		l.Remove(item2)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{10, 30}

		assert.Equal(t, item1, l.Front(), "front list element should be %v but got %v", item1, l.Front())
		assert.Equal(t, item3, l.Back(), "back list element should be %v but got %v", item3, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("remove the last element", func(t *testing.T) {
		l := NewList()

		item1 := l.PushBack(10)
		item2 := l.PushBack(20)
		item3 := l.PushBack(30)

		l.Remove(item3)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{10, 20}

		assert.Equal(t, item1, l.Front(), "front list element should be %v but got %v", item1, l.Front())
		assert.Equal(t, item2, l.Back(), "back list element should be %v but got %v", item2, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})
}

func TestList_MoveToFront(t *testing.T) {
	t.Run("move to front in empty list", func(t *testing.T) {
		l := NewList()
		item := &ListItem{Value: 10}
		l.MoveToFront(item)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := make([]int, 0)

		assert.Nil(t, l.Front(), "front list element should be nil but got %v", item, l.Front())
		assert.Nil(t, l.Back(), "back list element should be nil but got %v", item, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("move to front in single element list", func(t *testing.T) {
		l := NewList()
		item := l.PushBack(10)
		l.MoveToFront(item)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{10}

		assert.Equal(t, item, l.Front(), "front list element should be %v but got %v", item, l.Front())
		assert.Equal(t, item, l.Back(), "back list element should be %v but got %v", item, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("move the last element to front", func(t *testing.T) {
		l := NewList()
		item1 := l.PushBack(10)
		item2 := l.PushBack(20)
		l.MoveToFront(item2)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{20, 10}

		assert.Equal(t, item2, l.Front(), "front list element should be %v but got %v", item2, l.Front())
		assert.Equal(t, item1, l.Back(), "back list element should be %v but got %v", item1, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("move the middle element to front", func(t *testing.T) {
		l := NewList()
		l.PushBack(10)
		l.PushBack(20)
		item3 := l.PushBack(30)
		l.PushBack(40)
		item5 := l.PushBack(50)

		l.MoveToFront(item3)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{30, 10, 20, 40, 50}

		assert.Equal(t, item3, l.Front(), "front list element should be %v but got %v", item3, l.Front())
		assert.Equal(t, item5, l.Back(), "back list element should be %v but got %v", item5, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("move the first element to front", func(t *testing.T) {
		l := NewList()

		item5 := l.PushFront(50)
		l.PushFront(40)
		l.PushFront(30)
		l.PushFront(20)
		item1 := l.PushFront(10)

		l.MoveToFront(item1)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{10, 20, 30, 40, 50}

		assert.Equal(t, item1, l.Front(), "front list element should be %v but got %v", item1, l.Front())
		assert.Equal(t, item5, l.Back(), "back list element should be %v but got %v", item5, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})

	t.Run("move the element to front several times", func(t *testing.T) {
		l := NewList()

		item5 := l.PushFront(50)
		item4 := l.PushFront(40)
		l.PushFront(30)
		l.PushFront(20)
		l.PushFront(10)

		l.MoveToFront(item4)
		l.MoveToFront(item4)
		l.MoveToFront(item4)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		expectedElems := []int{40, 10, 20, 30, 50}

		assert.Equal(t, item4, l.Front(), "front list element should be %v but got %v", item4, l.Front())
		assert.Equal(t, item5, l.Back(), "back list element should be %v but got %v", item5, l.Back())
		assert.Equal(t, len(expectedElems), l.Len(), "list should have %d elements but got %d", len(expectedElems), l.Len())
		assert.Equal(t, expectedElems, elems, "list should have items %v but got %v", expectedElems, elems)
	})
}

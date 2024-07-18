package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	size  int
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.first == nil {
		l.last = item
	} else {
		item.Next = l.first
		l.first.Prev = item
	}
	l.first = item

	l.size++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.last == nil {
		l.first = item
	} else {
		item.Prev = l.last
		l.last.Next = item
	}
	l.last = item

	l.size++
	return item
}

func (l *list) Remove(i *ListItem) {
	if (i.Prev == nil && l.first != i) || (i.Next == nil && l.last != i) {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.first = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}

	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	if (i.Prev == nil && l.first != i) || (i.Next == nil && l.last != i) {
		return
	}
	if l.first == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}

	i.Prev = nil
	i.Next = l.first

	if l.first != nil {
		l.first.Prev = i
	}
	l.first = i
	if l.last == nil {
		l.last = i
	}
}

func NewList() List {
	return &list{}
}

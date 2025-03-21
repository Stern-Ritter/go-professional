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
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.len == 0 {
		l.last = item
	} else {
		l.first.Prev = item
		item.Next = l.first
	}
	l.first = item

	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.len == 0 {
		l.first = item
	} else {
		l.last.Next = item
		item.Prev = l.last
	}
	l.last = item

	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	switch i {
	case l.first, l.last:
		if i == l.first {
			l.first = i.Next
			if i.Next != nil {
				i.Next.Prev = nil
			}
		}
		if i == l.last {
			l.last = i.Prev
			if i.Prev != nil {
				i.Prev.Next = nil
			}
		}
	default:
		if i.Prev == nil || i.Next == nil {
			return
		}
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.len == 0 {
		return
	}

	switch i {
	case l.first:
		return
	case l.last:
		l.last = i.Prev
		if i.Prev != nil {
			i.Prev.Next = nil
		}
	default:
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}
		if i.Next != nil {
			i.Next.Prev = i.Prev
		}
	}

	l.first.Prev = i
	i.Next = l.first
	l.first = i
}

func NewList() List {
	return &list{}
}

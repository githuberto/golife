package life

import (
	"fmt"
	"testing"
)

func TestStateString(t *testing.T) {
	tables := []struct {
		s State
		e string
	}{
		{Alive, "O"},
		{Dead, "X"},
		{99, "?"},
	}

	for _, table := range tables {
		a := table.s.String()
		if a != table.e {
			t.Errorf("State.String() was incorrect, got: %s, want: %s.",
				a, table.e)
		}
	}
}

func TestCellString(t *testing.T) {
	tables := []struct {
		c Cell
		e string
	}{
		{Cell{State: Alive}, "O"},
		{Cell{State: Dead}, "X"},
	}

	for _, table := range tables {
		a := table.c.String()
		if a != table.e {
			t.Errorf("Cell.String() was incorrect, got: %s, want: %s.",
				a, table.e)
		}
	}
}

func TestMultiplex(t *testing.T) {
	ch := make(chan bool)
	ds := MakeMultiplex(ch, 4)

	vs := []bool{true, true, false}

	for _, v := range vs {
		ch <- v
		for i := 0; i < len(ds); i++ {
			a := <-ds[i]
			if a != v {
				t.Errorf("Multiplexer failed on channel %d, got: %t, want %t",
					i, a, v)
			}
		}
	}
}

func TestInBounds(t *testing.T) {
	table := []struct {
		i int
		j int
		e bool
	}{
		{0, 0, true},
		{1, 1, true},
		{-1, 0, false},
		{0, -1, false},
		{2, 0, false},
		{0, 2, true},
	}

	// Simple 2x3 grid.
	b := [][]Cell{
		{{}, {}, {}},
		{{}, {}, {}},
	}

	for _, tb := range table {
		a := InBounds(b, tb.i, tb.j)
		if a != tb.e {
			t.Errorf("InBounds(b, %d, %d) was incorrect, got %t, want %t.",
				tb.i, tb.j, a, tb.e)
		}
	}
}

func TestNewState(t *testing.T) {
	table := []struct {
		s State
		n int
		e State
	}{
		{Alive, 0, Dead},
		{Alive, 1, Dead},
		{Alive, 2, Alive},
		{Alive, 3, Alive},
		{Alive, 4, Dead},
		{Alive, 5, Dead},
		{Alive, 6, Dead},
		{Alive, 7, Dead},
		{Alive, 8, Dead},
		{Dead, 0, Dead},
		{Dead, 1, Dead},
		{Dead, 2, Dead},
		{Dead, 3, Alive},
		{Dead, 4, Dead},
		{Dead, 5, Dead},
		{Dead, 6, Dead},
		{Dead, 7, Dead},
		{Dead, 8, Dead},
	}

	for _, tb := range table {
		a := NextState(tb.s, tb.n)
		if a != tb.e {
			t.Errorf("NewState(%d, %d) was incorrect, got %d, want %d",
				tb.s, tb.n, a, tb.e)
		}
	}
}

type A struct {
	ch *chan int
}

func AddChannel(a *A, b *A, ch *chan int) {
	a.ch = ch
	b.ch = ch
}

func TestAddTube(t *testing.T) {
	ch := make(chan int, 1)
	fmt.Printf("Type: %T\n", ch)
	a := A{}
	b := A{}

	AddChannel(&a, &b, &ch)

	*a.ch <- 5

	select {
	//case v := <-dst.tubes[South].in:
	case v := <-*b.ch:
		if v != 5 {
			t.Errorf("What?")
		}
	default:
		t.Errorf("WRONG")
	}
}

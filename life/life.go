package life

import (
  "fmt"
  "strconv"
  "strings"
  "sync"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

// Enums
type State int
const (
  Alive State = iota
  Dead
)

func (s State) String() string {
  if s == Alive {
    return "O"
  } else if s == Dead {
    return "X"
  } else  {
    return "?"
  }
}

type Direction int
const (
  North Direction = iota
  NorthEast
  East
  SouthEast
  South
  SouthWest
  West
  NorthWest
)

// Tube with in and out channels for a cell.
type Tube struct {
  in chan State
  out chan State
}

// Cell has a state and tubes for each direction.
type Cell struct {
  State State
  tubes [8]Tube
  sync chan bool
}

func (c Cell) String() string {
  if c.State == Alive {
    return "O"
  } else {
    return "X"
  }
}

func NextState(s State, n int) State {
  if s == Alive {
    if (n == 2) || (n == 3) {
      return Alive
    }
    return Dead
  } else {
    if n == 3 {
      return Alive
    }
    return Dead
  }
}

func (c *Cell) Evolve(wg sync.WaitGroup) {
  defer wg.Done()

  for range c.sync {
    // Send state to everyone else.
    for i := 0; i < len(c.tubes); i++ {
      if c.tubes[i].out != nil {
        c.tubes[i].out <- c.State
      }
    }
    // Now read from everyone else.
    n := 0
    for i := 0; i < len(c.tubes); i++ {
      if c.tubes[i].in != nil {
        s := <-c.tubes[i].in
        if s == Alive {
          n++
        }
      }
    }
    c.State = NextState(c.State, n)
  }
}


func AddTube(src *Cell, dst *Cell, dir Direction) {
  t := make(chan State, 1)
  opp := (dir + 4) % 8
  src.tubes[dir].out = t
  dst.tubes[opp].in = t
}

func InBounds(b [][]Cell, i int, j int) bool {
  return 0 <= i && i < len(b) && 0 <= j && j < len(b[i])
}

func Dir(i int, j int) (d Direction) {
  if (i == 1) && (j == 0) {
    return North 
  } else if (i == 1) && (j == 1) {
    return NorthEast 
  } else if (i == 0) && (j == 1) {
    return East 
  } else if (i == -1) && (j == 1) {
    return SouthEast 
  } else if (i == -1) && (j == 0) {
    return South 
  } else if (i == -1) && (j == -1) {
    return SouthWest 
  } else if (i == 0) && (j == -1) {
    return West 
  } else if (i == 1) && (j == -1) {
    return NorthWest 
  }

  panic("INVALID INDEX: " +  strconv.Itoa(i) + ", " + strconv.Itoa(j))
  return North
}

func Multiplex(s chan bool, ds []chan bool) {
  // Read from source and forward to dests.
  for v := range s {
    var wg sync.WaitGroup
    wg.Add(len(ds))
    for i := 0; i < len(ds); i++ {
      // Do this in new goroutines so they don't block on the write.
      go func(i int){ 
        defer wg.Done()
        ds[i] <- v 
      }(i)
    }
    wg.Wait()
  }

  // Now close all the dests.
  for i := 0; i < len(ds); i++ {
    close(ds[i])
  }
}

func MakeMultiplex(s chan bool, n int) []chan bool {
  ds := make([]chan bool, n)
  for i := 0; i < len(ds); i++ {
    // No buffer because we want them to tick.
    ds[i] = make(chan bool)
  }
  go Multiplex(s, ds)

  return ds
}

// Board utils
func LinkBoard(b [][]Cell, s chan bool) {
  ds := MakeMultiplex(s, len(b) * len(b[0]))

  for i := 0; i < len(b); i++ {
    for j:= 0; j < len(b[i]); j++ {
      // Assign each cell the sync channel.
      b[i][j].sync = ds[i*len(b[0]) + j]

      // Do the delta thing.
      for id := -1; id <= 1; id++ {
        for jd := -1; jd <= 1; jd++ {
          if (id == 0) && (jd == 0) {
            continue
          }
          ni, nj := i+id, j+jd
          if !InBounds(b, ni, nj) {
            continue
          }

          AddTube(&b[i][j], &b[ni][nj], Dir(id, jd))
        }
      }
    }
  }
}

func PrintBoard(b [][]Cell) {
  fmt.Println("----------------")
  for i := 0; i < len(b); i++ {
    for j := 0; j < len(b[i]); j++ {
      fmt.Print(b[i][j])
    }
    fmt.Println()
  }
}

func MakeBoard(bs string) [][]Cell {
  rs := strings.Split(bs, "\n")
  b := make([][]Cell, len(rs))
  for i := 0; i < len(b); i++ {
    b[i] = make([]Cell, len(rs[i]))

    // Set each row to its corresponding state.
    for j := 0; j < len(b[i]); j++ {
      if rs[i][j] == 'O' {
        b[i][j].State = Alive
      } else if rs[i][j] == 'X' {
        b[i][j].State = Dead
      } else {
        panic(fmt.Sprintf("invalid state: %s", rs[i][j]))
      }
    }
  }
  return b
}

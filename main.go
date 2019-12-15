package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "sync"
  "time"
  "life"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}



func main() {
  if len(os.Args) != 2 {
    panic(fmt.Sprintf("usage: %s <board_file>", os.Args[0]))
  }

  path := os.Args[1]

  dat, err := ioutil.ReadFile(path)
  check(err)
  bs := string(dat)
  b := MakeBoard(bs)

  t := make(chan bool, 1)
  LinkBoard(b, t)

  var wg sync.WaitGroup
  wg.Add(len(b)*len(b[0]))

  for i := 0; i < len(b); i++ {
    for j :=0; j < len(b[i]); j++ {
      go b[i][j].Evolve(wg)
    }
  }

  for i := 0; i < 10; i++ {
    t <- true
    PrintBoard(b)
    time.Sleep(2*time.Second)
  }
  wg.Wait()
}

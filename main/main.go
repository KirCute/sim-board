package main

import (
	"io/fs"

	"github.com/KirCute/sim-board"
	_ "github.com/KirCute/sim-board/deck"
	"github.com/KirCute/sim-board/public"
)

func main() {
	f, err := fs.Sub(public.Public, "dist")
	if err != nil {
		panic(err)
	}
	sim_board.Run(f)
}

package main

import (
	"strconv"

	"github.com/zan8in/gologger"
	"github.com/zan8in/gologger/levels"
)

func main() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)

	gologger.Print().Msgf("\tgologger: sample test\t\n")
	gologger.Info().Str("alex", "to").Msg("running.")

	for i := 0; i < 10; i++ {
		gologger.Info().Str("count", strconv.Itoa(i)).Msg("running simulation step...")
	}

	gologger.Debug().Str("state", "running").Msg("planning running.")
	gologger.Warning().Str("state", "errored").Str("status", "404").Msg("could not run")
	gologger.Fatal().Msg("bye bye")
}

package main

import (
	"flag"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"
	"time"

	"github.com/merzzzl/image-duplicate-detector/internal/app"
)

func main() {
	var in string
	var sim string
	var uniq string
	var move bool

	flag.StringVar(&in, "i", "", "")
	flag.StringVar(&sim, "s", "", "")
	flag.StringVar(&uniq, "u", "", "")
	flag.BoolVar(&move, "m", false, "")
	flag.Parse()

	inPath := strings.Split(in, ",")

	l := app.GetLogger()

	l.Info().Strs("input", inPath).Str("unique", uniq).Str("similar", sim).Msg("start")

	matches, files, err := app.FindDuplicates(inPath)
	if err != nil {
		log.Fatal(err)
	}

	return

	moved := app.MoveMatches(matches, uniq, sim, move)
	fix, _ := app.FixDates([]string{uniq})

	l.Info().Int("files", files).Int("move", moved).Int("update", fix).Msg("end")
	time.Sleep(time.Second)
}

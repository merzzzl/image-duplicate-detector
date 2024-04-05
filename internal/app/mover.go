package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/merzzzl/image-duplicate-detector/internal/service"
)

func MoveMatches(m []MatchSet, uniq, sim string, move bool) int {
	oks := make(map[string]struct{})
	timeS := time.Now()

	var fn func(from, to string) error

	if move {
		fn = os.Rename
	} else {
		fn = copy
	}

	logger.SetStage("Moving")
	logger.SetStageStatus("0%")

	for i, match := range m {
		if _, ok := oks[match.Best.Name()]; !ok {
			name := fmt.Sprintf("%s/%d%s", uniq, match.Best.Index(), strings.ToLower(filepath.Ext(match.Best.Name())))
			if err := fn(match.Best.Name(), name); err != nil {
				logger.Err(err).Str("file", match.Best.Name()).Str("to", uniq).Msg("move/copy failed")
			}

			service.SetMeta(name, match.metadata)

			oks[match.Best.Name()] = struct{}{}
		}

		if match.Worst != nil {
			for _, worst := range match.Worst {
				if _, ok := oks[worst.Name()]; !ok {
					if err := fn(worst.Name(), fmt.Sprintf("%s/%d-%d%s", sim, match.Best.Index(), worst.Index(), strings.ToLower(filepath.Ext(worst.Name())))); err != nil {
						logger.Err(err).Str("file", worst.Name()).Str("to", sim).Msg("move/copy failed")
					}

					oks[worst.Name()] = struct{}{}
				}
			}
		}

		logger.SetStageStatus(calcProgressString(i+1, len(m), timeS))
	}

	logger.Info().TimeDiff("duration", time.Now(), timeS).Msg("end of move files")

	return len(oks)
}

func copy(from, to string) error {
	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(to)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = r.WriteTo(w)

	return err
}

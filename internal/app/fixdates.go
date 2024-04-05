package app

import (
	"os"
	"time"

	"github.com/merzzzl/image-duplicate-detector/internal/service"
)

func FixDates(inputPaths []string) (int, int) {
	var files []string
	var modified int

	logger.SetStage("Fixing")
	logger.SetStageStatus("0%")

	for _, path := range inputPaths {
		images, err := service.ListImage(path)
		if err != nil {
			return 0, 0
		}

		videos, err := service.ListVideo(path)
		if err != nil {
			return 0, 0
		}

		files = append(files, images...)
		files = append(files, videos...)
	}

	timeS := time.Now()

	for i, file := range files {
		logger.SetStageStatus(calcProgressString(i+1, len(files), timeS))

		t, ok := service.FileTime(file)
		if !ok {
			logger.Warn().Str("file", file).Msg("failed to get time from the file")

			continue
		}

		if err := os.Chtimes(file, t, t); err != nil {
			logger.Error().Err(err).Str("file", file).Msg("failed to change times of file")

			continue
		}

		modified++
	}

	return modified, len(files)
}

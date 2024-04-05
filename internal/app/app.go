package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/merzzzl/image-duplicate-detector/internal/service"
)

type Match struct {
	Best  *service.Object
	Worst *service.Object
	Hash  float64
	Sim   float64
}

type MatchSet struct {
	Best     *service.Object
	Worst    []*service.Object
	metadata map[string]interface{}
}

func FindDuplicates(inputPaths []string) ([]MatchSet, int, error) {
	var files []string

	for _, path := range inputPaths {
		images, err := service.ListImage(path)
		if err != nil {
			return nil, 0, err
		}

		videos, err := service.ListVideo(path)
		if err != nil {
			return nil, 0, err
		}

		files = append(files, images...)
		files = append(files, videos...)
	}

	objects := load(files)
	matches := compare(objects)

	for _, obj := range objects {
		obj.Close()
	}

	return groupMatches(matches), len(files), nil
}

func groupMatches(matches []Match) []MatchSet {
	set := make([]MatchSet, 0)
	timeS := time.Now()

	logger.SetStage("Grouping")
	logger.SetStageStatus("0%")

	for i, match1 := range matches {
		tmp := MatchSet{Best: match1.Best, Worst: nil}
		worstList := make(map[int64]struct{})
		worstList[match1.Best.Index()] = struct{}{}

		for j, match2 := range matches {
			logger.SetStageStatus(calcProgressString(i*len(matches)+j, len(matches)*len(matches), timeS))

			if _, ok := worstList[match2.Best.Index()]; ok && match2.Worst != nil {
				worstList[match2.Worst.Index()] = struct{}{}
				tmp.Worst = append(tmp.Worst, match2.Worst)
			}
		}

		tmp.metadata = service.MergeMeta(tmp.Worst)

		set = append(set, tmp)
	}

	return set
}

func calcProgressString[T int | int32 | int64](now, total T, start time.Time) string {
	p := (float32)(now) / ((float32)(total)) * 100.0
	d := time.Since(start).Minutes() / (float64)(now) * ((float64)(total) - (float64)(now))

	return fmt.Sprintf("%.2f%% ~%.0fmin", p, d)
}

func load(files []string) []*service.Object {
	objects := make([]*service.Object, 0, len(files))
	indexer := service.NewIndexer()
	semaphore := make(chan struct{}, 20)
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	timeS := time.Now()

	logger.SetStage("Loading")
	logger.SetStageStatus("0%")

	for i, file := range files {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(i int, file string) {
			defer func() {
				wg.Done()
				<-semaphore
			}()

			object, ok := service.Load(file, indexer)
			if !ok {
				logger.Warn().Str("file", file).Msg("failed to load file")

				return
			}

			mx.Lock()
			objects = append(objects, object)
			mx.Unlock()

			logger.Debug().Str("file", file).Msgf("file load")
			logger.SetStageStatus(calcProgressString(i+1, len(files), timeS))
		}(i, file)
	}

	wg.Wait()

	logger.Info().TimeDiff("duration", time.Now(), timeS).Msg("end of load files")

	return objects
}

func compare(objects []*service.Object) []Match {
	matches := make([]Match, 0, len(objects))
	timeS := time.Now()

	logger.SetStage("Compare")
	logger.SetStageStatus("0%")

	for i1 := 0; i1 < len(objects); i1++ {
		var found bool

		for i2 := i1 + 1; i2 < len(objects); i2++ {
			logger.SetStageStatus(calcProgressString(i1*len(objects)+i2, len(objects)*len(objects), timeS))

			if f := objects[i1].Distance(objects[i2]); f < 15 {
				if sim := objects[i1].Similar(objects[i2]); sim < 15 {
					logger.Info().Float64("dis", f).Float64("sim", sim).Strs("files", []string{objects[i1].Name(), objects[i2].Name()}).Msg("found similar")
					found = true

					if objects[i1].Compare(objects[i2]) < 0 {
						matches = append(matches, Match{
							Best:  objects[i2],
							Worst: objects[i1],
							Sim:   sim,
							Hash:  f,
						})

						objects = append(objects[:i1], objects[i1+1:]...)

						i1--

						break
					}

					matches = append(matches, Match{
						Best:  objects[i1],
						Worst: objects[i2],
						Sim:   sim,
						Hash:  f,
					})

					objects = append(objects[:i2], objects[i2+1:]...)

					i2--
				} else if f < 1 {
					logger.Warn().Float64("dis", f).Float64("sim", sim).Strs("files", []string{objects[i1].Name(), objects[i2].Name()}).Msg("small distance")
				}
			}
		}

		if !found {
			matches = append(matches, Match{
				Best: objects[i1],
			})
		}
	}

	logger.Info().TimeDiff("duration", time.Now(), timeS).Msg("end of compare objects")

	return matches
}

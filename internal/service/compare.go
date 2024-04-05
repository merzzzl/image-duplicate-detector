package service

import (
	"os"

	"gocv.io/x/gocv"
)

func compareSizes(f1, f2 string) int {
	s1, err1 := os.Stat(f1)
	s2, err2 := os.Stat(f2)

	if err1 != nil || err2 != nil {
		return 0
	}

	switch {
	case s1.Size() > s2.Size():
		return 1
	case s1.Size() < s2.Size():
		return -1
	default:
		return 0
	}
}

func comapreTimes(f1, f2 string) int {
	t1, ok1 := FileTime(f1)
	t2, ok2 := FileTime(f2)

	if !ok1 || !ok2 {
		return 0
	}

	switch {
	case t1.Before(t2):
		return 1
	case t1.After(t2):
		return -1
	default:
		return 0
	}
}

func comapreExtensions(f1, f2 string) int {
	w1 := extensionsWeight(f1)
	w2 := extensionsWeight(f2)

	switch {
	case w1 > w2:
		return 1
	case w1 < w2:
		return -1
	default:
		return 0
	}
}

func compareResolutions(f1, f2 string) int {
	r1, ok1 := mediaResolutions(f1)
	r2, ok2 := mediaResolutions(f2)

	if !ok1 || !ok2 {
		return 0
	}

	switch {
	case r1 > r2:
		return 1
	case r1 < r2:
		return -1
	default:
		return 0
	}
}

func mediaResolutions(inputPath string) (int, bool) {
	switch whichType(inputPath) {
	case fileTypeImage:
		return imageResolutions(inputPath)
	case fileTypeVideo:
		return videoResolutions(inputPath)
	default:
		return 0, false
	}
}

func imageResolutions(inputPath string) (int, bool) {
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return 0, false
	}
	defer img.Close()

	return img.Cols() * img.Rows(), true
}

func videoResolutions(inputPath string) (int, bool) {
	vdx, err := gocv.VideoCaptureFile(inputPath)
	if err != nil {
		return 0, false
	}
	defer vdx.Close()

	img := gocv.NewMat()
	defer img.Close()

	if !vdx.Read(&img) {
		return 0, false
	}

	return img.Cols() * img.Rows(), true
}

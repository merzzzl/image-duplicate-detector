package service

import (
	"image"

	"gocv.io/x/gocv"
)

func Load(inputPath string, i *Indexer) (*Object, bool) {
	if mat, ok := loadAsMat(inputPath); ok && !mat.Empty() {
		return &Object{
			name:  inputPath,
			index: i.GetVal(),
			raw:   mat,
			hash:  gocv.NewMat(),
			sfit:  gocv.NewMat(),
			orb:   gocv.NewMat(),
		}, true
	}

	return nil, false
}

func loadAsMat(inputPath string) (gocv.Mat, bool) {
	switch whichType(inputPath) {
	case fileTypeImage:
		return imageLoad(inputPath)
	case fileTypeVideo:
		return videoLoad(inputPath)
	default:
		return gocv.Mat{}, false
	}
}

func imageLoad(inputPath string) (gocv.Mat, bool) {
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return gocv.Mat{}, false
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	small := gocv.NewMat()
	defer small.Close()

	gocv.Resize(gray, &small, image.Point{400, 400}, 0, 0, gocv.InterpolationLinear)

	if small.Empty() {
		return gocv.Mat{}, false
	}

	return small.Clone(), true
}

func videoLoad(inputPath string) (gocv.Mat, bool) {
	vdx, err := gocv.VideoCaptureFile(inputPath)
	if err != nil {
		return gocv.Mat{}, false
	}
	defer vdx.Close()

	img := gocv.NewMat()
	defer img.Close()

	if !vdx.Read(&img) {
		return gocv.Mat{}, false
	}

	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	small := gocv.NewMat()
	defer small.Close()

	gocv.Resize(gray, &small, image.Point{400, 400}, 0, 0, gocv.InterpolationLinear)

	if small.Empty() {
		return gocv.Mat{}, false
	}

	return small.Clone(), true
}

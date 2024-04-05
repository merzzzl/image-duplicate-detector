package service

import (
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

type Object struct {
	name  string
	index int64
	raw   gocv.Mat
	hash  gocv.Mat
	sfit  gocv.Mat
	orb   gocv.Mat
	kp    []gocv.KeyPoint
}

var hash = contrib.AverageHash{}
var sift = gocv.NewSIFT()
var orb = gocv.NewORB()
var matcher = gocv.NewBFMatcherWithParams(gocv.NormL2, false)

func (o *Object) Close() {
	o.raw.Close()
	o.hash.Close()
	o.sfit.Close()
	o.orb.Close()
	o.kp = nil
}

func (o *Object) Name() string {
	return o.name
}

func (o *Object) Index() int64 {
	return o.index
}

func (o *Object) Distance(obj *Object) float64 {
	if !o.GenHash() || !obj.GenHash() {
		return 100
	}

	return hash.Compare(o.hash, obj.hash)
}

func (o *Object) GenHash() bool {
	if o.raw.Empty() {
		return false
	}

	if !o.hash.Empty() {
		return true
	}

	results := gocv.NewMat()
	defer results.Close()

	hash.Compute(o.raw, &results)

	o.hash = results.Clone()

	return true
}

func (o *Object) GenSIFT() bool {
	if o.raw.Empty() {
		return false
	}

	if !o.sfit.Empty() {
		return true
	}

	kp, ds := sift.DetectAndCompute(o.raw, gocv.NewMat())
	defer ds.Close()

	o.sfit = ds.Clone()
	o.kp = append(o.kp, kp...)

	return true
}

func (o *Object) GenORB() bool {
	if o.raw.Empty() {
		return false
	}

	if !o.orb.Empty() {
		return true
	}

	kp, ds := orb.DetectAndCompute(o.raw, gocv.NewMat())
	defer ds.Close()

	o.orb = ds.Clone()
	o.kp = append(o.kp, kp...)

	return true
}

func (o *Object) Similar(obj *Object) float64 {
	if !o.GenORB() || !o.GenSIFT() || !obj.GenORB() || !obj.GenSIFT() || o.sfit.Type() != obj.sfit.Type() || o.orb.Type() != obj.orb.Type() {
		return 100
	}

	var matches [][]gocv.DMatch

	matches = append(matches, matcher.KnnMatch(o.sfit, obj.sfit, 2)...)
	matches = append(matches, matcher.KnnMatch(o.orb, obj.orb, 2)...)

	var goodMatches []gocv.DMatch
	for _, dm := range matches {
		if len(dm) != 2 {
			continue
		}

		if dm[0].Distance < dm[1].Distance*0.75 {
			goodMatches = append(goodMatches, dm[0])
		}
	}

	numKeypoints := len(o.kp)
	if len(obj.kp) > numKeypoints {
		numKeypoints = len(obj.kp)
	}

	percentDiff := 100.0 * (1.0 - float64(len(goodMatches))/float64(numKeypoints))

	return percentDiff
}

func (o Object) Compare(obj *Object) int {
	if w := compareResolutions(o.name, obj.name); w != 0 {
		return w
	}

	if w := comapreExtensions(o.name, obj.name); w != 0 {
		return w
	}

	if w := comapreTimes(o.name, obj.name); w != 0 {
		return w
	}

	if w := compareSizes(o.name, obj.name); w != 0 {
		return w
	}

	return 0
}

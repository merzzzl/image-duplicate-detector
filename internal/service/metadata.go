package service

import (
	"time"

	"github.com/barasher/go-exiftool"
)

var (
	imageExtensionsWeight = map[string]int{
		".png":  10,
		".bmp":  10,
		".jpg":  5,
		".jpeg": 5,
	}
	videoExtensionsWeight = map[string]int{
		".mov": 8,
		".avi": 8,
		".mp4": 5,
		".m4v": 5,
		".mpg": 5,
	}
	imageTimeTags = []string{
		"DateTimeOriginal",
		"DateTimeDigitized",
		"DateTime",
		"ModifyDate",
		"SubsecTime",
		"SubsecTimeOriginal",
		"SubsecTimeDigitized",
	}
	videoTimeTags = []string{
		"MediaCreateDate",
		"TrackCreateDate",
		"CreationTime",
		"MediaModifiedDate",
		"TrackModifiedDate",
		"ModificationTime",
		"EncodingTime",
	}
	fileTimeTags = []string{
		"FileModifyDate",
		"FileAccessDate",
		"FileInodeChangeDate",
	}
	exif = newExifTool()
)

func newExifTool() *exiftool.Exiftool {
	tool, err := exiftool.NewExiftool()
	if err != nil {
		panic(err)
	}

	return tool
}

func extensionsWeight(ext string) int {
	if weight, ok := imageExtensionsWeight[ext]; ok {
		return weight
	}

	if weight, ok := videoExtensionsWeight[ext]; ok {
		return weight
	}

	return 0
}

func FileTime(path string) (time.Time, bool) {
	tstr := originTime(path)
	if tstr == "" {
		return time.Now(), false
	}

	location, err := time.LoadLocation("Local")
	if err != nil {
		return time.Now(), false
	}

	if len(tstr) == 25 {
		offset, err := time.ParseDuration(tstr[20:22] + "h" + tstr[23:25] + "m")
		if err != nil {
			return time.Now(), false
		}

		tstr = tstr[:19]
		location = time.FixedZone("exif", int(offset.Seconds()))
	}

	if len(tstr) > 19 {
		tstr = tstr[:19]
	}

	t, err := time.ParseInLocation("2006:01:02 15:04:05", tstr, location)
	if err != nil {
		return time.Now(), false
	}

	return t, true
}

func MergeMeta(objects []*Object) map[string]interface{} {
	var paths []string

	for _, o := range objects {
		paths = append(paths, o.name)
	}

	meta := exif.ExtractMetadata(paths...)

	var merged = make(map[string]interface{})

	for _, m := range meta {
		for f, v := range m.Fields {
			if _, ok := merged[f]; !ok {
				merged[f] = v
			}
		}
	}

	return merged
}

func SetMeta(file string, fields map[string]interface{}) {
	meta := exif.ExtractMetadata(file)

	for f, v := range fields {
		if _, ok := meta[0].Fields[f]; !ok {
			meta[0].Fields[f] = v
		}
	}    

	exif.WriteMetadata(meta)
}

func originTime(path string) string {
	meta := exif.ExtractMetadata(path)

	if len(meta) == 0 {
		return ""
	}

	filemeta := meta[0]

	for _, tag := range imageTimeTags {
		t, err := filemeta.GetString(tag)
		if err != nil || t == "" || t == "0000:00:00 00:00:00" {
			continue
		}

		return t
	}

	for _, tag := range videoTimeTags {
		t, err := filemeta.GetString(tag)
		if err != nil || t == "" || t == "0000:00:00 00:00:00" {
			continue
		}

		return t
	}

	for _, tag := range fileTimeTags {
		t, err := filemeta.GetString(tag)
		if err != nil || t == "" || t == "0000:00:00 00:00:00" {
			continue
		}

		return t
	}

	return ""
}

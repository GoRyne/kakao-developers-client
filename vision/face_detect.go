package vision

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"internal/common"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// Face represents data of the detected face.
type Face struct {
	FacialAttributes FacialAttributes `json:"facial_attributes"`
	FacialPoints     FacialPoints     `json:"facial_points"`
	Score            float64          `json:"score"`
	ClassIdx         int              `json:"class_idx"`
	X                float64          `json:"x"`
	Y                float64          `json:"y"`
	W                float64          `json:"w"`
	H                float64          `json:"h"`
	Pitch            float64          `json:"pitch"`
	Yaw              float64          `json:"yaw"`
	Roll             float64          `json:"roll"`
}

// FacialAttributes represents estimated gender, age of the detected face.
type FacialAttributes struct {
	Gender Gender  `json:"gender"`
	Age    float64 `json:"age"`
}

// Gender represents confidence score that the detected face is considered as male or female.
type Gender struct {
	Male   float64 `json:"male"`
	Female float64 `json:"female"`
}

// FacialPoints represents arrays of coordinates of the detected face. (a value between 0 and 1.0)
type FacialPoints struct {
	Jaw          [][]float64 `json:"jaw"`
	RightEyebrow [][]float64 `json:"right_eyebrow"`
	LeftEyebrow  [][]float64 `json:"left_eyebrow"`
	Nose         [][]float64 `json:"nose"`
	RightEye     [][]float64 `json:"right_eye"`
	LeftEye      [][]float64 `json:"left_eye"`
	Lip          [][]float64 `json:"lip"`
}

// FaceResult represents a document of a detected face result.
type FaceResult struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Faces  []Face `json:"faces"`
}

// FaceDetectResult represents a Face Detection result.
type FaceDetectResult struct {
	RID    string     `json:"rid"`
	Result FaceResult `json:"result"`
}

// String implements fmt.Stringer.
func (fr FaceDetectResult) String() string { return common.String(fr) }

// SaveAs saves fr to @filename
//
// The file extension must be .json.
func (fr FaceDetectResult) SaveAs(filename string) error { return common.SaveAsJSON(fr, filename) }

// FaceDetectInitializer is a lazy face detector.
type FaceDetectInitializer struct {
	AuthKey   string
	Filename  string
	ImageURL  string
	Threshold float64
	withFile  bool
}

// FaceDetect detects a face in the given image.
//
// Image can be either image URL or image file (JPG or PNG).
// Refer to https://developers.kakao.com/docs/latest/ko/vision/dev-guide#recog-face for more details.
func FaceDetect() *FaceDetectInitializer {
	return &FaceDetectInitializer{
		AuthKey:   common.KeyPrefix,
		ImageURL:  "",
		Filename:  "",
		Threshold: 0.7,
	}

}

// WithURL sets url to @url.
func (fi *FaceDetectInitializer) WithURL(url string) *FaceDetectInitializer {
	fi.ImageURL = url
	fi.withFile = false
	return fi
}

// WithFile sets image path to @filename.
func (fi *FaceDetectInitializer) WithFile(filename string) *FaceDetectInitializer {
	fi.Filename = filename
	fi.withFile = true
	return fi
}

// AuthorizeWith sets the authorization key to @key.
func (fi *FaceDetectInitializer) AuthorizeWith(key string) *FaceDetectInitializer {
	fi.AuthKey = common.FormatKey(key)
	return fi
}

// ThresholdAt sets the Threshold to @val. (a value between 0 and 1.0)
//
// Threshold is a reference value to detect as a face.
// If @val is set too high, some faces may not be able to be detected as a face.
// If @val is set too low, other area can be detected as a face.
func (fi *FaceDetectInitializer) ThresholdAt(val float64) *FaceDetectInitializer {
	if 0.1 <= val && val <= 1.0 {
		fi.Threshold = val
	} else {
		panic(errors.New("threshold must be between 0.1 and 1.0"))
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return fi
}

// Collect returns the face detection result.
func (fi *FaceDetectInitializer) Collect() (res FaceDetectResult, err error) {
	client := &http.Client{}
<<<<<<< HEAD
	var req *http.Request
=======
>>>>>>> upstream/master

	if fi.withFile {

		file, err := os.Open(fi.Filename)
		if err != nil {
			return res, err
		}

		if stat, _ := file.Stat(); 2*1024*1024 < stat.Size() {
			return res, err
		}

		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("threshold", fmt.Sprintf("%f", fi.Threshold))
		part, err := writer.CreateFormFile("image", fi.Filename)

		if err != nil {
			return res, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return res, err
		}
		writer.Close()

		req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/face/detect", prefix), body)
		if err != nil {
			return res, err
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())

	} else {
		req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/face/detect?threshold=%f&image_url=%s", prefix, fi.Threshold, fi.ImageURL), nil)
		if err != nil {
			return res, err
		}
	}

	req.Close = true

	req.Header.Add(common.Authorization, fi.AuthKey)
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, err
	}
	return

}

package loader

import (
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/r2dtools/mnist/downloader"
)

const mnistUri = "http://yann.lecun.com/exdb/mnist/"

const (
	trainImageFileName = "train-images-idx3-ubyte.gz"
	trainLabelFileName = "train-labels-idx1-ubyte.gz"
	testImageFileName  = "t10k-images-idx3-ubyte.gz"
	testLabelFileName  = "t10k-labels-idx1-ubyte.gz"
)

const (
	ImageWidth  = 28
	ImageHeight = 28
)
const ImageLength = ImageWidth * ImageHeight

const (
	imageMagic = 0x00000803
	labelMagic = 0x00000801
)

var (
	ErrInvalidFormat = errors.New("mnist: invalid format")
	ErrCountMismatch = errors.New("mnist: images and labels count mismatch")
)

type Image [ImageWidth * ImageHeight]byte
type Label uint8

type Set struct {
	Images []Image
	Labels []Label
}

func (s *Set) Count() int {
	return len(s.Labels)
}

func (s *Set) Get(i int) (Image, Label) {
	return s.Images[i], s.Labels[i]
}

func (s *Set) Slice(i, j int) *Set {
	return &Set{s.Images[i:j], s.Labels[i:j]}
}

type imageFileHeader struct {
	Magic     int32
	NumImages int32
	Height    int32
	Width     int32
}

type labelFileHeader struct {
	Magic     int32
	NumLabels int32
}

func LoadData(directory string) (train, test *Set, err error) {
	fileNames := []string{trainImageFileName, trainLabelFileName, testImageFileName, testLabelFileName}
	fileMap := make(map[string]string)

	if directory == "" {
		directory = os.TempDir()
	}

	for _, fileName := range fileNames {
		filePath := filepath.Join(directory, fileName)
		_, err := os.Stat(filePath)
		// download not existed files
		if err != nil && os.IsNotExist(err) {
			fileMap[fileName] = fmt.Sprintf("%s/%s", mnistUri, fileName)
		}
	}

	if err := downloader.DownloadFiles(directory, fileMap); err != nil {
		return nil, nil, err
	}

	train, err = loadSet(
		filepath.Join(directory, trainImageFileName),
		filepath.Join(directory, trainLabelFileName),
	)
	if err != nil {
		return nil, nil, err
	}

	test, err = loadSet(
		filepath.Join(directory, testImageFileName),
		filepath.Join(directory, testLabelFileName),
	)
	if err != nil {
		return nil, nil, err
	}

	return train, test, err
}

func loadSet(imageFilePath, labelFilePath string) (*Set, error) {
	images, err := loadImages(imageFilePath)
	if err != nil {
		return nil, err
	}

	labels, err := loadLabels(labelFilePath)
	if err != nil {
		return nil, err
	}

	if len(images) != len(labels) {
		return nil, ErrCountMismatch
	}

	return &Set{images, labels}, nil
}

func loadImages(imageFilePath string) ([]Image, error) {
	file, err := os.Open(imageFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	header := imageFileHeader{}
	if err = binary.Read(gzipReader, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != imageMagic ||
		header.Width != ImageWidth ||
		header.Height != ImageHeight {
		return nil, ErrInvalidFormat
	}

	images := make([]Image, header.NumImages)
	for i := int32(0); i < header.NumImages; i++ {
		image := Image{}
		if err := binary.Read(gzipReader, binary.BigEndian, &image); err != nil {
			return nil, err
		}

		images[i] = image
	}

	return images, nil
}

func loadLabels(labelFilePath string) ([]Label, error) {
	file, err := os.Open(labelFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	header := labelFileHeader{}
	if err = binary.Read(gzipReader, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != labelMagic {
		return nil, err
	}

	labels := make([]Label, header.NumLabels)
	for i := int32(0); i < header.NumLabels; i++ {
		err = binary.Read(gzipReader, binary.BigEndian, &labels[i])
		if err != nil {
			return nil, err
		}
	}

	return labels, nil
}

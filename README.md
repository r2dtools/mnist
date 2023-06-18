# MNIST Golang

<p>The MNIST database is a dataset of handwritten digits. It has 60,000 training samples, and 10,000 test samples. Each image is represented by 28x28 pixels, each containing a value 0 - 255 with its grayscale value.
For more information and to download the database, see <http://yann.lecun.com/exdb/mnist/>.</p>
<p>
Package provides a simple interface to parse and use the MNIST database. It can automatically download the database files (and cache) or you can download the database files manually to be easly loaded with the LoadData function.
Also the package supports convertation of base data structures to [gonum](https://www.gonum.org/) matrices. 
</p>

## Automatically download and parse MNIST database files
```go
import github.com/r2dtools/mnist/loader

train, test, err := loader.LoadData("")
if err != nil {
    panic(err)
}

trainSlice := train.Slice(0, 1000)

trainImages := train.Images
trainLabels := train.Labels

testImages := test.Images
testLabels := test.Labels
```

## Use already downloaded database files
```go
import github.com/r2dtools/mnist/loader

train, test, err := loader.LoadData("/working/directory")
if err != nil {
    panic(err)
}

....
```

## Convert base data structures to gonum Matrix
```go
import (
    github.com/r2dtools/mnist/dense
    github.com/r2dtools/mnist/loader
)

train, test, err := loader.LoadData("")
if err != nil {
    panic(err)
}

images := dense.NewImageDenses(train.Images, 1, 1000) // convert images data to a slice of gonum matrices with dimension 1x1000
labels := dense.NewLabelVecDense(train.Labels)

normalizedImages := dense.NewNormalizedImageDenses(train.Images, 1, 1000, 255) // convert images data to a slice of gonum matrices with dimension 1x1000. Devide all element by 255

....
```
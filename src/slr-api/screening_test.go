package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// TestArticle result for training and verifying our model
type TestArticle struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Include  bool   `json:"include"`
}

type ScreenTestSuite struct {
	suite.Suite
}

func (s *ScreenTestSuite) TestAccuracySmallSet() {
	file, err := ioutil.ReadFile("testdata/article_set_small_logging.json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)

	modelID := uuid.Must(uuid.NewV4())

	training := 1

	trainedInclude := 0
	trainedExclude := 0

	for _, art := range testData {
		switch art.Include {
		case true:
			if trainedInclude < training {
				err := TrainModel(modelID, art.Abstract, art.Title, art.Include)
				s.Require().NoError(err)
				trainedInclude++
			}
		case false:
			err := TrainModel(modelID, art.Abstract, art.Title, art.Include)
			s.Require().NoError(err)
			trainedExclude++
		}
	}

	abstractAccuracy := &AccuracyScore{Total: len(testData)}
	titleAccuracy := &AccuracyScore{Total: len(testData)}

	for _, art := range testData {
		res, err := GetScreeningMediaForProject(modelID, art.Title, art.Abstract)
		s.Require().NoError(err)
		abstractAccuracy.verifyResult(art.Include, res.Abstract.Class)
		titleAccuracy.verifyResult(art.Include, res.Title.Class)

	}

	s.Equal(9, abstractAccuracy.Correct)
	s.Equal(0, abstractAccuracy.Incorrect)
	s.Equal(8, titleAccuracy.Correct)
	s.Equal(1, titleAccuracy.Incorrect)
}

type AccuracyScore struct {
	Total         int
	Correct       int
	Incorrect     int
	FalsePositive int
	FalseNegative int
}

func (s *ScreenTestSuite) TestAccuracyLargeSet() {
	file, err := ioutil.ReadFile("testdata/article_set_large_logging.json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)

	modelID := uuid.Must(uuid.NewV4())

	training := 40
	trainedSet := []int{}
	for i, art := range testData {
		if i > training-1 {
			break
		}
		trainedSet = append(trainedSet, i)
		err := TrainModel(modelID, art.Abstract, art.Title, art.Include)
		s.Require().NoError(err)
	}
	spew.Dump(trainedSet)

	abstractAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}
	titleAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}
	titleAndAbstractAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}

	for i, art := range testData {
		if isInTrainedSet(trainedSet, i) {
			continue
		}
		res, err := GetScreeningMediaForProject(modelID, art.Title, art.Abstract)
		s.Require().NoError(err)
		abstractAccuracy.verifyResult(art.Include, res.Abstract.Class)
		titleAccuracy.verifyResult(art.Include, res.Title.Class)
		titleAndAbstractAccuracy.verifyResult(art.Include, res.AbstractAndTitle.Class)
		if i%50 == 0 {
			log.Infof("Title and Abstract:")
			titleAndAbstractAccuracy.PrintAccuracy()
			log.Infof("Abstract:")
			abstractAccuracy.PrintAccuracy()
			log.Infof("Title:")
			titleAccuracy.PrintAccuracy()
		}
	}
	log.Infof("Title and Abstract:")
	titleAndAbstractAccuracy.PrintAccuracy()
	log.Infof("Abstract:")
	abstractAccuracy.PrintAccuracy()
	log.Infof("Title:")
	titleAccuracy.PrintAccuracy()
}

func isInTrainedSet(trainedSet []int, i int) bool {
	for _, j := range trainedSet {
		if j == i {
			return true
		}
	}
	return false
}

func (s *AccuracyScore) verifyResult(shouldBeIncluded bool, definedClass string) {
	if shouldBeIncluded && definedClass == "Include" {
		s.Correct++
	} else if !shouldBeIncluded && definedClass == "Exclude" {
		s.Correct++
	} else {
		if shouldBeIncluded {
			s.FalseNegative++
		} else {
			s.FalsePositive++
		}
		s.Incorrect++
	}
}

func (s *AccuracyScore) PrintAccuracy() {
	log.Infof("Correct: %d, %f", s.Correct, calcPercentage(s.Correct, s.Total))
	log.Infof("Incorrect: %d, %f", s.Incorrect, calcPercentage(s.Incorrect, s.Total))
	log.Infof("False Positive: %d, %f", s.FalsePositive, calcPercentage(s.FalsePositive, s.Incorrect))
	log.Infof("False Negative: %d, %f", s.FalseNegative, calcPercentage(s.FalseNegative, s.Incorrect))
}

func calcPercentage(x, y int) float64 {
	return float64(x) / float64(y) * 100
}

func TestScreenTestSuite(t *testing.T) {
	test := &ScreenTestSuite{}
	suite.Run(t, test)
}

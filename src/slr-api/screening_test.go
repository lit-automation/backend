package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/gofrs/uuid"
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

func (s *ScreenTestSuite) TestAccuracy() {
	file, err := ioutil.ReadFile("testdata/article_set.json")
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

	correctAbstract := 0
	wrongAbstract := 0

	correctTitle := 0
	wrongTitle := 0

	for _, art := range testData {
		res, err := GetScreeningMediaForProject(modelID, art.Title, art.Abstract)
		s.Require().NoError(err)
		if art.Include && res.Abstract.Class == "Include" {
			correctAbstract++
		} else if !art.Include && res.Abstract.Class == "Exclude" {
			correctAbstract++
		} else {
			wrongAbstract++
		}

		if art.Include && res.Title.Class == "Include" {
			correctTitle++
		} else if !art.Include && res.Title.Class == "Exclude" {
			correctTitle++
		} else {
			wrongTitle++
		}
	}

	s.Equal(9, correctAbstract)
	s.Equal(0, wrongAbstract)
	s.Equal(8, correctTitle)
	s.Equal(1, wrongTitle)
}

func TestScreenTestSuite(t *testing.T) {
	test := &ScreenTestSuite{}
	suite.Run(t, test)
}

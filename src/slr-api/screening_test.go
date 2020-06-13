package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

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

	tfIDFAccuracy := &AccuracyScore{Total: len(testData)}
	tfAccuracy := &AccuracyScore{Total: len(testData)}

	for _, art := range testData {
		res, err := GetScreeningMediaForProject(modelID, art.Title, art.Abstract)
		s.Require().NoError(err)
		tfIDFAccuracy.verifyResult(art.Include, res.Tfidf.Class)
		tfAccuracy.verifyResult(art.Include, res.Tf.Class)
	}

	s.Equal(9, tfIDFAccuracy.Correct)
	s.Equal(0, tfIDFAccuracy.Incorrect)
	s.Equal(9, tfAccuracy.Correct)
	s.Equal(0, tfAccuracy.Incorrect)
}

type AccuracyScore struct {
	Total         int
	Correct       int
	TruePositive  int
	TrueNegative  int
	Incorrect     int
	FalsePositive int
	FalseNegative int
}

type AccuracyScoreFloat struct {
	Total         float64
	Correct       float64
	TruePositive  float64
	TrueNegative  float64
	Incorrect     float64
	FalsePositive float64
	FalseNegative float64
}

func (s *ScreenTestSuite) createBalancedTestData(testData []TestArticle) {
	res := []TestArticle{}
	counterP := 0
	counterN := 0
	for _, a := range testData {
		if a.Include && counterP > 50 {
			continue
		}
		if a.Include {
			counterP++
		}
		if !a.Include && counterN > 50 {
			continue
		}
		if !a.Include {
			counterN++
		}
		res = append(res, a)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(res), func(i, j int) { res[i], res[j] = res[j], res[i] })

	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	file2, err := os.Create("testsetbalanced.json")
	if err != nil {
		panic(err)
	}
	_, err = file2.Write(b)
	if err != nil {
		panic(err)
	}
}

func (s *ScreenTestSuite) TestOutputCreator() {
	file, err := ioutil.ReadFile("article_set_large_loggin_100_unbalancedtfidf.json")
	s.Require().NoError(err)
	testData := make(map[int][]*AccuracyScore)
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)

	for k, v := range testData {
		res := AccuracyScoreFloat{
			Total: float64(v[0].Total),
		}
		for _, x := range v {
			res.FalseNegative += float64(x.FalseNegative)
			res.FalsePositive += float64(x.FalsePositive)
			res.Correct += float64(x.Correct)
			res.Incorrect += float64(x.Incorrect)
			res.TrueNegative += float64(x.TrueNegative)
			res.TruePositive += float64(x.TruePositive)
		}
		res.FalseNegative = res.FalseNegative / float64(len(v))
		res.FalsePositive = res.FalsePositive / float64(len(v))
		res.Correct = res.Correct / float64(len(v))
		res.Incorrect = res.Incorrect / float64(len(v))
		res.TrueNegative = res.TrueNegative / float64(len(v))
		res.TruePositive = res.TruePositive / float64(len(v))

		precision := float64(res.TruePositive) / float64(res.TruePositive+res.FalsePositive)
		recall := float64(res.TruePositive) / float64(res.TruePositive+res.FalseNegative)
		f1 := float64(res.TruePositive*2) / float64(res.TruePositive*2+res.FalsePositive+res.FalseNegative)
		predErr := float64(res.FalsePositive+res.FalseNegative) / float64(res.TruePositive+res.FalsePositive+res.TrueNegative+res.FalseNegative)

		fmt.Printf(`%d,"%f","%f","%f","%f","%f","%f","%f","%f","%f","%f"`, k, f1, precision, recall, predErr, res.Correct, res.Incorrect, res.FalseNegative, res.FalsePositive, res.TrueNegative, res.TruePositive)
	}
}

func (s *ScreenTestSuite) TestAccuracyLargeSet() {
	fileName := "article_set_large_loggin_100_unbalanced"
	file, err := ioutil.ReadFile("testdata/" + fileName + ".json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)
	trainingAmounts := []int{2, 6, 10, 20, 30, 50, 90}
	resultMapTF := make(map[int][]*AccuracyScore)
	resultMapTFIDF := make(map[int][]*AccuracyScore)
	for _, x := range trainingAmounts {
		// Add crossvalidation of 5 runs
		runs := 5
		remainder := len(testData) - x
		incr := remainder / runs

		for j := 0; j < runs; j++ {
			startI := j * incr

			training := x
			modelID := uuid.Must(uuid.NewV4())
			log.Infof("Processing set of size: %d, model: %s", training, modelID.String())

			trainedSet := []int{}
			counterP := 0
			counterN := 0
			for i, art := range testData {
				if i < startI {
					continue
				}
				if len(trainedSet) > training-1 {
					break
				}
				if x == 2 {
					if art.Include && counterP >= 1 {
						continue
					}
					if !art.Include && counterN >= 1 {
						continue
					}
					if art.Include {
						counterP++
					} else {
						counterN++
					}
				}
				trainedSet = append(trainedSet, i)
				err := TrainModel(modelID, art.Abstract, art.Title, art.Include)
				s.Require().NoError(err)
			}
			log.Infof("Predicting")
			s.predict(testData, trainedSet, modelID, resultMapTF, resultMapTFIDF)
		}
	}
	s.writeMapToDisk(fileName+"tf.json", resultMapTF)
	s.writeMapToDisk(fileName+"tfidf.json", resultMapTFIDF)
}

func (s *ScreenTestSuite) writeMapToDisk(name string, res map[int][]*AccuracyScore) {
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	file2, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	_, err = file2.Write(b)
	if err != nil {
		panic(err)
	}
}

func (s *ScreenTestSuite) predict(testData []TestArticle, trainedSet []int, modelID uuid.UUID, resultMapTF, resultMapTFIDF map[int][]*AccuracyScore) {
	tfIDFAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}
	termFrequencyAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}
	processed := 0
	for i, art := range testData {
		if isInTrainedSet(trainedSet, i) {
			continue
		}
		processed++
		res, err := GetScreeningMediaForProject(modelID, art.Title, art.Abstract)
		s.Require().NoError(err)
		tfIDFAccuracy.verifyResult(art.Include, res.Tfidf.Class)
		termFrequencyAccuracy.verifyResult(art.Include, res.Tf.Class)
	}
	log.Infof("TF:")
	termFrequencyAccuracy.PrintAccuracy()

	resultMapTF[len(trainedSet)] = append(resultMapTF[len(trainedSet)], termFrequencyAccuracy)
	log.Infof("TFIDF:")
	tfIDFAccuracy.PrintAccuracy()
	resultMapTFIDF[len(trainedSet)] = append(resultMapTFIDF[len(trainedSet)], tfIDFAccuracy)
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
		s.TruePositive++
	} else if !shouldBeIncluded && definedClass == "Exclude" {
		s.Correct++
		s.TrueNegative++
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
	log.Infof("True Positive: %d, %f", s.TruePositive, calcPercentage(s.TruePositive, s.Correct))
	log.Infof("True Negative: %d, %f", s.TrueNegative, calcPercentage(s.TrueNegative, s.Correct))
	log.Infof("False Positive: %d, %f", s.FalsePositive, calcPercentage(s.FalsePositive, s.Incorrect))
	log.Infof("False Negative: %d, %f", s.FalseNegative, calcPercentage(s.FalseNegative, s.Incorrect))
	log.Infof("Precision: %f", float64(s.TruePositive)/float64(s.TruePositive+s.FalsePositive))
	log.Infof("Recall: %f", float64(s.TruePositive)/float64(s.TruePositive+s.FalseNegative))
	log.Infof("F1: %f", float64(s.TruePositive*2)/float64(s.TruePositive*2+s.FalsePositive+s.FalseNegative))
	log.Infof("Error: %f", float64(s.FalsePositive+s.FalseNegative)/float64(s.TruePositive+s.FalsePositive+s.TrueNegative+s.FalseNegative))
}

func calcPercentage(x, y int) float64 {
	return float64(x) / float64(y) * 100
}

func TestScreenTestSuite(t *testing.T) {
	test := &ScreenTestSuite{}
	suite.Run(t, test)
}

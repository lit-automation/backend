package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/wimspaargaren/slr-automation/src/packages/database"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
)

const (
	currentSet = "article_set_large_rl_personalization_100_unbalanced"
)

var (
	screeningTestSetSizes = []int{6, 8, 10, 20, 30, 40, 50, 75, 90}
)

// TestArticle result for training and verifying our model
type TestArticle struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Include  bool   `json:"include"`
	Article  *models.Article
}

type ScreenTestSuite struct {
	suite.Suite

	DB *gorm.DB

	transact *gorm.DB

	DBManager *DBManager
}

func (s *ScreenTestSuite) SetupTest() {
	s.transact = s.DB.Begin()
	s.DBManager = NewDBManager(s.transact)
}

func (s *ScreenTestSuite) TearDownTest() {
	s.transact.Rollback()
}

func (s *ScreenTestSuite) TestAccuracySmallSet() {
	file, err := ioutil.ReadFile("testdata/article_set_small_logging.json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)

	testData = s.prepareTestData(testData)
	modelID := uuid.Must(uuid.NewV4())

	err = TrainModel(modelID, testData[0].Article, true, testData[0].Include)
	s.Require().NoError(err)

	err = TrainModel(modelID, testData[2].Article, true, testData[2].Include)
	s.Require().NoError(err)

	tfIDFAccuracy := &AccuracyScore{Total: len(testData)}
	tfAccuracy := &AccuracyScore{Total: len(testData)}

	for _, art := range testData {
		res, err := GetScreeningMediaForProject(modelID, art.Article, true)
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
	file, err := ioutil.ReadFile(currentSet + "tfidf.json")
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
		if res.TruePositive+res.FalsePositive == 0 {
			precision = 1
		}
		recall := float64(res.TruePositive) / float64(res.TruePositive+res.FalseNegative)
		if res.TruePositive+res.FalseNegative == 0 {
			recall = 1
		}
		f1 := float64(res.TruePositive*2) / float64(res.TruePositive*2+res.FalsePositive+res.FalseNegative)
		if float64(res.TruePositive*2+res.FalsePositive+res.FalseNegative) == 0 {
			f1 = 1
		}
		predErr := float64(res.FalsePositive+res.FalseNegative) / float64(res.TruePositive+res.FalsePositive+res.TrueNegative+res.FalseNegative)
		fmt.Printf(`%d,"%f","%f","%f","%f","%f","%f","%f","%f","%f","%f"`, k, f1, precision, recall, predErr, res.Correct, res.Incorrect, res.FalseNegative, res.FalsePositive, res.TrueNegative, res.TruePositive)
		fmt.Println()
	}
}

func (s *ScreenTestSuite) prepareTestData(testData []TestArticle) []TestArticle {
	fmt.Println("Preparing data")
	for i, art := range testData {
		dbArtcle := &models.Article{
			Title:       art.Title,
			Abstract:    art.Abstract,
			Keywords:    []byte("{}"),
			Metadata:    []byte("{}"),
			DocAbstract: []byte("{}"),
			DocFullText: []byte("{}"),
			CitedBy:     []byte("[]"),
		}
		err := dbArtcle.SetAbstractDoc()
		s.Require().NoError(err)
		err = s.DBManager.ArticleDB.Add(context.Background(), dbArtcle)
		s.Require().NoError(err)
		testData[i].Article = dbArtcle
	}
	fmt.Println("Done preparing")
	return testData
}

func (s *ScreenTestSuite) TestActiveLearningPerDocument() {
	toTrainSet := []int{1}

	fileName := currentSet
	file, err := ioutil.ReadFile("testdata/" + fileName + ".json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)
	testData = s.prepareTestData(testData)
	modelID := uuid.Must(uuid.NewV4())
	trainedSet := []int{}
	resultMapTF := make(map[int][]*AccuracyScore)
	resultMapTFIDF := make(map[int][]*AccuracyScore)
	for len(trainedSet) < 90 {
		for i := range toTrainSet {
			index := toTrainSet[i]
			trainedSet = append(trainedSet, index)
			err := TrainModel(modelID, testData[index].Article, true, testData[index].Include)
			s.Require().NoError(err)
		}
		fmt.Println(len(trainedSet))
		result := s.predictActiveSet(testData, trainedSet, modelID, resultMapTF, resultMapTFIDF)

		sort.Slice(result, func(i, j int) bool {
			return result[i].Confidence < result[j].Confidence
		})
		toTrainSet = []int{}
		toTrainSet = append(toTrainSet, result[0].Index)
	}
	s.writeMapToDisk(fileName+"tf.json", resultMapTF)
	s.writeMapToDisk(fileName+"tfidf.json", resultMapTFIDF)
}

func (s *ScreenTestSuite) TestActiveLearning() {
	fileName := currentSet
	file, err := ioutil.ReadFile("testdata/" + fileName + ".json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)
	testData = s.prepareTestData(testData)
	trainingAmounts := append(screeningTestSetSizes, 0)
	modelID := uuid.Must(uuid.NewV4())
	trainedSet := []int{}
	// Init for test
	toTrainSet := []int{7, 6, 4, 69}
	resultMapTF := make(map[int][]*AccuracyScore)
	resultMapTFIDF := make(map[int][]*AccuracyScore)
	for _, x := range trainingAmounts {
		for i := range toTrainSet {
			index := toTrainSet[i]
			trainedSet = append(trainedSet, index)
			err := TrainModel(modelID, testData[index].Article, true, testData[index].Include)
			s.Require().NoError(err)
		}
		fmt.Println(len(trainedSet))
		result := s.predictActiveSet(testData, trainedSet, modelID, resultMapTF, resultMapTFIDF)

		sort.Slice(result, func(i, j int) bool {
			return result[i].Confidence < result[j].Confidence
		})
		toTrainSet = []int{}
		for _, y := range result {
			if y.Confidence == 0.5 {
				fmt.Println(x, y.Confidence)
			}
		}
		for j := 0; j < x-len(trainedSet); j++ {
			toTrainSet = append(toTrainSet, result[j].Index)
		}
	}
	s.writeMapToDisk(fileName+"tf.json", resultMapTF)
	s.writeMapToDisk(fileName+"tfidf.json", resultMapTFIDF)
}

type ActiveLearningPredictor struct {
	Index      int
	Confidence float64
}

func (s *ScreenTestSuite) predictActiveSet(testData []TestArticle, trainedSet []int, modelID uuid.UUID, resultMapTF, resultMapTFIDF map[int][]*AccuracyScore) []ActiveLearningPredictor {
	activeLearningPredictor := []ActiveLearningPredictor{}
	tfIDFAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}
	termFrequencyAccuracy := &AccuracyScore{Total: len(testData) - len(trainedSet)}
	processed := 0
	for i, art := range testData {
		if isInTrainedSet(trainedSet, i) {
			continue
		}
		processed++
		res, err := GetScreeningMediaForProject(modelID, art.Article, true)
		s.Require().NoError(err)
		activeLearningPredictor = append(activeLearningPredictor, ActiveLearningPredictor{
			Confidence: res.Tfidf.Confidence,
			Index:      i,
		})
		tfIDFAccuracy.verifyResult(art.Include, res.Tfidf.Class)
		termFrequencyAccuracy.verifyResult(art.Include, res.Tf.Class)
	}

	resultMapTF[len(trainedSet)] = append(resultMapTF[len(trainedSet)], termFrequencyAccuracy)
	log.Infof("TFIDF:")
	tfIDFAccuracy.PrintAccuracy()
	resultMapTFIDF[len(trainedSet)] = append(resultMapTFIDF[len(trainedSet)], tfIDFAccuracy)
	return activeLearningPredictor
}

func (s *ScreenTestSuite) TestAccuracyLargeSet() {
	fileName := currentSet
	file, err := ioutil.ReadFile("testdata/" + fileName + ".json")
	s.Require().NoError(err)
	testData := []TestArticle{}
	err = json.Unmarshal([]byte(file), &testData)
	s.Require().NoError(err)

	testData = s.prepareTestData(testData)

	trainingAmounts := append([]int{2}, screeningTestSetSizes...)
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
				err := TrainModel(modelID, art.Article, true, art.Include)
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
		res, err := GetScreeningMediaForProject(modelID, art.Article, true)
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
	m := []interface{}{
		models.Article{},
		models.Project{},
		models.User{},
	}
	db, err := database.ConnectTest("slr", m)
	if err != nil {
		t.FailNow()
	}
	test := &ScreenTestSuite{
		DB: db,
	}
	suite.Run(t, test)
}

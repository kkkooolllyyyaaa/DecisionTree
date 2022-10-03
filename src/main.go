package main

import (
	"DecisionTree/src/dataWorker"
	"DecisionTree/src/decisionTree"
	"fmt"
	"os"
)

func main() {

	studsTrain, headerTrain, err := dataWorker.Parse("resources/data_train.txt", ";")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	dt := decisionTree.DecisionTree{Header: headerTrain, Students: dataWorker.Students{Studs: studsTrain}}
	root := dt.Init(30)

	studsTest, _, _ := dataWorker.Parse("resources/data_test.txt", ";")
	rows := float64(len(studsTest))
	p := 0.0
	tp := 0.0
	fp := 0.0
	fn := 0.0

	toWrite := make([][]string, len(studsTest))
	for i, value := range studsTest {
		actual := value.Grade
		predicted, _ := dt.Predict(root, value)

		if good(actual) == good(predicted) {
			p += 1.0
		}
		if good(actual) && good(predicted) {
			tp += 1.0
		}
		if !good(actual) && good(predicted) {
			fp += 1.0
		}
		if good(actual) && !good(predicted) {
			fn += 1.0
		}
		toWrite[i] = make([]string, 2)
		toWrite[i][0] = goodString(actual)
		toWrite[i][1] = goodString(predicted)
	}
	err = dataWorker.Write("resources/result.txt", ",", toWrite, 2)
	if err != nil {
		err.Error()
	}

	fmt.Println("Accuracy:", p/rows)
	fmt.Println("Precision:", tp/(tp+fp))
	fmt.Println("Recall:", tp/(tp+fn))
}

func good(grade string) bool {
	return grade < "4"
}

func goodString(grade string) string {
	if good(grade) {
		return "1"
	} else {
		return "0"
	}
}

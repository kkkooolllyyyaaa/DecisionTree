package decisionTree

import (
	"DecisionTree/src/dataWorker"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Node struct {
	children map[string]*Node
	idx      int
	class    string
}

func NewLeaf(class string) *Node {
	return &Node{class: class}
}

func NewNode() *Node {
	node := Node{class: "NO_CLASS", children: make(map[string]*Node)}
	return &node
}

type DecisionTree struct {
	Header   dataWorker.Header
	Students dataWorker.Students
	root     *Node
}

var selectedAttrs []int
var n int

type IdxToGainRation struct {
	idx       int
	gainRatio float64
}

var maxDepth int

func (dt DecisionTree) Init(maxDepthInput int) *Node {
	n = len(dt.Students.Studs[0].Attributes)
	rand.Seed(time.Now().UnixNano())
	selectedAttrs = selectAttributes(0, n, int(math.Sqrt(float64(n))))
	maxDepth = maxDepthInput
	dt.root = NewNode()
	dt.root.children = buildNode(dt.Students, dt.root, 0)
	return dt.root
}

func (dt DecisionTree) PrintSelected() {
	fmt.Println("Выбранные аттрибуты:")
	for _, val := range selectedAttrs {
		fmt.Print(dt.Header.AttributesNames[val], " ")
	}
	fmt.Println()
}

func buildNode(students dataWorker.Students, parent *Node, depth int) (children map[string]*Node) {
	children = make(map[string]*Node)
	gradesList := students.Grades()
	grades := extractDistinct(gradesList)

	infoT := info(gradesList, grades)
	max := 0.0
	maxIdx := 0
	for _, idx := range selectedAttrs {
		infoX := infoX(students, idx)
		gainRatio := (infoT - infoX) / splitInfoX(students, idx)
		if gainRatio > max {
			max = gainRatio
			maxIdx = idx
		}
	}
	parent.idx = maxIdx

	columnValues := students.ColumnAsSlice(maxIdx)
	columnDistinct := extractDistinct(columnValues)
	for _, value := range columnDistinct {
		selectedStudents := students.SelectWhereEq(maxIdx, value)
		selectedStudentsGrades := extractDistinct(selectedStudents.Grades())

		if len(selectedStudentsGrades) == 1 {
			children[value] = NewLeaf(selectedStudentsGrades[0])
		} else if depth > maxDepth {
			maxGrades := 0
			grade := ""
			for _, val := range selectedStudentsGrades {
				freq := frequencyOfValue(selectedStudents.Grades(), val)
				if freq > maxGrades {
					maxGrades = freq
					grade = val
				}
			}
			children[value] = NewLeaf(grade)
		} else {
			node := NewNode()
			node.children = buildNode(selectedStudents, node, depth+1)
			children[value] = node
		}
	}
	return children
}

func (dt DecisionTree) traversePath(student dataWorker.Student) (path []string) {
	var current = dt.root
	for current.children != nil {
		idx := current.idx
		value := student.Attributes[idx]
		path = append(path, value)
		current = current.children[value]
	}
	return
}

func (dt DecisionTree) Predict(root *Node, student dataWorker.Student) (class string, err error) {
	var current = root
	for current.children != nil {
		idx := current.idx
		value := student.Attributes[idx]
		if _, ok := current.children[value]; ok {
			current = current.children[value]
		} else {
			min := 2147483647.0
			var foundKey string
			for k := range current.children {
				kInt, _ := strconv.Atoi(k)
				valueInt, _ := strconv.Atoi(value)
				diff := math.Abs(float64(kInt - valueInt))
				if diff < min {
					foundKey = k
					min = diff
				}
			}
			current = current.children[foundKey]
		}
	}
	return current.class, nil
}

func splitInfoX(students dataWorker.Students, idx int) (result float64) {
	valuesList := students.ColumnAsSlice(idx)
	values := extractDistinct(valuesList)

	for _, value := range values {
		selectedStuds := students.SelectWhereEq(idx, value)
		div := float64(len(selectedStuds.Studs)) / float64(len(students.Studs))
		result -= div * math.Log2(div)
	}
	return
}

func infoX(students dataWorker.Students, idx int) (result float64) {
	valuesList := students.ColumnAsSlice(idx)
	values := extractDistinct(valuesList)

	for _, value := range values {
		selectedStuds := students.SelectWhereEq(idx, value)
		selectedStudsGrades := selectedStuds.Grades()
		selectedStudsGradesDistinct := extractDistinct(selectedStudsGrades)

		result += float64(len(selectedStuds.Studs)) / float64(len(students.Studs)) *
			info(selectedStudsGrades, selectedStudsGradesDistinct)
	}
	return
}

func info(valuesList []string, valuesDistinct []string) (acc float64) {
	for _, value := range valuesDistinct {
		freq := frequencyOfValue(valuesList, value)
		div := float64(freq) / float64(len(valuesList))
		acc -= div * (math.Log2(div))
	}
	return
}

func frequencyOfValue(values []string, target string) (count int) {
	for _, value := range values {
		if value == target {
			count++
		}
	}
	return
}

type void struct{}

var member void

func extractDistinct(of []string) (result []string) {
	set := make(map[string]void)
	for _, value := range of {
		set[value] = member
	}

	for k := range set {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

func selectAttributes(min, max, n int) (result []int) {
	set := make(map[int]void)
	for len(set) < n {
		set[rand.Intn(max-min)+min] = member
	}

	for k := range set {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return
}

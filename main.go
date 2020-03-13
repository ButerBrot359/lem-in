package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	Nodes           = map[string][]string{}
	Pathes          [][]string
	unRelatedRoutes [][]int
	count           int
)

type Ant struct {
	id           int
	path         []string
	currentPlace int
	step         int
}

func main() {

	arg := os.Args[1:]
	if len(arg) != 1 {
		fmt.Println("Введите имя файла")
		return
	}
	filename := arg[0]
	inputStringArrayFromFile := readFile(filename)
	antCount, err := getAntCount(inputStringArrayFromFile)
	if !err {
		fmt.Println("bad format")
		return
	}

	nodeNames, err := getNodeNames(inputStringArrayFromFile)
	if !err {
		fmt.Println("bad format when get")
		return
	}

	nodeCoordinates, err := getNodeCoordinates(inputStringArrayFromFile, nodeNames)
	if !err {
		fmt.Println("bad format when getCoordinates")
		return
	}
	// fmt.Println("NodeCoordinatess =", nodeCoordinates)
	nodeConnections, err := getNodeConnections(inputStringArrayFromFile, nodeNames)
	if !err {
		fmt.Println("bad format when find connections")
		return
	}
	// fmt.Println("NodeConnections =", nodeConnections)
	fillMap(nodeNames, nodeConnections)
	path := []string{nodeNames[0]}
	end := nodeNames[len(nodeNames)-1]
	findPath(path[0], end, path)

	picture(nodeCoordinates, nodeNames)
	unrelatedArrays := findArrayOfNotRealtedPaths()
	heights := findCountOfAuntMovesForEveryArray(unrelatedArrays, antCount)
	index := getMinIndex(heights)

	var array []int
	for i := range Pathes {
		findUnrelated(array, i)
	}
	heights = findCountOfAuntMovesForEveryArray(unRelatedRoutes, antCount)
	index = getMinIndex(heights)

	sendAnts(index, antCount, nodeNames)

}

//запускаем муравьев---------------------------------------------------------------------
func sendAnts(routeIndex int, antCount int, nodeNames []string) {
	end := nodeNames[len(nodeNames)-1]
	route := unRelatedRoutes[routeIndex]
	//надо сделать сортировку по длине
	for i := 0; i < len(route)-1; i++ {
		for j := i + 1; j < len(route); j++ {
			// fmt.Println("here", height(route[i]), height(route[j]))
			if height(route[i]) > height(route[j]) {
				n := route[i]
				route[i] = route[j]
				route[j] = n
			}
		}
	}
	// cсортировку закончил

	var Routes [][]string
	for i := range route {
		var path []string
		path = append(Pathes[route[i]], end)
		Routes = append(Routes, path)
	}
	// находим массив сколько надо муравьев запустить на каждый путь
	var heights []int
	var heightx, newHeight []int
	flagForZero := false
	for _, f := range route {
		// fmt.Println(f, height(f))
		heightx = append(heightx, height(f))
		newHeight = append(newHeight, height(f))
		if height(f) == 0 {
			flagForZero = true
		}
	}
	if flagForZero { // при случае нулевого прохода
		// fmt.Println("Zerro is here")
		ants := make([]Ant, antCount)
		var a []string
		a = append(a, end)
		for i := range ants {
			ants = append(ants, Ant{id: i, path: a, currentPlace: 0, step: 0})
		}
		sx := make([][]string, len(ants))

		for i := range ants {
			Ant1 := ants[i]
			for Ant1.currentPlace < len(Ant1.path) {
				s := "L" + strconv.Itoa(Ant1.id+1) + "-" + Ant1.path[Ant1.currentPlace]
				if len(s) != 0 {
					sx[i] = append(sx[i], s)
				}
				Ant1.currentPlace++
			}
		}

		sf := make([]string, 1)
		for i, v := range sx {
			for j := range v {
				sf[j+ants[i].step] += v[j] + " "
			}
		}
		for i, v := range sf {
			fmt.Println(i, v)
		}
		return
	}

	n := antCount
	for n > 0 {
		index := getMinIndex(newHeight)
		newHeight[index]++
		n--
	}
	max := newHeight[0]
	for i := range heightx {
		if newHeight[i] != heightx[i] && newHeight[i] > max {
			max = newHeight[i]
		}
		heights = append(heights, newHeight[i]-heightx[i])
	}

	var ants []Ant
	sum := 0
	for _, v := range heights {
		sum += v
	}
	count := 0
	steps := 0
	for sum > 0 {
		for i, v := range heights {
			if v != 0 {
				heights[i]--
				ants = append(ants, Ant{id: count, path: Routes[i], step: steps})
				count++
			}
		}
		steps++
		sum = 0
		for _, v := range heights {
			sum += v
		}
		// fmt.Println(heights, sum)
	}
	// for _, v := range ants {
	// 	fmt.Println(v)
	// }
	sx := make([][]string, len(ants))

	for i := range ants {
		Ant1 := ants[i]
		for Ant1.currentPlace < len(Ant1.path) {
			s := "L" + strconv.Itoa(Ant1.id+1) + "-" + Ant1.path[Ant1.currentPlace]
			sx[i] = append(sx[i], s)
			Ant1.currentPlace++
		}
	}

	sf := make([]string, max)
	for i, v := range sx {
		for j := range v {
			sf[j+ants[i].step] += v[j] + " "
		}
	}
	for _, v := range sf {
		fmt.Println(v)
	}

}

func findCountOfAuntMovesForEveryArray(unRelatedRoutes [][]int, antCount int) []int {
	var heights []int
	flagForZero := false
	for _, v := range unRelatedRoutes {

		var heightx, newHeight []int

		for _, f := range v {
			heightx = append(heightx, height(f))
			newHeight = append(newHeight, height(f))
			if height(f) == 0 {
				flagForZero = true
			}
		}
		if flagForZero {
			heights = append(heights, 1)
			continue
		}

		n := antCount
		for n > 0 {
			index := getMinIndex(newHeight)
			newHeight[index]++
			n--
		}
		max := 0
		for i := range heightx {
			if heightx[i] != newHeight[i] {
				if max < newHeight[i] {
					max = newHeight[i]
				}
			}
		}

		heights = append(heights, max)
	}

	return heights
}

func getMinIndex(s []int) int {
	min := s[0]
	mini := 0
	for i, v := range s {
		if v < min {
			min = v
			mini = i
		}
	}
	return mini
}

func getMin(s []int) int {
	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min
}

func height(index int) int {
	return len(Pathes[index])
}

func findArrayOfNotRealtedPaths() [][]int {
	arrayOfNotRelatedPathes := make([][][]string, len(Pathes))
	for i := range arrayOfNotRelatedPathes {
		arrayOfNotRelatedPathes[i] = make([][]string, 0)
	}

	arrayOfNotRelatedPathesNames := make([][]int, len(Pathes))
	for i := range arrayOfNotRelatedPathes {
		arrayOfNotRelatedPathes[i] = append(arrayOfNotRelatedPathes[i], Pathes[i])
		arrayOfNotRelatedPathesNames[i] = append(arrayOfNotRelatedPathesNames[i], i)
		for j := i; j < len(Pathes); j++ {
			flag := true
			v := Pathes[j]
			for _, f := range arrayOfNotRelatedPathes[i] {
				if !isPathesUnique(v, f) {
					flag = false
				}
			}
			if flag && len(v) != 0 {
				arrayOfNotRelatedPathes[i] = append(arrayOfNotRelatedPathes[i], v)
				arrayOfNotRelatedPathesNames[i] = append(arrayOfNotRelatedPathesNames[i], j)
			}
		}
	}

	var new [][]int
	new = append(new, arrayOfNotRelatedPathesNames[0])
	for _, v := range arrayOfNotRelatedPathesNames {
		flag := true
		for _, f := range new {
			count := 0
			for i := range f {
				for j := range v {
					if f[i] == v[j] {
						count++
					}
				}
			}
			if count == len(f) {
				flag = false
			}
		}
		if flag {
			new = append(new, v)
		}

	}
	for i, v := range Pathes {
		if len(v) == 0 {
			var a []int
			a = append(a, i)
			fmt.Println(a)
			new = append(new, a)
		}
	}

	//arrayOfNotRelatedPathesNames = new
	return arrayOfNotRelatedPathesNames
}

func findUnrelated(array []int, index int) {
	if len(array) == 0 {
		array = append(array, index)
	}
	flag := true
	for i := index + 1; i < len(Pathes); i++ {
		flag = true
		for _, v := range array {
			if !isPathesUniqueByIndex(i, v) {
				flag = false
			}
		}
		if flag {
			copy := array
			copy = append(copy, i)
			// fmt.Println(copy)
			findUnrelated(copy, i)
			continue
		}
	}
	flag = false
	for _, v := range unRelatedRoutes {
		if isSubclasses(v, array) {
			flag = true
		}
	}
	if !flag || len(unRelatedRoutes) == 0 {
		unRelatedRoutes = append(unRelatedRoutes, array)
	}

}

func isSubclasses(s1 []int, s2 []int) bool {
	min := len(s1)
	if len(s2) < min {
		min = len(s2)
	}
	count := 0
	for i := range s1 {
		for j := range s2 {
			if s1[i] == s2[j] {
				count++
			}
		}
	}
	if count == min {
		return true
	}
	return false
}

func isPathesUnique(Path1 []string, Path2 []string) bool {
	for _, v := range Path1 {
		for _, f := range Path2 {
			if v == f {
				return false
			}
		}
	}
	return true
}

func isPathesUniqueByIndex(Index1 int, Index2 int) bool {
	for _, v := range Pathes[Index1] {
		for _, f := range Pathes[Index2] {
			if v == f {
				return false
			}
		}
	}
	return true
}

func correctPathes(end string) {
	for _, pps := range Pathes {
		pps[len(pps)-1] = end
		// fmt.Println(pps[1:])
		// pps = append(pps[:0], pps[1:]...)
	}
}

func picture(nodeCoordinates [][]int, nodeNames []string) {
	maxX := nodeCoordinates[0][0]
	maxY := nodeCoordinates[0][1]
	for _, v := range nodeCoordinates {
		if maxX < v[0] {
			maxX = v[0]
		}
		if maxY < v[1] {
			maxY = v[1]
		}
	}
	sTop := "\\"
	sBottom := "/"
	sBody := "|"
	for i := 0; i < (maxX+5)*3; i++ {
		sTop += "-"
		sBody += " "
		sBottom += "-"
	}
	sBody += "|"
	sTop += "/"
	sBottom += "\\"
	sFinal := sBottom + "\n"
	for i := 0; i < (maxY+1)*3; i++ {
		sFinal += sBody + "\n"
	}
	sFinal += sTop + "\n"
	// fmt.Print(sFinal)
	putNodesOnCoordinates(sFinal, nodeCoordinates, nodeNames)

}

func putNodesOnCoordinates(input string, nodeCoordinates [][]int, nodeNames []string) {
	inputSplit := strings.Split(input, "\n")
	for i := range nodeCoordinates {
		x := nodeCoordinates[i][0]
		y := nodeCoordinates[i][1]
		inputSplit = putNodeNamesOnCoordinate(inputSplit, x, y, nodeNames[i])
	}
	sFinals := ""
	for _, v := range inputSplit {
		sFinals += v + "\n"
	}

}

func putNodeNamesOnCoordinate(inputSplit []string, x int, y int, name string) []string {
	x = (x + 1) * 3
	y = (y + 1) * 3
	for i := range inputSplit {
		if i == y {
			inputSplit[i] = inputSplit[i][:x] + "[" + name + "]" + inputSplit[i][x+len(name)+2:]
		}
	}
	return inputSplit
}

func fillMap(nodeNames []string, nodeConnections [][]string) {
	for i := range nodeNames {
		Nodes[nodeNames[i]] = nodeConnections[i]
	}
}

func findPath(current string, end string, path []string) {
	if current == end {
		count++
		copy := make([]string, len(path)-2)
		for i := 0; i < len(path)-2; i++ {
			copy[i] = path[i+1]
		}
		Pathes = append(Pathes, copy)
		return
	}
	for _, v := range Nodes[current] {
		flag := true
		for _, f := range path {
			if v == f {
				flag = false
			}
		}
		if flag {
			copy := path
			copy = append(copy, v)
			findPath(v, end, copy)
		}

	}

}

func getNodeConnections(s []string, nodeNames []string) ([][]string, bool) {
	nodeConnections := make([][]string, len(nodeNames))
	for i := range nodeConnections {
		nodeConnections[i] = make([]string, 0)
	}
	for i := range s {
		pps := strings.Split(s[i], "-")
		if len(pps) == 2 {
			// fmt.Println(pps)
			for j := 0; j < len(nodeNames); j++ {
				if nodeNames[j] == pps[0] {
					nodeConnections[j] = append(nodeConnections[j], pps[1])
				}
				if nodeNames[j] == pps[1] {
					nodeConnections[j] = append(nodeConnections[j], pps[0])
				}
			}
		}

	}
	return nodeConnections, true
}

func readFile(name string) []string {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	result := strings.Split(string(data), "\n")
	return result
}

func getAntCount(s []string) (int, bool) {
	ns := s[0]
	antCount, err := strconv.Atoi(ns)
	if err != nil {
		return 0, false
	}
	return antCount, true
}

func getNodeNames(s []string) ([]string, bool) {
	var nodeNames []string
	var start, finish int
	for i, pp := range s {
		if pp == "##start" {
			start = i + 1
		}
		if pp == "##end" {
			finish = i + 1
		}
	}
	if start == 0 || finish == 0 {
		return nodeNames, false
	}
	for i := range s {
		pps := strings.Split(s[i], " ")
		if len(pps) == 3 && i != finish && i != start {
			nodeNames = append(nodeNames, pps[0])
		}
	}
	pps := strings.Split(s[finish], " ")
	if len(pps) == 3 {
		nodeNames = append(nodeNames, pps[0])
	} else {
		return nodeNames, false
	}

	copy := []string{}
	pps = strings.Split(s[start], " ")
	if len(pps) == 3 {
		copy = append(copy, pps[0])
	} else {
		return nodeNames, false
	}
	for _, v := range nodeNames {
		copy = append(copy, v)
	}
	nodeNames = copy

	return nodeNames, true
}

func getNodeCoordinates(s []string, nodeNames []string) ([][]int, bool) {
	nodeCoordinates := make([][]int, len(nodeNames))
	for i := range nodeCoordinates {
		nodeCoordinates[i] = make([]int, 0)
	}
	for _, v := range s {
		pps := strings.Split(v, " ")
		for i, f := range nodeNames {

			if pps[0] == f {
				// nodeCoordinates[i] = append(nodeCoordinates[i], pps[0])
				if len(pps) == 3 {
					x, err1 := strconv.Atoi(pps[1])
					y, err2 := strconv.Atoi(pps[2])
					if err1 != nil || err2 != nil {
						return nodeCoordinates, false
					}
					nodeCoordinates[i] = append(nodeCoordinates[i], x)
					nodeCoordinates[i] = append(nodeCoordinates[i], y)
				}
			}
		}
	}
	return nodeCoordinates, true
}

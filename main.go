package main

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	set "github.com/deckarep/golang-set"
	"github.com/llgcode/draw2d/draw2dimg"
)

type LRTuple struct {
	dir  float64
	dist int
}

type Tuple struct {
	rad float64
	gx  int
	gy  int
}

var ArrInput []LRTuple
var xyInput []Tuple
var iInput []Tuple
var xyOutput []interface{}

func printErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

//parse and read textfile
func readFile() {
	file, err := os.Open("input.txt")
	printErr(err)
	defer file.Close()

	reader := csv.NewReader(file)
	result, err := reader.Read()
	printErr(err)

	log.Println(len(result))

	for z := range result {
		trimResult := strings.Trim(result[z], " ")
		regStr := regexp.MustCompile("(\\D+)(\\d+)")
		resultz := regStr.FindAllStringSubmatch(trimResult, -1)

		headResult := resultz[0]

		convDist, err := strconv.Atoi(headResult[2])
		printErr(err)
		radian := 0.00
		if headResult[1] == "R" {
			radian = (math.Pi / 180) * 90
		} else {
			radian = (math.Pi / 180) * -90
		}

		ArrInput = append(ArrInput, LRTuple{radian, convDist})

	}
	total := 0
	for i := 0; i < len(ArrInput); i++ {
		rads := strconv.FormatFloat(ArrInput[i].dir, 'E', -1, 64)
		dist := strconv.Itoa(ArrInput[i].dist)

		total = total + ArrInput[i].dist
		log.Println("readFile: " + strconv.Itoa(i) + ", " + rads + ", " + dist)

	}
	log.Println(strconv.Itoa(total))

}

//calculate headings and steps
func parseHeadings() {

	first := Tuple{0.00, 0, 0}
	xyInput = append(xyInput, first)

	for bb := 0; bb < len(ArrInput); bb++ {

		rads := xyInput[bb].rad + ArrInput[bb].dir
		gx := xyInput[bb].gx + (int(math.Sin(rads)) * ArrInput[bb].dist)
		gy := xyInput[bb].gy + (int(math.Cos(rads)) * ArrInput[bb].dist)

		next := Tuple{rads, gx, gy}
		xyInput = append(xyInput, next)
	}

	for k := 0; k < len(xyInput); k++ {
		rads := strconv.FormatFloat(xyInput[k].rad, 'E', 4, 64)
		gx := strconv.Itoa(xyInput[k].gx)
		gy := strconv.Itoa(xyInput[k].gy)
		fmt.Printf("%s, %s, %s\n", rads, gx, gy)
	}

}

//quick drawing to check for intersection
func drawPath() {
	dest := image.NewRGBA(image.Rect(0, 0, 1200, 1200.0))
	gc := draw2dimg.NewGraphicContext(dest)
	gc.SetStrokeColor(color.Black)
	gc.SetLineWidth(1)

	gc.MoveTo(0, 0)
	for i := 0; i < len(xyInput); i++ {
		gc.LineTo(float64(xyInput[i].gx), float64(xyInput[i].gy))

	}

	gc.Close()
	gc.FillStroke()

	draw2dimg.SaveToPngFile("pathline.png", dest)
}

//hashing function
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//calculate all plotpoints for Santa's path
func iIntersection() {

	first := Tuple{0.00, 0, 0}
	iInput = append(iInput, first)
	count := 0
	for k := 0; k < len(ArrInput); k++ {
		for jj := 0; jj < ArrInput[k].dist; jj++ {
			rads := xyInput[k+1].rad
			gx := iInput[count].gx + (int(math.Sin(rads)) * 1)
			gy := iInput[count].gy + (int(math.Cos(rads)) * 1)

			next := Tuple{rads, gx, gy}
			iInput = append(iInput, next)
			count++

		}

	}

	log.Println(strconv.Itoa(count))
}

//hash and load Santa's path into sets
func firstIntersection() {
	intersectionSet := set.NewSet()
	for j := 0; j < len(iInput); j++ {
		Coords := "(" + strconv.Itoa(iInput[j].gx) + ", " + strconv.Itoa(iInput[j].gy) + ")"
		hashedCoords := GetMD5Hash(Coords)
		flag := intersectionSet.Contains(hashedCoords)

		if flag == true {
			gx := strconv.Itoa(iInput[j].gx)
			gy := strconv.Itoa(iInput[j].gy)
			strFlag := strconv.FormatBool(flag)
			fmt.Printf("(%s, %s) - %s -(%v) \n", gx, gy, hashedCoords, strFlag)
		} else {
			gx := strconv.Itoa(iInput[j].gx)
			gy := strconv.Itoa(iInput[j].gy)
			strFlag := strconv.FormatBool(flag)
			fmt.Printf("(%s, %s) - %s -(%v) \n", gx, gy, hashedCoords, strFlag)

		}
		intersectionSet.Add(hashedCoords)
	}
	log.Println(strconv.Itoa(len(iInput)))
}

func main() {
	readFile()
	parseHeadings()
	drawPath()
	iIntersection()
	firstIntersection()
}

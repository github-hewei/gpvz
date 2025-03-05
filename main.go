package main

import (
	"gpvz/gpvz"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//dir, _ := os.Getwd()
	//gameImgPath := path.Join(dir, "pvz_tools", "res", "game10.png")
	//gameImg := gocv.IMRead(gameImgPath, gocv.IMReadColor)
	//defer gameImg.Close()
	//
	//fmt.Println("rows", gameImg.Rows())
	//fmt.Println("cols", gameImg.Cols())
	//
	//rect := image.Rect(600, 50, 800, 0)
	//topRight := gameImg.Region(rect)
	//gcv.IMWrite("debug_5.png", topRight)

	//sunHSVLower := gocv.NewScalar(27, 95, 220, 0)
	//sunHSVUpper := gocv.NewScalar(35, 148, 244, 0)
	//contour := pvz_tools.NewContour(sunHSVLower, sunHSVUpper, 3, 100, 1)
	//result := contour.FindRectangles(gameImg)
	//fmt.Println("result", result)

	//sunHSVLower := gocv.NewScalar(117, 66, 105, 0)
	//sunHSVUpper := gocv.NewScalar(120, 68, 130, 0)
	//contour := pvz_tools.NewContour(sunHSVLower, sunHSVUpper, 3, 100, 1)
	//result := contour.FindRectangles(gameImg)
	//fmt.Println("result", result)

	tools := gpvz.NewGPvz()
	go tools.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	tools.Quit()
}

package gpvz

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

// Contour 轮廓查找器
type Contour struct {
	Debug    bool
	HsvLower gocv.Scalar
	HsvUpper gocv.Scalar
	openSize int
	minArea  int
}

// NewContour 创建一个轮廓查找器
func NewContour(lower gocv.Scalar, upper gocv.Scalar, args ...int) *Contour {
	openSize := 3

	if len(args) > 0 {
		openSize = args[0]
	}

	minArea := 100
	if len(args) > 1 {
		minArea = args[1]
	}

	debug := false
	if len(args) > 2 {
		debug = true
	}

	return &Contour{
		Debug:    debug,
		HsvLower: lower,
		HsvUpper: upper,
		openSize: openSize,
		minArea:  minArea,
	}
}

// FindRectangles 从图像中查找矩形
func (c *Contour) FindRectangles(img gocv.Mat) []image.Rectangle {
	// 第一步，转换到HSV颜色空间（更易提取特定颜色）
	hsv := gocv.NewMat()
	defer hsv.Close()
	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)
	_ = c.Debug && gocv.IMWrite("debug_1.png", hsv)

	// 第二步，创建二值化掩膜找出黄色区域
	mask := gocv.NewMat()
	gocv.InRangeWithScalar(hsv, c.HsvLower, c.HsvUpper, &mask)
	_ = c.Debug && gocv.IMWrite("debug_2.png", mask)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(c.openSize, c.openSize))
	//gocv.Dilate(mask, &mask, kernel)

	// 第三步，开运算去除小噪点
	gocv.MorphologyEx(mask, &mask, gocv.MorphOpen, kernel)
	_ = c.Debug && gocv.IMWrite("debug_3.png", mask)

	// 第四步，膨胀操作提升轮廓的边界
	gocv.Dilate(mask, &mask, kernel)
	_ = c.Debug && gocv.IMWrite("debug_4.png", mask)

	// 第五步，查找轮廓
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	result := make([]image.Rectangle, 0)
	if contours.Size() == 0 {
		return result
	}

	for i := 0; i < contours.Size(); i++ {
		// 计算轮廓面积， 过滤掉过小的轮廓
		area := gocv.ContourArea(contours.At(i))
		if area < float64(c.minArea) {
			continue
		}

		// 获取轮廓的边界框
		rect := gocv.BoundingRect(contours.At(i))
		result = append(result, rect)
	}

	if c.Debug {
		// 在图像上绘制轮廓
		for _, rect := range result {
			gocv.Rectangle(&img, rect, color.RGBA{R: 255}, 1)
		}
		gocv.IMWrite("debug_5.png", img)
	}

	return result
}

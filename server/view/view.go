package view

import (
	"ccproj/server/db"
	"ccproj/server/types"
	"fmt"
	"io/ioutil"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


var app *tview.Application

// var agents map[string]nmsServer.AgentRegistration
var logView *tview.TextView
var mM *tview.List
var menuStack []*tview.List
var logManager = db.NewLogManager()
var agents map[string]types.Agent


// StartGUI inicia a interface gráfica no terminal
func StartGUI(agentMap map[string]types.Agent, lm *db.LogManager) {
	
	logManager = lm
	agents = agentMap

	app = tview.NewApplication()
	logView = tview.NewTextView().SetDynamicColors(true).SetScrollable(true)

	mainMenu := tview.NewList().
		AddItem("Clients", "View connected clients", 'c', showClientsMenu).
		AddItem("Metrics", "View all metrics", 'm', showAllMetrics).
		AddItem("Logs", "View server logs", 'l', showLogs).
		AddItem("Quit", "", 'q', func() {
			app.Stop()
		})

	mM = mainMenu

	mainMenu.SetBorder(true).SetTitle("Main Menu").SetTitleAlign(tview.AlignLeft)

	menuStack = []*tview.List{mainMenu}

	if err := app.SetRoot(mainMenu, true).Run(); err != nil {
		panic(err)
	}
}

// showClientsMenu exibe o menu de clientes conectados
func showClientsMenu() {
	clientsMenu := tview.NewList()

	for id, agent := range agents {
		agentID := id
		clientsMenu.AddItem(fmt.Sprintf("%s - %s", agent.AgentID, agent.AgentIP), "", 0, func() {
			showClientFiles(agentID)
		})
	}

	clientsMenu.AddItem("Back", "", 'b', func() {
		app.SetRoot(popMenu(), true)
	})

	clientsMenu.SetBorder(true).SetTitle("Clients").SetTitleAlign(tview.AlignLeft)
	pushMenu(clientsMenu)
	app.SetRoot(clientsMenu, true)
}

// showClientFiles exibe os arquivos .txt disponíveis para um cliente
func showClientFiles(agentID string) {
	clientFilesMenu := tview.NewList()
	clientDir := fmt.Sprintf("../client_metrics/%s", agentID)

	clientFilesMenu.AddItem("Client's Log", "", 'l', func() {
		showClientLogs(agentID)
	})

	files, err := ioutil.ReadDir(clientDir)
	if err != nil {
		clientFilesMenu.AddItem("Error reading client directory", "", 0, nil)
	} else {
		for _, file := range files {
			if !file.IsDir() && file.Mode().IsRegular() && file.Name() != "" {
				fileName := file.Name()
				clientFilesMenu.AddItem(fileName, "", 0, func() {
					showFileContent(clientDir, fileName)
				})
			}
		}
	}

	clientFilesMenu.AddItem("Back", "", 'b', func() {
		app.SetRoot(popMenu(), true)
	})

	clientFilesMenu.SetBorder(true).SetTitle(fmt.Sprintf("Files for %s", agentID)).SetTitleAlign(tview.AlignLeft)
	pushMenu(clientFilesMenu)
	app.SetRoot(clientFilesMenu, true)
}

// showClientLogs exibe os logs do cliente
func showClientLogs(clientID string) {
	logs, err := logManager.GetAllLogs(clientID)
	if err != nil {
		logs = []string{fmt.Sprintf("Error reading logs: %v", err)}
	}

	logContent := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	for _, log := range logs {
		fmt.Fprintln(logContent, log)
	}

	logContent.SetBorder(true).SetTitle(fmt.Sprintf("Logs for %s", clientID)).SetTitleAlign(tview.AlignLeft)
	logContent.SetDoneFunc(func(key tcell.Key) {
		app.SetRoot(popMenu(), true)
	})

	pushMenu(logContent)
	app.SetRoot(logContent, true)
}

func showAllMetrics() {
	metricsMenu := tview.NewList()
	clientMetricsDir := "../client_metrics"

	folders, err := ioutil.ReadDir(clientMetricsDir)
	if err != nil {
		metricsMenu.AddItem("Error reading client_metrics directory", "", 0, nil)
	} else {
		for _, folder := range folders {
			if folder.IsDir() {
				folderName := folder.Name()
				metricsMenu.AddItem(folderName, "", 0, func() {
					showClientFiles(folderName)
				})
			}
		}
	}

	metricsMenu.AddItem("Back", "", 'b', func() {
		app.SetRoot(popMenu(), true)
	})

	metricsMenu.SetBorder(true).SetTitle("All Metrics").SetTitleAlign(tview.AlignLeft)
	pushMenu(metricsMenu)
	app.SetRoot(metricsMenu, true)
}

// showFileContent exibe o conteúdo de um arquivo .txt
func showFileContent(clientDir, fileName string) {
	filePath := fmt.Sprintf("%s/%s", clientDir, fileName)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		content = []byte(fmt.Sprintf("Error reading file: %v", err))
	}

	fileContentView := tview.NewTextView().
		SetText(string(content)).
		SetDynamicColors(true).
		SetScrollable(true)

	fileContentView.SetBorder(true).SetTitle(fileName).SetTitleAlign(tview.AlignLeft)
	fileContentView.SetDoneFunc(func(key tcell.Key) {
		app.SetRoot(popMenu(), true)
	})

	pushMenu(fileContentView)
	app.SetRoot(fileContentView, true)
}

// showLogs exibe os logs do servidor em tempo real
func showLogs() {
	logs, err := logManager.GetAllGeneralLogs()
    if err != nil {
        logs = []string{fmt.Sprintf("Error reading logs: %v", err)}
    }

    logContent := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
    for _, log := range logs {
        fmt.Fprintln(logContent, log)
    }

    logContent.SetBorder(true).SetTitle("Logs").SetTitleAlign(tview.AlignLeft)
    logContent.SetDoneFunc(func(key tcell.Key) {
        app.SetRoot(popMenu(), true)
    })

    pushMenu(logContent)
    app.SetRoot(logContent, true)
}


// AddLog adiciona uma entrada de log
func AddLog(log string) {
	fmt.Fprintln(logView, log)
	app.Draw()
}

// pushMenu adiciona um menu à stack
func pushMenu(menu tview.Primitive) {
	if list, ok := menu.(*tview.List); ok {
		menuStack = append(menuStack, list)
	}
}

// popMenu remove e retorna o menu do topo da stack
func popMenu() tview.Primitive {
	if len(menuStack) == 0 {
		return mM
	}
	menu := menuStack[len(menuStack)-1]
	menuStack = menuStack[:len(menuStack)-1]
	return menu
}

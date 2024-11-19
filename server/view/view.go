package view

import (
    "fmt"
    "io/ioutil"
    "github.com/gdamore/tcell/v2"
    //"ccproj/server/nmsServer"
    "github.com/rivo/tview"
    "strings"
)

var agents = map[string]string{
	"Client1": "192.168.1.10",
	"Client2": "192.168.1.20",
	"Client3": "192.168.1.30",
}

var app *tview.Application
//var agents map[string]nmsServer.AgentRegistration
var logView *tview.TextView
var mM *tview.List
var menuStack []*tview.List

// StartGUI inicia a interface gráfica no terminal
func StartGUI() {
    app = tview.NewApplication()
    //agents = make(map[string]nmsServer.AgentRegistration)
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

    // for id, agent := range agents {
    //     agentID := id
    //     clientsMenu.AddItem(fmt.Sprintf("%s - %s", agent.AgentID, agent.IPv4.String()), "", 0, func() {
    //         showClientFiles(agentID)
    //     })
    // }

	for name, ip := range agents {
        clientName := name
        clientIP := ip
        clientsMenu.AddItem(fmt.Sprintf("%s - %s", clientName, clientIP), "", 0, func() {
            showClientFiles(clientName)
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
    clientDir := fmt.Sprintf("../ClientMetrics/%s", agentID)

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

func showAllMetrics() {
    metricsMenu := tview.NewList()
    clientMetricsDir := "../ClientMetrics"

    folders, err := ioutil.ReadDir(clientMetricsDir)
    if err != nil {
        metricsMenu.AddItem("Error reading ClientMetrics directory", "", 0, nil)
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
    logView.SetBorder(true).SetTitle("Logs").SetTitleAlign(tview.AlignLeft)
    logView.SetDoneFunc(func(key tcell.Key) {
        app.SetRoot(popMenu(), true)
    })
    pushMenu(logView)
    app.SetRoot(logView, true)
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
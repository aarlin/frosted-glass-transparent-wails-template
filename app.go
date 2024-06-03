package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
    "log"
    "sync"
    "syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
    "github.com/lxn/win"
)

// App struct
type App struct {
	ctx context.Context
  	entities map[string]*Entity
    logLines int
	mu       sync.Mutex
}

// Entity struct to hold entity data
type Entity struct {
	HP              float64
	DPS10s          float64
	DPS60s          float64
	DPSOnThatEnemy  float64
	Name            string
	LastCombatTime  time.Time
}

// NewApp creates a new App application struct
func NewApp() *App {
  return &App{
		entities: make(map[string]*Entity),
        logLines: 0,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
    hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr("wuthering-waves-dps-meter"))
    win.SetWindowLong(hwnd, win.GWL_EXSTYLE, win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)

    a.ctx = ctx
}

func (a *App) LogLine() {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.logLines++
}

func (a *App) GetLogLineCount() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    return a.logLines
}

func (a *App) StartDpsTracker() {
  go a.startLogReader()
}

func (a *App) startLogReader() {
	paths := []string{
		"C:\\Program Files\\Epic Games\\WutheringWavesj3oFh\\Wuthering Waves Game\\Client\\Saved\\Logs\\Client.log",
		"C:\\Wuthering Waves\\Wuthering Waves Game\\Client\\Saved\\Logs\\Client.log",
	}

	var file *os.File
	var err error

	for _, path := range paths {
		file, err = os.Open(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	currentTime := time.Now()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		a.parseLine(line, currentTime)
		time.Sleep(200 * time.Millisecond)
	}
}

func (a *App) parseLine(line string, currentTime time.Time) {
	// Regex patterns
    a.LogLine()

	reAliveValue := regexp.MustCompile(`LifeValue: (\d+)`)
	reDeadValue := regexp.MustCompile(`CharacterAiComponent\.SetEnable`)
	reCombatStartEntity := regexp.MustCompile(`\[CombatInfo\].*\[EntityId:(\d+):Monster:BP_([^_]*)`)
	reCombatEndEntity := regexp.MustCompile(`\[CombatInfo\].*\[EntityId:(\d+):Vision:BP_([^_]*)`)


	// Parse timestamp
	logTimestamp, err := parseTimestamp(line)
	if err != nil {
		return
	}

	// Extract and process LifeValue
    if reAliveValue.MatchString(line) || reDeadValue.MatchString(line) {
        log.Println(line)
        var entityID string
        var entityName string
        var pastHP float64
        if reAliveValue.MatchString(line) {
            log.Println("Start combat with monster")
            lifeValueStr := reAliveValue.FindStringSubmatch(line)[1]
            pastHP, _ = strconv.ParseFloat(lifeValueStr, 64)
            entityID := reCombatStartEntity.FindStringSubmatch(line)[1]
            log.Println(entityID)
            entityName := reCombatStartEntity.FindStringSubmatch(line)[2]
            log.Println(entityName)
        } else {
            pastHP = 0
            entityID := reCombatEndEntity.FindStringSubmatch(line)[1]
            log.Println(entityID)
            entityName := reCombatEndEntity.FindStringSubmatch(line)[2]
            log.Println(entityName)
        }


        a.modifyEntity(entityID, entityName, pastHP, currentTime, logTimestamp)
    }

   a.updateUI()
}

func (a *App) modifyEntity(entityID string, entityName string, pastHP float64, currentTime time.Time, logTimestamp time.Time) {
    if _, exists := a.entities[entityID]; !exists {
        a.entities[entityID] = &Entity{
            HP:             pastHP,
            Name:           entityName,
            LastCombatTime: logTimestamp,
        }
    }

    entity := a.entities[entityID]
    if (currentTime.Sub(logTimestamp).Seconds() < 10) && (currentTime.Sub(logTimestamp).Seconds() > 0) {
        entity.DPS10s = (pastHP - entity.HP) / currentTime.Sub(logTimestamp).Seconds()
    }
    if (currentTime.Sub(logTimestamp).Seconds() < 60) && (currentTime.Sub(logTimestamp).Seconds() > 0) {
        entity.DPS60s = (pastHP - entity.HP) / currentTime.Sub(logTimestamp).Seconds()
    }
    if (entity.LastCombatTime.Sub(logTimestamp).Seconds() != 0) {
        entity.DPSOnThatEnemy = (pastHP - entity.HP) / entity.LastCombatTime.Sub(logTimestamp).Seconds()
    }
    entity.HP = pastHP
    entity.LastCombatTime = logTimestamp
}

func parseTimestamp(line string) (time.Time, error) {
	re := regexp.MustCompile(`\[(\d{4}\.\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{2})`)
	match := re.FindStringSubmatch(line)
	if match != nil {
		return time.Parse("2006.01.02-15.04.05", match[1])
	}
	return time.Time{}, fmt.Errorf("timestamp not found")
}

// updateUI updates the DPS values in the UI
func (a *App) updateUI() {
	// This function should send the updated DPS values to the frontend for display
	// Use Wails binding to call a frontend method

	runtime.EventsEmit(a.ctx, "rcv:entities", a.entities)
	runtime.EventsEmit(a.ctx, "rcv:logLines", a.GetLogLineCount())
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

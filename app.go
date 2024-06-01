package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// App struct
type App struct {
	ctx context.Context
  entities map[string]*Entity
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
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
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
	reLifeValue := regexp.MustCompile(`LifeValue: (\d+)`)
	reEntityID := regexp.MustCompile(`\[EntityId:(\d+):Monster:BP_([^_]*)`)

	// Parse timestamp
	time, err := parseTimestamp(line)
	if err != nil {
		return
	}

	// Extract and process LifeValue
	if reLifeValue.MatchString(line) {
		lifeValueStr := reLifeValue.FindStringSubmatch(line)[1]
		pastHP, _ := strconv.ParseFloat(lifeValueStr, 64)

		entityID := reEntityID.FindStringSubmatch(line)[1]
		entityName := reEntityID.FindStringSubmatch(line)[2]

		if _, exists := a.entities[entityID]; !exists {
			a.entities[entityID] = &Entity{
				HP:             pastHP,
				Name:           entityName,
				LastCombatTime: time,
			}
		}

		entity := a.entities[entityID]
		if (currentTime.Sub(time).Seconds() < 10) && (currentTime.Sub(time).Seconds() > 0) {
			entity.DPS10s = (pastHP - entity.HP) / currentTime.Sub(time).Seconds()
		}
		if (currentTime.Sub(time).Seconds() < 60) && (currentTime.Sub(time).Seconds() > 0) {
			entity.DPS60s = (pastHP - entity.HP) / currentTime.Sub(time).Seconds()
		}
		if (entity.LastCombatTime.Sub(time).Seconds() != 0) {
			entity.DPSOnThatEnemy = (pastHP - entity.HP) / entity.LastCombatTime.Sub(time).Seconds()
		}
		entity.HP = pastHP
		entity.LastCombatTime = time
	}

	// Update UI
	a.updateUI()
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

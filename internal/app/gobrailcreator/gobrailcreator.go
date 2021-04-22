package gobrailcreator

// * gets configuration from config/reader (creation plan)
// * can add elements (by json or objects?)
// * can provide creation to run railroad

import (
	"fmt"
	"strings"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	"github.com/gen2thomas/gobrail/internal/boardsapi"
	"github.com/gen2thomas/gobrail/internal/raildevicesapi"
	"github.com/gen2thomas/gobrail/internal/railplan"
)

// RailRunner is an interface to poll the rail
type RailRunner interface {
	Run() (err error)
}

type i2cAdaptor interface {
	i2c.Connector
	gobot.Connection
}

// AdaptorType represents the supported adaptors
type AdaptorType uint8

const (
	digisparkType AdaptorType = iota
	raspiType
	tinkerboardType
	unknownType
)

// RecipeFiles contains additional files to add to menu card
type RecipeFiles struct {
	boards  []string
	devices []string
}

var adaptorTypeToStringMap = map[AdaptorType]string{digisparkType: "digispark", raspiType: "raspi", tinkerboardType: "tinkerboard", unknownType: "typUnknown"}
var adaptorStringToTypeMap = map[string]AdaptorType{"digispark": digisparkType, "raspi": raspiType, "tinkerboard": tinkerboardType, "unknown": unknownType}

var lastGobot *gobot.Robot

// Create will create a static device connection for run
// before creating, the old gobot robot will be stopped
// after creating the devices, a new gobot robot will be created and started
func Create(daemonMode bool, name string, adaptorType AdaptorType, planFile string, recipeFiles RecipeFiles) (runner RailRunner, err error) {
	if err = Stop(); err != nil {
		return
	}

	fmt.Printf("\n------ Init gobot adaptor (%s) ------\n", adaptorType)
	var adaptor i2cAdaptor
	if adaptor, err = createAdaptor(adaptorType); err != nil {
		return
	}
	fmt.Printf("\n------ Init APIs ------\n")
	boardsAPI := boardsapi.NewBoardsAPI(adaptor)
	deviceAPI := raildevicesapi.NewRailDevicesAPI(boardsAPI)
	fmt.Printf("\n------ Read Plan (%s) ------\n", planFile)
	var plan railplan.RailPlan
	if plan, err = railplan.ReadMenuCard(planFile); err != nil {
		return
	}
	if len(recipeFiles.boards) > 0 {
		fmt.Printf("\n------ Read and add some board recipes ------\n")
		for _, boardFile := range recipeFiles.boards {
			if err = plan.AddBoardRecipe(boardFile); err != nil {
				return
			}
		}
	}
	if len(recipeFiles.devices) > 0 {
		fmt.Printf("\n------ Read and add some device recipes ------\n")
		for _, deviceFile := range recipeFiles.devices {
			if err = plan.AddDeviceRecipe(deviceFile); err != nil {
				return
			}
		}
	}
	fmt.Printf("\n------ Init boards from recipe list ------\n")
	for _, boardRecipe := range plan.BoardRecipes {
		if err = boardsAPI.AddBoard(boardRecipe); err != nil {
			return
		}
	}
	fmt.Printf("\n------ Init devices from recipe list ------\n")
	for _, deviceRecipe := range plan.DeviceRecipes {
		if err = deviceAPI.AddDevice(deviceRecipe); err != nil {
			return
		}
	}
	fmt.Printf("\n------ Map inputs to outputs ------\n")
	if err = deviceAPI.ConnectNow(); err != nil {
		return
	}

	fmt.Printf("\n------ Debugging ------\n")
	boardsAPI.ShowAllConfigs()
	boardsAPI.ShowAllUsedInputs()

	if daemonMode {
		// cyclic call of "Run()" is done by daemon program
		lastGobot = gobot.NewRobot(name,
			[]gobot.Connection{adaptor},
			boardsAPI.GobotDevices(),
		)
		// very important for daemon mode
		lastGobot.AutoRun = false
	} else {
		work := func() {
			gobot.Every(10*time.Millisecond, func() {
				if err := deviceAPI.Run(); err != nil {
					fmt.Println(err)
				}
			})
		}

		lastGobot = gobot.NewRobot(name,
			[]gobot.Connection{adaptor},
			boardsAPI.GobotDevices(),
			work,
		)
	}

	if err = lastGobot.Start(); err != nil {
		return
	}

	return deviceAPI, nil
}

// Stop stops the gobot robot, when available
func Stop() (err error) {
	if lastGobot != nil {
		if lastGobot.Running() {
			fmt.Printf("\n------ Stop gobot (%s) ------\n", lastGobot.Name)
			err = lastGobot.Stop()
		}
		lastGobot = nil
	}
	return
}

// ParseAdaptorType try get adaptor type from string
func ParseAdaptorType(adaptorString string) (a AdaptorType, err error) {
	var ok bool
	if a, ok = adaptorStringToTypeMap[strings.ToLower(adaptorString)]; !ok {
		err = fmt.Errorf("Unknown adaptor %s", adaptorString)
	}
	return
}

func (a AdaptorType) String() string {
	return adaptorTypeToStringMap[a]
}

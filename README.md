## Hardware Setup
Please check: https://github.com/wchill/SwitchInputEmulator.git

## Software Setup

1. Launch go package: SwitchKeyboard (main.go)
   * This connects to UART (please modify the path/port), syncs the controller, and launches a REST server at address `:80`
   * The REST server forwards all incoming requests to the UART
2. Launch `switch.ahk`
    * This forwards all keyboards command to the go REST server
    * This forwards all raw mouse delta data to the go REST server
    * F9 to pause/unpause the script
    * If the script is running, it disables all key presses, hides the cursor and clips the cursor

Note that these 2 layers may add some delay (< 1 frame) but I wasn't able to find an easier way to achieve said AHK functionalities in Go.
package main

import (
	"flag"
	"fmt"
	"golang.org/x/term"
	"encoding/json"
	"io/ioutil"
	"log"
	"bufio"
	"os"
	"os/exec"
	"strings"
)

// Raw input keycodes
var up byte = 65
var down byte = 66
var escape byte = 27
var enter byte = 13
var keys = map[byte]bool {
	up: true,
	down: true,
}

type Config struct {
    Types   []string `json:"types"`
    Scopes  []string `json:"scopes"`
}

func loadConfig() (*Config, error) {
    data, err := ioutil.ReadFile("convcom.json")
    if err != nil {
        return nil, err
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

type Menu struct {
	Prompt  	string
	CursorPos 	int
	MenuItems 	[]*MenuItem
}

type MenuItem struct {
	Text     string
	ID       string
	SubMenu  *Menu
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		Prompt: prompt,
		MenuItems: make([]*MenuItem, 0),
	}
}

// AddItem will add a new menu option to the menu list
func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID: id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

// renderMenuItems prints the menu item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(m.MenuItems) -1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems) - 1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		if index == m.CursorPos {
			cursor = fmt.Sprintf("\033[1;33m> \033[0m") 
			menuItemText = fmt.Sprintf("\033[1;33m%s\033[0m", menuItemText) 
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

// Display will display the current menu options and awaits user selection
// It returns the users selected choice
func (m *Menu) Display() string {
	defer func() {
		// Show cursor again.
		fmt.Printf("\033[?25h")
	}()

	fmt.Printf("\033[1;35m%s\033[0m\n", m.Prompt )

	m.renderMenuItems(false)

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	for {
		keyCode := getInput()
		if keyCode == escape {
			return ""
		} else if keyCode == 24 { // Ctrl + X
            fmt.Println("\nProcess exited.")
            os.Exit(0) // Exit the program
        } else if keyCode == enter {
			menuItem := m.MenuItems[m.CursorPos]
			fmt.Println("\r")
			return menuItem.ID
		} else if keyCode == up {
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if keyCode == down {
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		}
	}
}

// getInput will read raw input from the terminal
func getInput() byte {
	// Open the terminal
	fd := int(os.Stdin.Fd())

	// Set terminal to raw mode
	oldState, err := term.GetState(fd)
	if err != nil {
		log.Fatal(err)
	}

	oldState, err = term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}

	// Read input
	var readBytes [3]byte
	n, err := os.Stdin.Read(readBytes[:])
	if err != nil {
		log.Fatal(err)
	}

	// Restore the terminal state
	if err := term.Restore(fd, oldState); err != nil {
		log.Fatal(err)
	}

	if n == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		if readBytes[0] == 24 {
            return 24 // Ctrl + X
        }
		return readBytes[0]
	}

	return 0
}

// commitAndPush creates a Git commit with the provided message and pushes it to the remote repository.
func (c Choices) commitAndPush(dryrun bool) error {
	// Execute git commit
	fmt.Printf("\033[2J\033[H")
	fmt.Printf("\033[4;37m\033[1;37mConventional Commit\033[0m\033[0m\n\n")
	fmt.Printf("* \033[46m%s\033[0m\033[42m%s\033[0m\033[41m%s\033[0m: %s\n\n", c.TypeChoice, c.ScopeChoice, c.BreakChoice, c.CommitMessage)
	commit := fmt.Sprintf("Commit ... %s%s%s: %s\n", c.TypeChoice, c.ScopeChoice, c.BreakChoice, c.CommitMessage)
	if dryrun {
		fmt.Printf(commit)
		return nil
	} 
	
	commitCmd := exec.Command("git", "commit", "-am", commit)
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Execute git push
	//pushCmd := exec.Command("git", "push")
	//if err := pushCmd.Run(); err != nil {
	//	return fmt.Errorf("failed to push changes: %w", err)
	//}

	return nil
}

// Struct to hold user choices
type Choices struct {
	TypeChoice    string
	ScopeChoice   string
	BreakChoice   string
	CommitMessage string
}

func (c Choices) renderCommit() {
	fmt.Printf("\033[2J\033[H")
	fmt.Printf("\033[4;37m\033[1;37mConventional Commit\033[0m\033[0m\n\n")
	if c.TypeChoice == "" {
		fmt.Printf("* \n\n")
	} else {
	fmt.Printf("* \033[46m%s\033[0m\033[42m%s\033[0m\033[41m%s\033[0m: %s\n\n", c.TypeChoice, c.ScopeChoice, c.BreakChoice, c.CommitMessage)
	}
}

// createConfigFile creates a config file with the specified name if it does not already exist.
func createConfigFile() error {
	fileName := "convcom.json"
	// Check if the file already exists
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return fmt.Errorf("config file %s already exists", fileName)
	}

	// Define the configuration data
	config := Config{
		Types:  []string{"build", "ci", "chore", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"},
		Scopes: []string{},
	}

	// Open the file for writing
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// Encode the config to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print with indent
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	fmt.Printf("Config file %s created successfully.\n", fileName)
	return nil
}

func main() {

	initFlag := flag.Bool("init", false, "Creates a standard configuration file")
	commitFlag := flag.Bool("commit", false, "Creates a git commit")
	dryrunFlag := flag.Bool("dryrun", false, "Simulates git commit")
	flag.Parse()

	switch {
		case *initFlag:
			if err := createConfigFile(); err != nil {
				fmt.Println("Error:", err)
			}
		case *commitFlag:
			config, err := loadConfig()
			if err != nil {
				log.Fatalf("Error loading config: %v", err)
			}

			// Initialize Choices struct
			choices := Choices{}

			choices.renderCommit()

			menuType := NewMenu("Choose a type")
			for _, ctype := range config.Types {
				menuType.AddItem(ctype, ctype)
			}

			choices.TypeChoice = menuType.Display()
			choices.renderCommit()

			menuScope := NewMenu("Choose a scope")
			menuScope.AddItem("none", "")
			for _, cscope := range config.Scopes {
				menuScope.AddItem(cscope, "("+cscope+")")
			}

			choices.ScopeChoice = menuScope.Display()
			choices.renderCommit()

			menuBreak := NewMenu("Breaking change?")

			menuBreak.AddItem("no", "")
			menuBreak.AddItem("yes", "!")

			choices.BreakChoice = menuBreak.Display()
			choices.renderCommit()

			// Prompt for commit message
			fmt.Print("Enter commit message: ")
			reader := bufio.NewReader(os.Stdin)
			choices.CommitMessage, _ = reader.ReadString('\n')
			choices.CommitMessage = strings.TrimSpace(choices.CommitMessage)

			choices.renderCommit()

			menuCommit := NewMenu("Push?")

			menuCommit.AddItem("no", "no")
			menuCommit.AddItem("yes", "yes")

			if menuCommit.Display() == "yes" {
				choices.commitAndPush(*dryrunFlag)
			}
		default:
			fmt.Println("No valid flag provided. Use -init to create a configuration file or -commit (-dryrun) create a conventional commit.")
	}
}
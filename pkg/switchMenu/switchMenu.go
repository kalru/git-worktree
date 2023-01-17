package switchMenu

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc, path string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) Path() string        { return i.path }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func get_worktrees() []string {
	cmd := exec.Command("sh", "-c", "git worktree list")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %s", err)
	}
	worktrees := strings.Split(string(out[:]), "\n")
	return worktrees[:len(worktrees)-1]
}

func get_master_branch() string {
	// git branch -l master main should return one of the 2 branches.
	// This assumes one of these is the primary branch
	cmd := exec.Command("sh", "-c", "git branch -l master main")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %s", err)
	}
	lines := strings.Split(string(out[:]), "\n")
	return lines[0][2:]
}

func get_files_changed(branch string) string {
	// the last line of the command: git diff --stat
	cmd := exec.Command("sh", "-c", fmt.Sprintf("git diff %s --stat", branch))
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %s", err)
	}
	lines := strings.Split(string(out[:]), "\n")
	if len(lines) < 2 {
		return fmt.Sprintf("No changes relative to %s", get_master_branch())
	}
	return lines[len(lines)-2]
}

func open(path string) {
	editor := viper.GetString("editor")
	if editor != "code" && editor != "pycharm" {
		fmt.Println("\nPlease set editor to one of the following options: [code, pycharm, ]")
		os.Exit(1)
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", editor, path))
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %s", err)
	}
	log.Printf("%s", out)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			fs := m.list.FilterState()
			if fs == list.Unfiltered || fs == list.FilterApplied {
				// log.Printf("Index: %v, Selected item: %v", m.list.Index(), m.list.SelectedItem())
				// Not sure why it is so hard to get path here, but this is the only way I see now...
				open(fmt.Sprintf("%s", reflect.ValueOf(m.list.SelectedItem()).FieldByName("path")))
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Run() {
	if viper.GetBool("debug") {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	items := []list.Item{}

	for _, tree := range get_worktrees() {
		re_branch := regexp.MustCompile(`\[([^][]*)]`)
		branch := re_branch.FindString(tree)
		branch = branch[1 : len(branch)-1]
		re_path := regexp.MustCompile(`^([^\s]+)`)
		path := re_path.FindString(tree)
		items = append(items, item{title: branch, desc: get_files_changed(branch), path: path})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Select Worktree"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

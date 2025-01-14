package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/gen2brain/beeep"
	"github.com/icza/gox/timex"
	"github.com/mattn/go-mastodon"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

type URL struct {
	Text    string
	URL     string
	Classes []string
}

//Runs commands prefixed !CMD!
func CmdToString(cmd string) (string, error) {
	cmd = strings.TrimPrefix(cmd, "!CMD!")
	parts := strings.Split(cmd, " ")
	s, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	return strings.TrimSpace(string(s)), err
}

func getURLs(text string) []URL {
	doc := html.NewTokenizer(strings.NewReader(text))
	var urls []URL

	for {
		n := doc.Next()
		switch n {
		case html.ErrorToken:
			return urls

		case html.StartTagToken:
			token := doc.Token()
			if token.Data == "a" {
				url := URL{}
				var appendUrl = true
				for _, a := range token.Attr {
					switch a.Key {
					case "href":
						url.URL = a.Val
						url.Text = a.Val
					case "class":
						url.Classes = strings.Split(a.Val, " ")

						if strings.Contains(a.Val, "hashtag") {
							appendUrl = false
						}
					}
				}
				if appendUrl {
					urls = append(urls, url)
				}
			}
		}
	}
}

func cleanTootHTML(content string) (string, []URL) {
	stripped := bluemonday.NewPolicy().AllowElements("p", "br").AllowAttrs("href", "class").OnElements("a").Sanitize(content)
	urls := getURLs(stripped)
	stripped = bluemonday.NewPolicy().AllowElements("p", "br").Sanitize(content)
	stripped = strings.ReplaceAll(stripped, "<br>", "\n")
	stripped = strings.ReplaceAll(stripped, "<br/>", "\n")
	stripped = strings.ReplaceAll(stripped, "<p>", "")
	stripped = strings.ReplaceAll(stripped, "</p>", "\n\n")
	stripped = strings.TrimSpace(stripped)
	stripped = html.UnescapeString(stripped)
	return stripped, urls
}

func openEditor(app *tview.Application, content string) (string, error) {
	editor, exists := os.LookupEnv("EDITOR")
	if !exists || editor == "" {
		editor = "vi"
	}
	args := []string{}
	parts := strings.Split(editor, " ")
	if len(parts) > 1 {
		args = append(args, parts[1:]...)
		editor = parts[0]
	}
	f, err := ioutil.TempFile("", "tut")
	if err != nil {
		return "", err
	}
	if content != "" {
		_, err = f.WriteString(content)
		if err != nil {
			return "", err
		}
	}
	args = append(args, f.Name())
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var text []byte
	app.Suspend(func() {
		err = cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
		f.Seek(0, 0)
		text, err = ioutil.ReadAll(f)
	})
	f.Close()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(text)), nil
}

func copyToClipboard(text string) bool {
	if clipboard.Unsupported {
		return false
	}
	clipboard.WriteAll(text)
	return true
}

func openCustom(app *tview.Application, program string, args []string, terminal bool, url string) {
	args = append(args, url)
	if terminal {
		openInTerminal(app, program, args...)
	} else {
		exec.Command(program, args...).Start()
	}
}

func openURL(app *tview.Application, conf MediaConfig, pc OpenPatternConfig, url string) {
	for _, m := range pc.Patterns {
		if m.Compiled.Match(url) {
			args := append(m.Args, url)
			if m.Terminal {
				openInTerminal(app, m.Program, args...)
			} else {
				exec.Command(m.Program, args...).Start()
			}
			return
		}
	}
	args := append(conf.LinkArgs, url)
	if conf.LinkTerminal {
		openInTerminal(app, conf.LinkViewer, args...)
	} else {
		exec.Command(conf.LinkViewer, args...).Start()
	}
}

func reverseFiles(filenames []string) []string {
	if len(filenames) == 0 {
		return filenames
	}
	var f []string
	for i := len(filenames) - 1; i >= 0; i-- {
		f = append(f, filenames[i])
	}
	return f
}

type runProgram struct {
	Name     string
	Args     []string
	Terminal bool
}

func newRunProgram(name string, args ...string) runProgram {
	return runProgram{
		Name: name,
		Args: args,
	}
}

func openMediaType(app *tview.Application, conf MediaConfig, filenames []string, mediaType string) {
	terminal := []runProgram{}
	external := []runProgram{}

	switch mediaType {
	case "image":
		if conf.ImageReverse {
			filenames = reverseFiles(filenames)
		}
		if conf.ImageSingle {
			for _, f := range filenames {
				args := append(conf.ImageArgs, f)
				c := newRunProgram(conf.ImageViewer, args...)
				if conf.ImageTerminal {
					terminal = append(terminal, c)
				} else {
					external = append(external, c)
				}
			}
		} else {
			args := append(conf.ImageArgs, filenames...)
			c := newRunProgram(conf.ImageViewer, args...)
			if conf.ImageTerminal {
				terminal = append(terminal, c)
			} else {
				external = append(external, c)
			}
		}
	case "video", "gifv":
		if conf.VideoReverse {
			filenames = reverseFiles(filenames)
		}
		if conf.VideoSingle {
			for _, f := range filenames {
				args := append(conf.VideoArgs, f)
				c := newRunProgram(conf.VideoViewer, args...)
				if conf.VideoTerminal {
					terminal = append(terminal, c)
				} else {
					external = append(external, c)
				}
			}
		} else {
			args := append(conf.VideoArgs, filenames...)
			c := newRunProgram(conf.VideoViewer, args...)
			if conf.VideoTerminal {
				terminal = append(terminal, c)
			} else {
				external = append(external, c)
			}
		}
	case "audio":
		if conf.AudioReverse {
			filenames = reverseFiles(filenames)
		}
		if conf.AudioSingle {
			for _, f := range filenames {
				args := append(conf.AudioArgs, f)
				c := newRunProgram(conf.AudioViewer, args...)
				if conf.AudioTerminal {
					terminal = append(terminal, c)
				} else {
					external = append(external, c)
				}
			}
		} else {
			args := append(conf.AudioArgs, filenames...)
			c := newRunProgram(conf.AudioViewer, args...)
			if conf.AudioTerminal {
				terminal = append(terminal, c)
			} else {
				external = append(external, c)
			}
		}
	}
	go func() {
		for _, ext := range external {
			exec.Command(ext.Name, ext.Args...).Run()
		}
	}()
	for _, term := range terminal {
		openInTerminal(app, term.Name, term.Args...)
	}
}

func openInTerminal(app *tview.Application, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var err error
	app.Suspend(func() {
		err = cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
	})
	return err
}

func downloadFile(url string) (string, error) {
	f, err := ioutil.TempFile("", "tutfile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", nil
	}

	return f.Name(), nil
}

func getConfigDir() string {
	home, _ := os.LookupEnv("HOME")
	xdgConfig, exists := os.LookupEnv("XDG_CONFIG_HOME")
	if !exists {
		xdgConfig = home + "/.config"
	}
	xdgConfig += "/tut"
	return xdgConfig
}

func testConfigPath(name string) (string, error) {
	xdgConfig := getConfigDir()
	path := xdgConfig + "/" + name
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return path, nil
}

func GetAccountsPath() (string, error) {
	return testConfigPath("accounts.yaml")
}

func GetConfigPath() (string, error) {
	return testConfigPath("config.ini")
}

func CheckPath(input string, inclHidden bool) (string, bool) {
	info, err := os.Stat(input)
	if err != nil {
		return "", false
	}
	if !inclHidden && strings.HasPrefix(info.Name(), ".") {
		return "", false
	}

	if info.IsDir() {
		if input == "/" {
			return input, true
		}
		return input + "/", true
	}
	return input, true
}

func IsDir(input string) bool {
	info, err := os.Stat(input)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FindFiles(s string) []string {
	input := filepath.Clean(s)
	if len(s) > 2 && s[len(s)-2:] == "/." {
		input += "/."
	}
	var files []string
	path, exists := CheckPath(input, true)
	if exists {
		files = append(files, path)
	}

	base := filepath.Base(input)
	inclHidden := strings.HasPrefix(base, ".") || (len(input) > 1 && input[len(input)-2:] == "/.")
	matches, _ := filepath.Glob(input + "*")
	if strings.HasSuffix(path, "/") {
		matchesDir, _ := filepath.Glob(path + "*")
		matches = append(matches, matchesDir...)
	}
	for _, f := range matches {
		p, exists := CheckPath(f, inclHidden)
		if exists && p != path {
			files = append(files, p)
		}
	}
	return files
}

func ColorKey(c *Config, pre, key, end string) string {
	color := ColorMark(c.Style.TextSpecial2)
	normal := ColorMark(c.Style.Text)
	key = TextFlags("b") + key + TextFlags("-")
	if c.General.ShortHints {
		pre = ""
		end = ""
	}
	text := fmt.Sprintf("%s%s%s%s%s%s", normal, pre, color, key, normal, end)
	return text
}

func TextFlags(s string) string {
	return fmt.Sprintf("[::%s]", s)
}

func ColorMark(color tcell.Color) string {
	return fmt.Sprintf("[#%06x]", color.Hex())
}

func FormatUsername(a mastodon.Account) string {
	if a.DisplayName != "" {
		return fmt.Sprintf("%s (%s)", a.DisplayName, a.Acct)
	}
	return a.Acct
}

func SublteText(style StyleConfig, text string) string {
	subtle := ColorMark(style.Subtle)
	return fmt.Sprintf("%s%s", subtle, text)
}

func FloorDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func OutputDate(status time.Time, today time.Time, long, short string, relativeDate int) string {
	ty, tm, td := today.Date()
	sy, sm, sd := status.Date()

	format := long
	sameDay := false
	displayRelative := false

	if ty == sy && tm == sm && td == sd {
		format = short
		sameDay = true
	}

	todayFloor := FloorDate(today)
	statusFloor := FloorDate(status)

	if relativeDate > -1 && !sameDay {
		days := int(todayFloor.Sub(statusFloor).Hours() / 24)
		if relativeDate == 0 || days <= relativeDate {
			displayRelative = true
		}
	}
	var dateOutput string
	if displayRelative {
		y, m, d, _, _, _ := timex.Diff(statusFloor, todayFloor)
		if y > 0 {
			dateOutput = fmt.Sprintf("%s%dy", dateOutput, y)
		}
		if dateOutput != "" || m > 0 {
			dateOutput = fmt.Sprintf("%s%dm", dateOutput, m)
		}
		if dateOutput != "" || d > 0 {
			dateOutput = fmt.Sprintf("%s%dd", dateOutput, d)
		}
	} else {
		dateOutput = status.Format(format)
	}
	return dateOutput
}

func Notify(nc NotificationConfig, t NotificationType, title string, body string) {
	switch t {
	case NotificationFollower:
		if !nc.NotificationFollower {
			return
		}
	case NotificationFavorite:
		if !nc.NotificationFavorite {
			return
		}
	case NotificationMention:
		if !nc.NotificationMention {
			return
		}
	case NotificationBoost:
		if !nc.NotificationBoost {
			return
		}
	case NotificationPoll:
		if !nc.NotificationPoll {
			return
		}
	case NotificationPost:
		if !nc.NotificationPost {
			return
		}
	default:
		return
	}

	beeep.Notify(title, body, "")
}

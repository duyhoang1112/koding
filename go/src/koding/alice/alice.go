package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bobappleyard/readline"
	"io"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

type Entry struct {
	Name  string `json:"name"`
	VHost string `json:"vhost"`
}

var commands = map[string]func(args []string){
	"overview":    func(args []string) { PrintElement("overview", "overview") },
	"vhost":       func(args []string) { ChangeVHost(args) },
	"nodes":       func(args []string) { ReadList("node", false) },
	"connections": func(args []string) { ReadList("connection", false) },
	"channels":    func(args []string) { ReadList("channel", false) },
	"exchanges":   func(args []string) { ReadList("exchange", true) },
	"queues":      func(args []string) { ReadList("queue", true) },
	"vhosts":      func(args []string) { ReadList("vhost", false) },
	"users":       func(args []string) { ReadList("user", false) },
	"trace":       func(args []string) { Trace() },
	"graph":       func(args []string) { CreateGraph() },
	"exit":        func(args []string) { exit = true },
	"quit":        func(args []string) { exit = true },
}

var exit = false
var abort = false
var vhost = "/"
var listEntries []Entry = nil
var listKind string
var entryPathPrefix string

func main() {
	readline.CatchSigint = false
	readline.Completer = func(query, ctx string) []string {
		completions := make([]string, 0)
		for command := range commands {
			if strings.HasPrefix(command, query) {
				completions = append(completions, command)
			}
		}
		return completions
	}

	go func() {
		signals := make(chan os.Signal, 2)
		signal.Notify(signals, syscall.SIGINT)
		for _ = range signals {
			abort = true
		}
	}()

	for !exit {
		input, err := readline.String(readline.EscapePrompt("\x1b[1malice " + vhost + " > \x1b[0m"))
		if err == io.EOF {
			fmt.Println()
			break
		}
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		readline.AddHistory(input)

		parts := strings.Split(input, " ")
		command := commands[parts[0]]
		index, err := strconv.Atoi(input)
		if command != nil {
			command(parts[1:])
		} else if err == nil {
			if listEntries == nil {
				fmt.Println("No list.")
			} else if index < 0 || index >= len(listEntries) {
				fmt.Println("Index out of bounds.")
			} else {
				entry := listEntries[index]

				path := listKind + "s"
				if entry.VHost != "" {
					path += "/" + QueryEscape(entry.VHost)
				}
				path += "/" + QueryEscape(entry.Name)

				PrintElement(path, listKind+" "+entry.Name)
			}
		} else {
			fmt.Println("Sorry, unknown command.")
		}
	}
}

func ChangeVHost(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: vhost <name>")
		return
	}
	if Get("vhosts/"+QueryEscape(args[0]), "", nil) {
		vhost = args[0]
	} else {
		fmt.Println("No such vhost.")
	}
}

func ReadList(kind string, withVHost bool) {
	path := kind + "s"
	query := "columns=name"
	if withVHost {
		path += "/" + vhost
		query += ",vhost"
	}
	if !Get(path, query, &listEntries) {
		return
	}
	format := fmt.Sprintf("%%%dd: %%s\n", int(math.Ceil(math.Log10(float64(len(listEntries))))))
	for i, entry := range listEntries {
		fmt.Printf(format, i, entry.Name)
	}
	listKind = kind
}

func Trace() {
	fmt.Println("Under construction.")
	return

	var data map[string]interface{}
	Get("vhosts/"+QueryEscape(vhost), "", &data)
	if data["tracing"] != true {
		fmt.Println("Please enable firehose tracing for this vhost first.")
		return
	}
}

func CreateGraph() {
	fmt.Println("Under construction.")
	return

	file, err := os.Create("graph.dot")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write([]byte("digraph G {"))

	var exchanges []Entry = nil
	Get("exchanges/"+vhost, "columns=name", &exchanges)
	for _, exchange := range exchanges {
		fmt.Fprintf(file, "%s;", exchange.Name)
	}

	file.Write([]byte("}"))
}

func PrintElement(path, name string) {
	var data map[string]interface{}
	if !Get(path, "", &data) {
		return
	}
	fmt.Print(name + ":")
	PrettyPrint(data, 2, GetMaxKeyLength(data))
}

func Get(path, query string, v interface{}) bool {
	return DoRequest("GET", path, query, nil, v)
}

func Put(path string, data map[string]interface{}) bool {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return DoRequest("PUT", path, "", bytes.NewReader(body), nil)
}

func DoRequest(method, path, query string, body io.Reader, v interface{}) bool {
	client := &http.Client{}
	req, err := http.NewRequest(method, "http://web0.beta.system.aws.koding.com:55672", body)
	if err != nil {
		panic(err)
	}
	req.URL.Opaque = "/api/" + path
	req.URL.RawQuery = query
	req.SetBasicAuth("guest", "x1srTA7!%Vb}$n|S")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return false
	}
	if v != nil {
		body := make([]byte, resp.ContentLength)
		i := int64(0)
		abort = false
		for {
			if abort {
				fmt.Print(" Aborted.\n")
				return false
			}
			n, err := resp.Body.Read(body[i:])
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			i += int64(n)
			fmt.Printf("\x1b[GReceiving... (%d%%)", i*100/resp.ContentLength)
		}
		fmt.Print("\x1b[G\x1b[K")
		err = json.Unmarshal(body, v)
		if err != nil {
			panic(err)
		}
	}
	return true
}

func PrettyPrint(value interface{}, indentation, keyLength int) {
	format := fmt.Sprintf(strings.Repeat(" ", indentation)+"%%-%dv   ", keyLength+1)
	switch casted := value.(type) {
	case map[string]interface{}:
		fmt.Print("\n")
		keys := make([]string, len(casted))
		i := 0
		for key := range casted {
			keys[i] = key
			i += 1
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf(format, strings.Replace(key, "_", " ", -1)+":")
			PrettyPrint(casted[key], indentation+2, keyLength-2)
		}
	case []interface{}:
		fmt.Print("\n")
		for i, entry := range casted {
			fmt.Printf(format, strconv.Itoa(i)+":")
			PrettyPrint(entry, indentation+2, keyLength-2)
		}
	default:
		fmt.Printf("%v\n", casted)
	}
}

func GetMaxKeyLength(value interface{}) int {
	maxLength := 0
	switch casted := value.(type) {
	case map[string]interface{}:
		for key, entry := range casted {
			if len(key) > maxLength {
				maxLength = len(key)
			}
			subLength := GetMaxKeyLength(entry) + 2
			if subLength > maxLength {
				maxLength = subLength
			}
		}
	case []interface{}:
		maxLength := int(math.Ceil(math.Log10(float64(len(casted)))))
		for _, entry := range casted {
			subLength := GetMaxKeyLength(entry) + 2
			if subLength > maxLength {
				maxLength = subLength
			}
		}
	}
	return maxLength
}

func QueryEscape(s string) string {
	s = strings.Replace(s, " ", "%20", -1)
	s = strings.Replace(s, "/", "%2f", -1)
	s = strings.Replace(s, "@", "%40", -1)
	return s
}

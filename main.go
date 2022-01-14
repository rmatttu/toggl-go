package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"toggl-go/api"
	"toggl-go/year"
)

// Args is commandline args
type Args struct {
	Wait    *int
	Timeout *int
	Token   *string
	Email   *string
	Since   *string
}

func (u *Args) getJoinedArgs() string {
	text := []string{
		"Wait: " + strconv.Itoa(*u.Wait),
		"Timeout: " + strconv.Itoa(*u.Timeout),
		"Token: " + *u.Token,
		"Email: " + *u.Email,
		"Since: " + *u.Since,
	}
	return strings.Join(text, ", ")
}

// NullWriter is dummy target for loggger
// logger write to null
type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	// sw1 := flag.Bool("sw1", false, "オプションを付けるとsw1はtrue (default: false")
	var args Args
	args.Wait = flag.Int("wait", 3, "Sleep wait seconds")
	args.Timeout = flag.Int("timeout", 30, "Timeout seconds")
	args.Token = flag.String("token", "__token__", "API token")
	args.Email = flag.String("email", "your_mail_address@example.com", "Your email address")
	args.Since = flag.String("since", "1990", "Target since year")
	flag.Parse()

	log.SetOutput(new(NullWriter))

	log.Print("args: ", args.getJoinedArgs())
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Print("hello world") // 2021/02/01 10:13:42.326903 main.go:56: hello world

	conf := api.ClientConfigure{
		Timeout:  *args.Timeout,
		APIToken: *args.Token,
	}

	sinceYear, err := strconv.Atoi(*args.Since)
	if err != nil {
		panic(err)
	}
	nowYear := time.Now().Year()

	mainWorkspace, _ := api.FetchMainWorkspace(conf)
	for targetYear := sinceYear; targetYear <= nowYear; targetYear++ {
		target, _ := year.New(sinceYear)
		log.Print("target: " + target.String())
		for page := 1; true; page++ {
			req := api.DetailedReportRequest{
				UserAgent:   *args.Email,
				WorkspaceID: mainWorkspace.ID,
				Since:       target.Since,
				Until:       target.Until,
				Page:        page,
			}

			time.Sleep(time.Duration(*args.Wait) * time.Second)
			res, responseRaw, err := api.FetchDetailedReport(conf, req)
			if err != nil {
				panic(err)
			}

			fmt.Fprintf(os.Stdout, "%s\n", responseRaw)

			if res.TotalCount < res.PerPage*page {
				break
			}
		}
	}

}

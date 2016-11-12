package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/fatih/structs"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

const (
	name     = "f-prot"
	category = "av"
)

type pluginResults struct {
	ID   string      `json:"id" structs:"id,omitempty"`
	Data ResultsData `json:"f-prot" structs:"f-prot"`
}

// FPROT json object
type FPROT struct {
	Results ResultsData `json:"f-prot"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" structs:"infected"`
	Result   string `json:"result" structs:"result"`
	Engine   string `json:"engine" structs:"engine"`
	Updated  string `json:"updated" structs:"updated"`
}

// ParseFprotOutput convert fprot output into ResultsData struct
func ParseFprotOutput(fprotout string, err error) (ResultsData, error) {

	if err != nil {
		return ResultsData{}, err
	}

	fprot := ResultsData{Infected: false}
	colonSeparated := []string{}

	lines := strings.Split(fprotout, "\n")
	// Extract Virus string and extract colon separated lines into an slice
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, ":") {
				colonSeparated = append(colonSeparated, line)
			}
			if strings.Contains(line, "[Found virus]") {
				result := extractVirusName(line)
				if len(result) != 0 {
					fprot.Result = result
					fprot.Infected = true
				} else {
					fmt.Println("[ERROR] Virus name extracted was empty: ", result)
					os.Exit(2)
				}
			}
		}
	}
	// fmt.Println(lines)

	// Extract FPROT Details from scan output
	if len(colonSeparated) != 0 {
		for _, line := range colonSeparated {
			if len(line) != 0 {
				keyvalue := strings.Split(line, ":")
				if len(keyvalue) != 0 {
					switch {
					case strings.Contains(keyvalue[0], "Virus signatures"):
						fprot.Updated = parseUpdatedDate(strings.TrimSpace(keyvalue[1]))
					case strings.Contains(line, "Engine version"):
						fprot.Engine = strings.TrimSpace(keyvalue[1])
					}
				}
			}
		}
	} else {
		fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
		os.Exit(2)
	}

	// fprot.Updated = getUpdatedDate()

	return fprot, nil
}

// extractVirusName extracts Virus name from scan results string
func extractVirusName(line string) string {
	r := regexp.MustCompile(`<(.+)>`)
	res := r.FindStringSubmatch(line)
	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	utils.Assert(err)
	return string(updated)
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

func parseUpdatedDate(date string) string {
	layout := "200601021504"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func printMarkDownTable(fprot FPROT) {

	fmt.Println("#### F-PROT")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": fprot.Results.Infected,
		"Result":   fprot.Results.Result,
		"Engine":   fprot.Results.Engine,
		"Updated":  fprot.Results.Updated,
	})
	table.Markdown = true
	table.Print()
}

func updateAV(ctx context.Context) error {
	fmt.Println("Updating F-PROT...")
	fmt.Println(utils.RunCommand(ctx, "/opt/f-prot/fpupdate"))
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

func main() {

	var elastic string

	cli.AppHelpTemplate = utils.AppHelpTemplate
	app := cli.NewApp()
	app.Name = "fprot"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice F-PROT AntiVirus Plugin"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.BoolFlag{
			Name:   "post, p",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.StringFlag{
			Name:        "elasitcsearch",
			Value:       "",
			Usage:       "elasitcsearch address for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH",
			Destination: &elastic,
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  10,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update virus definitions",
			Action: func(c *cli.Context) error {
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Int("timeout"))*time.Second)
				defer cancel()
				return updateAV(ctx)
			},
		},
	}
	app.Action = func(c *cli.Context) error {

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Int("timeout"))*time.Second)
		defer cancel()

		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			utils.Assert(err)
		}
		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		results, err := ParseFprotOutput(utils.RunCommand(ctx, "/usr/local/bin/fpscan", "-r", path))
		utils.Assert(err)

		fprot := FPROT{
			Results: results,
		}

		// upsert into Database
		elasticsearch.InitElasticSearch(elastic)
		elasticsearch.WritePluginResultsToDatabase(elasticsearch.PluginResults{
			ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
			Name:     name,
			Category: category,
			Data:     structs.Map(fprot.Results),
		})

		if c.Bool("table") {
			printMarkDownTable(fprot)
		} else {
			fprotJSON, err := json.Marshal(fprot)
			utils.Assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("X-Malice-ID", utils.Getopt("MALICE_SCANID", utils.GetSHA256(path))).
					Send(string(fprotJSON)).
					End(printStatus)

				return nil
			}
			fmt.Println(string(fprotJSON))
		}
		return nil
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}

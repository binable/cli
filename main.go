package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const apiURL = "https://xn--mll-hoa.io/api/fetch"

var version = "dev"

type request struct {
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

type wasteType struct {
	Dates    []string `json:"dates"`
	Next     *string  `json:"next"`
	NextDays *string  `json:"nextDays"`
	Last     *string  `json:"last"`
}

type response struct {
	TypeMap map[string]string    `json:"typeMap"`
	Empty   bool                 `json:"empty"`
	Address struct {
		Street      string `json:"street"`
		HouseNumber string `json:"houseNumber"`
		Zip         string `json:"zip"`
		City        string `json:"city"`
		Country     string `json:"country"`
	} `json:"address"`
	Provider struct {
		Name string `json:"name"`
		URI  string `json:"uri"`
	} `json:"provider"`
	Paper                          wasteType `json:"paper"`
	Bio                            wasteType `json:"bio"`
	ResidualWaste                  wasteType `json:"residualWaste"`
	ResidualWasteLessFrequent      wasteType `json:"residualWasteLessFrequent"`
	ResidualWasteContainer         wasteType `json:"residualWasteContainer"`
	ResidualWasteContainerLessFreq wasteType `json:"residualWasteContainerLessFrequent"`
	ReusableMaterials              wasteType `json:"reusableMaterials"`
	ChristmasTree                  wasteType `json:"christmasTree"`
	Toxic                          wasteType `json:"toxic"`
	Diaper                         wasteType `json:"diaper"`
	HedgeTreeTrimming              wasteType `json:"hedgeTreeTrimming"`
	FleaMarket                     wasteType `json:"fleaMarket"`
}

type entry struct {
	key   string
	label string
	wt    wasteType
}

var colorReset = "\033[0m"
var colorBold = "\033[1m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorCyan = "\033[36m"
var colorGray = "\033[90m"

func init() {
	if os.Getenv("NO_COLOR") != "" || !isTerminal() {
		colorReset = ""
		colorBold = ""
		colorGreen = ""
		colorYellow = ""
		colorCyan = ""
		colorGray = ""
	}
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func fetch(req request) (*response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var result response
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid response: %w\n%s", err, string(data))
	}
	return &result, nil
}

func formatDate(s string) string {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return s
	}
	return t.Format("02.01.2006 (Mon)")
}

func daysColor(days string) string {
	switch days {
	case "0":
		return colorGreen + colorBold + "heute!" + colorReset
	case "1":
		return colorGreen + "morgen" + colorReset
	}
	return colorCyan + "in " + days + " Tagen" + colorReset
}

func printTable(entries []entry, typeMap map[string]string, jsonOut bool) {
	if jsonOut {
		out := map[string]interface{}{}
		for _, e := range entries {
			if e.wt.Next == nil {
				continue
			}
			label := typeMap[e.key]
			if label == "" {
				label = e.label
			}
			out[label] = map[string]interface{}{
				"next":     *e.wt.Next,
				"nextDays": e.wt.NextDays,
				"dates":    e.wt.Dates,
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		return
	}

	maxLabel := 0
	for _, e := range entries {
		if e.wt.Next == nil {
			continue
		}
		label := typeMap[e.key]
		if label == "" {
			label = e.label
		}
		if len(label) > maxLabel {
			maxLabel = len(label)
		}
	}

	fmt.Println()
	for _, e := range entries {
		if e.wt.Next == nil {
			continue
		}
		label := typeMap[e.key]
		if label == "" {
			label = e.label
		}
		padding := strings.Repeat(" ", maxLabel-len(label))
		days := ""
		if e.wt.NextDays != nil {
			days = "  " + daysColor(*e.wt.NextDays)
		}
		fmt.Printf("  %s%s%s%s  %s%s\n",
			colorBold, label, colorReset,
			padding,
			formatDate(*e.wt.Next),
			days,
		)
	}
	fmt.Println()
}

func run() int {
	fs := flag.NewFlagSet("binable", flag.ExitOnError)

	var street, house, zip, city, country string
	var jsonOut, all, ver bool

	fs.StringVar(&street, "street", "", "")
	fs.StringVar(&street, "s", "", "")
	fs.StringVar(&house, "house", "", "")
	fs.StringVar(&house, "n", "", "")
	fs.StringVar(&zip, "zip", "", "")
	fs.StringVar(&zip, "z", "", "")
	fs.StringVar(&city, "city", "", "")
	fs.StringVar(&city, "c", "", "")
	fs.StringVar(&country, "country", "DE", "")
	fs.StringVar(&country, "C", "DE", "")
	fs.BoolVar(&jsonOut, "json", false, "")
	fs.BoolVar(&jsonOut, "j", false, "")
	fs.BoolVar(&all, "all", false, "")
	fs.BoolVar(&all, "a", false, "")
	fs.BoolVar(&ver, "version", false, "")
	fs.BoolVar(&ver, "v", false, "")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Verwendung: binable --street <Straße> --house <Nr> --zip <PLZ> --city <Stadt>\n\nOptionen:\n")
		fmt.Fprintf(os.Stderr, "  -s, --street <Straße>    Straßenname\n")
		fmt.Fprintf(os.Stderr, "  -n, --house <Nr>         Hausnummer\n")
		fmt.Fprintf(os.Stderr, "  -z, --zip <PLZ>          Postleitzahl\n")
		fmt.Fprintf(os.Stderr, "  -c, --city <Stadt>       Stadt\n")
		fmt.Fprintf(os.Stderr, "  -C, --country <Land>     Länderkürzel (Standard: DE)\n")
		fmt.Fprintf(os.Stderr, "  -a, --all                Alle Termine des Jahres anzeigen\n")
		fmt.Fprintf(os.Stderr, "  -j, --json               Ausgabe als JSON\n")
		fmt.Fprintf(os.Stderr, "  -v, --version            Version anzeigen\n")
		fmt.Fprintf(os.Stderr, "\nBeispiel:\n  binable --street \"Schürhornweg\" --house 1 --zip 33649 --city Bielefeld\n")
	}

	fs.Parse(os.Args[1:])

	if ver {
		fmt.Println("binable", version)
		return 0
	}

	if street == "" || zip == "" || city == "" {
		fs.Usage()
		return 1
	}

	req := request{
		Street:      street,
		HouseNumber: house,
		Zip:         zip,
		City:        city,
		Country:     country,
	}

	fmt.Fprintf(os.Stderr, "%sLade Abfuhrdaten für %s %s, %s %s …%s\n",
		colorGray, street, house, zip, city, colorReset)

	result, err := fetch(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
		return 1
	}

	if result.Empty {
		fmt.Fprintln(os.Stderr, "Keine Daten gefunden für diese Adresse.")
		return 1
	}

	if !jsonOut {
		fmt.Printf("\n%sAnbieter:%s %s\n", colorBold, colorReset, result.Provider.Name)
		fmt.Printf("%sAdresse:%s  %s %s, %s %s, %s\n",
			colorBold, colorReset,
			result.Address.Street, result.Address.HouseNumber,
			result.Address.Zip, result.Address.City,
			result.Address.Country,
		)
	}

	entries := []entry{
		{"residualWaste", "Restmüll", result.ResidualWaste},
		{"bio", "Biomüll", result.Bio},
		{"paper", "Papiermüll", result.Paper},
		{"reusableMaterials", "Wertstoff", result.ReusableMaterials},
		{"residualWasteLessFrequent", "Restmüll (groß)", result.ResidualWasteLessFrequent},
		{"residualWasteContainer", "Restmüll Container", result.ResidualWasteContainer},
		{"residualWasteContainerLessFrequent", "Restmüll Container (groß)", result.ResidualWasteContainerLessFreq},
		{"christmasTree", "Weihnachtsbäume", result.ChristmasTree},
		{"toxic", "Schadstoffe", result.Toxic},
		{"diaper", "Windeln", result.Diaper},
		{"hedgeTreeTrimming", "Heckenschnitt", result.HedgeTreeTrimming},
		{"fleaMarket", "Trödelmarkt", result.FleaMarket},
	}

	// Sort by next date
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].wt.Next == nil {
			return false
		}
		if entries[j].wt.Next == nil {
			return true
		}
		return *entries[i].wt.Next < *entries[j].wt.Next
	})

	printTable(entries, result.TypeMap, jsonOut)

	if all && !jsonOut {
		fmt.Printf("%sAlle Termine:%s\n\n", colorBold, colorReset)
		for _, e := range entries {
			if len(e.wt.Dates) == 0 {
				continue
			}
			label := result.TypeMap[e.key]
			if label == "" {
				label = e.label
			}
			fmt.Printf("  %s%s%s\n", colorBold, label, colorReset)
			for _, d := range e.wt.Dates {
				fmt.Printf("    %s  %s\n", colorGray+"•"+colorReset, formatDate(d))
			}
			fmt.Println()
		}
	}

	return 0
}

func main() {
	os.Exit(run())
}

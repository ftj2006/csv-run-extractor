package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var fileCount int

type Config struct {
	RunLimits struct {
		ThrottleField    string  `yaml:"throttleField"`
		MinThrottleValue float64 `yaml:"minThrottleValue"`
		MinCount         int     `yaml:"minCount"`
	} `yaml:"runLimits"`
	Fields []string `yaml:"fields"`
}

func saveFile(filePath string, records [][]string) {
	fileCount++

	filePath = strings.Replace(filePath, ".csv", " - Run %d.csv", 1)
	csvFile, err := os.Create(fmt.Sprintf(filePath, fileCount))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	err = csvWriter.WriteAll(records)

	fmt.Printf("Saved Run %d to file %s\n", fileCount, csvFile.Name())
	if err != nil {
		log.Fatal(err)
	}
}

func buildRecord(record []string, fields map[string]int, config Config) []string {
	newRecord := []string{}

	for _, field := range config.Fields {
		fieldIndex, ok := fields[field]

		var value string
		if ok {
			value = record[fieldIndex]
		} else {
			value = "NOT FOUND"
		}

		newRecord = append(newRecord, value)
	}
	return newRecord
}
func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func processCsv(config Config, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	header, _ := csvReader.Read()

	fields := make(map[string]int)

	for i, name := range header {
		fields[name] = i
	}

	tPosCount := 0
	prevRecord := []string{}
	buffer := [][]string{}
	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		tPosStr := record[fields[config.RunLimits.ThrottleField]]

		tPos, err := strconv.ParseFloat(tPosStr, 64)
		if err != nil {
			log.Fatal("Unable to parse throttle position: ", err)
		}

		if tPos > config.RunLimits.MinThrottleValue {
			if tPosCount == 0 {
				buffer = append(buffer, config.Fields)
				buffer = append(buffer, buildRecord(prevRecord, fields, config))
			}
			tPosCount++

			buffer = append(buffer, buildRecord(record, fields, config))
		} else { //This is where the run ends
			buffer = append(buffer, buildRecord(record, fields, config))

			if tPosCount > config.RunLimits.MinCount {
				saveFile(filePath, buffer)
			}

			//Reset buffer and position count
			buffer = [][]string{}
			tPosCount = 0
		}
		prevRecord = record
	}
}

func loadConfig() (Config, error) {
	var config Config

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	filePath := filepath.Join(exPath, "config.yaml")

	fmt.Println(filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	file := os.Args[1]
	fileCount = 0

	config, err := loadConfig()
	if err != nil {
		log.Fatal("Unable to load config", err)
	}

	processCsv(config, file)
}

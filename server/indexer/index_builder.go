package indexer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type EnumIndexBuilder struct {
	Property string
	Config   IndexConfig
	ValueIds map[string][]Ids
}

type Ids struct {
	DataRowId int
}

func (enumBuilder *EnumIndexBuilder) AddIndexValue(dataRowId int, value string) {
	enumBuilder.ValueIds[value] = append(enumBuilder.ValueIds[value], Ids{dataRowId})
}

func (enumBuilder *EnumIndexBuilder) WriteIndex(fileId int, outDir string) (*EnumIndex, error) {
	values := make(map[string]*EnumValue)
	count := 0
	var err error
	for name, value := range enumBuilder.ValueIds {
		values[name], err = writeValueIndex(fileId, count, value, outDir)
		if err != nil {
			return nil, fmt.Errorf("failed to write value index: %v", err)
		}
		count++
	}
	return &EnumIndex{
		Values: values,
		Kind:   "enum",
	}, nil
}

func writeValueIndex(fileId int, valueId int, ids []Ids, outDir string) (*EnumValue, error) {
	fileName := strconv.Itoa(fileId) + "-" + strconv.Itoa(valueId) + ".csv"
	filePath := filepath.Join(outDir, fileName)

	f, err := os.Create(filePath)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to create file: %v", err)

	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"dataRowId"})
	for _, record := range ids {
		row := []string{strconv.Itoa(record.DataRowId)}
		if err := w.Write(row); err != nil {
			return nil, fmt.Errorf("error writing record to file: %v", err)
		}
	}

	return &EnumValue{
		Count: len(ids),
		Url:   fileName,
	}, nil
}

func createIndexBuilder(property string, indexConfig IndexConfig) interface{} {
	switch indexConfig.Kind {
	case "enum":
		return EnumIndexBuilder{
			Property: property,
			Config:   indexConfig,
			ValueIds: make(map[string][]Ids),
		}
	}
	return nil
}

/**
 * Write indexes using the index builders and returns a `IndexRoot.indexes` map
 */
func writeIndexes(indexBuilders []interface{}, outDir string) (map[string]interface{}, error) {
	indexes := make(map[string]interface{})
	count := 0
	var err error
	for _, b := range indexBuilders {
		switch t := b.(type) {
		case EnumIndexBuilder:
			indexes[t.Property], err = t.WriteIndex(count, outDir)
			if err != nil {
				return nil, fmt.Errorf("failed to write index: %v", err)
			}
		default:
			continue
		}
		count++
	}

	return indexes, nil
}

/**
* Writes the data.csv file under `outDir` and returns its path.
 */
func writeResultsData(data []map[string]string, outDir string) (string, error) {
	fileName := "resultsData.csv"
	filePath := filepath.Join(outDir, fileName)

	f, err := os.Create(filePath)
	if err != nil {
		f.Close()
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	var keys []string
	for k := range data[0] {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	w.Write(keys)
	for _, record := range data {
		row := make([]string, 0, 1+len(keys))

		for _, k := range keys {
			row = append(row, record[k])
		}

		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("error writing record to file: %v", err)
		}
	}

	return fileName, nil
}

func writeIndexRoot(indexRoot IndexRoot, outDir string) error {
	fileName := "indexRoot.json"
	filePath := filepath.Join(outDir, fileName)

	b, err := json.Marshal(indexRoot)
	if err != nil {
		return fmt.Errorf("error while marshalling: %v", err)
	}

	err = os.WriteFile(filePath, b, 0644)
	if err != nil {
		return fmt.Errorf("error while writing the indexRoot: %v", err)
	}

	return nil
}

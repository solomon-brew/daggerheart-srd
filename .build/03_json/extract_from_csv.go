package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type regularField struct {
	orig string
	norm string
}

type groupedField struct {
	index int
	sub   string
	orig  string
}

func normalizeHeader(header string) string {
	header = strings.TrimSpace(strings.ToLower(header))
	var b strings.Builder
	underscore := false
	for _, r := range header {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			underscore = false
			continue
		}
		if !underscore {
			b.WriteRune('_')
			underscore = true
		}
	}
	out := strings.Trim(b.String(), "_")
	out = regexp.MustCompile(`_+`).ReplaceAllString(out, "_")
	return out
}

func groupDigitFields(fieldnames []string) ([]regularField, map[string][]groupedField) {
	regular := []regularField{}
	grouped := map[string][]groupedField{}
	pattern := regexp.MustCompile(`^(.*?)[ _-](\d+)[ _-](.*)$`)
	for _, orig := range fieldnames {
		m := pattern.FindStringSubmatch(orig)
		if m == nil {
			regular = append(regular, regularField{orig: orig, norm: normalizeHeader(orig)})
			continue
		}
		prefix, indexStr, subfield := m[1], m[2], m[3]
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			regular = append(regular, regularField{orig: orig, norm: normalizeHeader(orig)})
			continue
		}
		groupKey := normalizeHeader(prefix)
		subKey := normalizeHeader(subfield)
		if subKey == "" {
			subKey = "value"
		}
		grouped[groupKey] = append(grouped[groupKey], groupedField{
			index: index,
			sub:   subKey,
			orig:  orig,
		})
	}
	for key := range grouped {
		sort.Slice(grouped[key], func(i, j int) bool {
			return grouped[key][i].index < grouped[key][j].index
		})
	}
	return regular, grouped
}

func readCSV(path string) ([]string, [][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	rows := [][]string{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, row)
	}
	return header, rows, nil
}

func csvToJSON(csvPath, jsonPath string) error {
	header, rows, err := readCSV(csvPath)
	if err != nil {
		return err
	}
	if len(header) == 0 {
		return nil
	}
	regular, grouped := groupDigitFields(header)

	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		record := map[string]any{}
		// Regular fields
		for _, f := range regular {
			idx := indexOf(header, f.orig)
			if idx < 0 || idx >= len(row) {
				continue
			}
			val := strings.Trim(row[idx], " \t\r\n")
			if val == "" {
				continue
			}
			record[f.norm] = val
		}

		// Grouped fields
		for groupKey, fields := range grouped {
			groupRows := map[int]map[string]string{}
			for _, gf := range fields {
				idx := indexOf(header, gf.orig)
				if idx < 0 || idx >= len(row) {
					continue
				}
				val := strings.Trim(row[idx], " \t\r\n")
				if val == "" {
					continue
				}
				if _, ok := groupRows[gf.index]; !ok {
					groupRows[gf.index] = map[string]string{}
				}
				groupRows[gf.index][gf.sub] = val
			}
			if len(groupRows) == 0 {
				continue
			}
			indexes := make([]int, 0, len(groupRows))
			for idx := range groupRows {
				indexes = append(indexes, idx)
			}
			sort.Ints(indexes)
			groupList := make([]map[string]string, 0, len(indexes))
			for _, idx := range indexes {
				if len(groupRows[idx]) == 0 {
					continue
				}
				groupList = append(groupList, groupRows[idx])
			}
			if len(groupList) > 0 {
				record[groupKey] = groupList
			}
		}

		if len(record) > 0 {
			out = append(out, record)
		}
	}

	if err := os.MkdirAll(filepath.Dir(jsonPath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	data = append([]byte{0xEF, 0xBB, 0xBF}, data...)
	return os.WriteFile(jsonPath, data, 0644)
}

func indexOf(haystack []string, needle string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}

func main() {
	csvDir := ".build/02_csv"
	jsonDir := ".build/03_json"
	entries, err := os.ReadDir(csvDir)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", csvDir, err)
		return
	}
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		fmt.Printf("Error creating %s: %v\n", jsonDir, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".csv") {
			continue
		}
		csvPath := filepath.Join(csvDir, entry.Name())
		jsonPath := filepath.Join(jsonDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".json")
		fmt.Printf("Converting %s -> %s\n", csvPath, jsonPath)
		if err := csvToJSON(csvPath, jsonPath); err != nil {
			fmt.Printf("Error converting %s: %v\n", csvPath, err)
		}
	}
}

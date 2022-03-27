package golog

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ColorMatch struct {
	Text  string
	Color string

	c Color
}
type ColorFile struct {
	Rule []*ColorMatch
}

func (slf *ColorFile) ColorFromText(text string) Color {

	for _, rule := range slf.Rule {
		if strings.Contains(text, rule.Text) {
			return rule.c
		}
	}

	return NoColor
}

func (slf *ColorFile) Load(data string) error {

	err := json.Unmarshal([]byte(data), slf)
	if err != nil {
		return err
	}

	for _, rule := range slf.Rule {

		rule.c = matchColor(rule.Color)

		if rule.c == NoColor {
			return fmt.Errorf("color name not exists: %s", rule.Text)
		}

	}

	return nil
}

func NewColorFile() *ColorFile {
	return &ColorFile{}
}

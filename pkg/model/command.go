package model

type Command struct {
	Command string   `toml:"command"`
	Options []string `toml:"options"`
}

package oled

import (
	"os/exec"
	"fmt"
	"strings"
)

type direction string

const (
	Left               direction = "left"
	Right              direction = "right"
	DiagonalLeft       direction = "diagonal-left"
	DiagonalRight      direction = "diagonal-right"
	Rows               int       = 8
	CharactersPerRow   int       = 21
	PixelsPerRow       int       = 128
	PixelsPerCharacter int       = 8
	PixelsPerColumn    int       = Rows * PixelsPerCharacter
)

type Screen struct {
	commandQueue []string
}

func NewScreen() Screen {
	s := Screen{}
	s.InitializeScreen()
	s.ClearScreen()
	s.PowerOn()
	s.ResetCursor()
	s.ScrollOff()
	s.DimOff()

	return s
}

func (s *Screen) InitializeScreen() {
	s.runCommand("-i")
}

func (s *Screen) ClearScreen() {
	s.runCommand("-c")
}

func (s *Screen) ResetCursor() {
	s.MoveCursorToCharacter(0, 0)
}

func (s *Screen) MoveCursorToCharacter(row, character int) {
	s.addCommand("cursor", fmt.Sprintf("%d,%d", row, character))
}

func (s *Screen) MoveCursorToPixel(row, pixel int) {
	s.addCommand("cursorPixel", fmt.Sprintf("%d,%d", row, pixel))
}

func (s *Screen) PowerOn() {
	s.addCommand("power", "on")
}

func (s *Screen) PowerOff() {
	s.addCommand("power", "off")
}

func (s *Screen) InvertColorsOn() {
	s.addCommand("invert", "on")
}

func (s *Screen) InvertColorsOff() {
	s.addCommand("invert", "off")
}

func (s *Screen) DimOn() {
	s.addCommand("dim", "on")
}

func (s *Screen) DimOff() {
	s.addCommand("dim", "off")
}

func (s *Screen) WriteString(input string) {
	input = strings.Replace(input, "\n", "\\n", -1)
	s.addCommand("write", input)
}

func (s *Screen) WriteByte(input byte) {
	s.addCommand("writeByte", fmt.Sprintf("%#x", input))
}

func (s *Screen) ScrollOn(towards direction) {
	s.addCommand("scroll", string(towards))
}

func (s *Screen) ScrollOff() {
	s.addCommand("scroll", "stop")
}

func (s *Screen) DrawImage(path string) {
	s.addCommand("draw", path)
}

func (s *Screen) Commit() {
	cmd := exec.Command("oled-exp", s.commandQueue...)
	cmd.Run()

	s.ResetCommands()
}

func (s *Screen) ResetCommands() {
	s.commandQueue = []string{}
}

func (s *Screen) addCommand(arg ...string) {
	s.commandQueue = append(s.commandQueue, arg...)

}

func (s *Screen) runCommand(arg ...string) {
	cmd := exec.Command("oled-exp", arg...)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

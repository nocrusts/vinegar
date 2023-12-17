// Package wine implements wine program command routines for
// interacting with a wineprefix [Prefix]
package wine

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The program used for Wine.
var Wine = "wine"

// Prefix is a representation of a wineprefix, which is where
// WINE stores its data, which is equivalent to WINE's C:\ drive.
type Prefix struct {
	// Output specifies the descendant prefix commmand's
	// Stderr and Stdout together.
	//
	// Wine will always write to stderr instead of stdout,
	// Stdout is combined just to make that certain.
	Output io.Writer

	dir string
}

// Wine Vulkan information stored inside winevulkan.json.
type WineVulkan struct {
	FileFormatVersion string `json:"file_format_version"`
	ICD               ICD    `json:"ICD"`
}

// Installable Client Driver info containing the library path and API version.
type ICD struct {
	LibraryPath string `json:"library_path"`
	ApiVersion  string `json:"api_version"`
}

// New returns a new Prefix.
func New(dir string, out io.Writer) Prefix {
	return Prefix{
		Output: out,
		dir:    dir,
	}
}

// WineLook checks for [Wine] with exec.LookPath, and returns
// true if [Wine] is present and has no problems.
func WineLook() bool {
	_, err := exec.LookPath(Wine)
	return err == nil
}

// Dir retrieves the Prefix's directory.
func (p *Prefix) Dir() string {
	return p.dir
}

// Get supported vulkan version.
func (p *Prefix) VulkanVersion() string {
	winevk_info := filepath.Join(p.Dir(), "drive_c", "windows", "syswow64", "winevulkan.json")

	f, err := os.ReadFile(winevk_info)
	if err != nil {
		return ""
	}

	wv := WineVulkan{}
	if err := json.Unmarshal(f, &wv); err != nil {
		return ""
	}

	return wv.ICD.ApiVersion
}

func (p *Prefix) VulkanSupported() bool {
	vulkan_version := p.VulkanVersion()

	if vulkan_version == "" {
		return false
	}

	return strings.Split(vulkan_version, ".")[1] >= "1"
}

// Wine makes a new Cmd with wine as the named program.
func (p *Prefix) Wine(exe string, arg ...string) *Cmd {
	arg = append([]string{exe}, arg...)

	return p.Command(Wine, arg...)
}

// Kill runs Command with 'wineserver -k' as the named program.
func (p *Prefix) Kill() {
	log.Println("Killing wineprefix")

	_ = p.Command("wineserver", "-k").Run()
}

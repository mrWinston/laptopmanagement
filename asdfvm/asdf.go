package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// PluginState represents the different states that can be specified for a given asdf plugin
type PluginState int

const (
  // UNKNOWN represents an invalid or unknown state
	UNKNOWN PluginState = iota
  // PRESENT ensures that the plugin is installed in any version
	PRESENT
  // ABSENT makes sure that the plugin isn't installed
	ABSENT
  // LATEST ensures that the plugin added and has the latest version installed
	LATEST
)

// ParsePluginState parses a state name into a PluginState enum
func ParsePluginState(name string) PluginState {
	switch strings.ToLower(name) {
	case "present":
		return PRESENT
	case "absent":
		return ABSENT
	case "latest":
		return LATEST
	default:
		return UNKNOWN
	}
}

// EnsureAsdfPlugin makes sure that the specified asdf plugin has the state specified
func EnsureAsdfPlugin(name string, url string, stateName string, version string, isDefault bool) (changed bool, err error) {
	desiredState := ParsePluginState(stateName)

	switch desiredState {
	case PRESENT, LATEST:
		return EnsureAsdfPluginInstalled(name, url, version, isDefault)
	case ABSENT:
		return EnsureAsdfPluginRemoved(name)
  default:
		return false, NewInvalidStateError(stateName)
	}
}

// EnsureAsdfPluginRemoved removes the specified asdf plugin. returns true only if something was changed.
func EnsureAsdfPluginRemoved(name string) (changes bool, err error) {
  isInstalled, err := IsAlreadyInstalled(name)
  if err != nil {
    return false, err
  }
  if !isInstalled {
    return false, nil
  }
  
	rmPluginCmd := exec.Command("asdf", "plugin", "remove", name)
	rawOut, err := rmPluginCmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("Error running asdf: %w. Output: %s", err, rawOut)
	}
	return true, nil
}

// EnsureAsdfPluginInstalled adds the given asdf plugin and installs the program in the given version.
func EnsureAsdfPluginInstalled(name string, url string, version string, isDefault bool) (changed bool, err error) {
	isInstalled, err := IsAlreadyInstalled(name)
	if err != nil {
		return false, err
	}

	if isInstalled {
		return false, nil
	}

	err = AddAsdfPlugin(name, url)
	if err != nil {
		return true, err
	}
	err = InstallAsdfPlugin(name, version)
	if err != nil {
		return true, err
	}
  
  if isDefault {
    err = SetAsdfPluginVersionGlobal(name, "latest")
  }
	return true, err
}

// IsAlreadyInstalled returns true, if and only if an asdf plugin with the given name is already added.
func IsAlreadyInstalled(name string) (bool, error) {
	asdfCmd := exec.Command("asdf", "plugin", "list")
	rawOut, err := asdfCmd.CombinedOutput()
	if err != nil {
		return true, err
	}
	out := string(rawOut)
	pluginNames := strings.Split(out, "\n")
	for _, plugin := range pluginNames {
		if plugin == name {
			return true, nil
		}
	}
	return false, nil
}

// AddAsdfPlugin adds the asdf plugin with the given name and url
func AddAsdfPlugin(name string, url string) error {
	addPluginCmd := exec.Command("asdf", "plugin", "add", name, url)
	rawOut, err := addPluginCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running asdf: %w. Output: %s", err, rawOut)
	}
	return nil
}

// InstallAsdfPlugin installs the specified version of the program managed by the specified asdf plugin name
func InstallAsdfPlugin(name string, version string) error {
	addPluginCmd := exec.Command("asdf", "install", name, version)
	rawOut, err := addPluginCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running asdf: %w. Output: %s", err, rawOut)
	}
	return nil
}

// SetAsdfPluginVersionGlobal sets the global version for a specified plugin
func SetAsdfPluginVersionGlobal(name string, version string) error {
	addPluginCmd := exec.Command("asdf", "global", name, version)
	rawOut, err := addPluginCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running asdf: %w. Output: %s", err, rawOut)
	}
	return nil
}

package main

import "github.com/baulk/bloat/modules/semver"

type ScoopLink struct {
	URL  string `json:"url,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type ScoopArchitecture struct {
	X64 *ScoopLink `json:"64bit,omitempty"`
	X86 *ScoopLink `json:"32bit,omitempty"`
}

type ScoopMetadata struct {
	Version      string             `json:"version"`
	Description  string             `json:"description,omitempty"`
	Homepage     string             `json:"homepage,omitempty"`
	Depends      string             `json:"depends,omitempty"`
	URL          string             `json:"url,omitempty"`
	Hash         string             `json:"hash,omitempty"`
	Architecture *ScoopArchitecture `json:"architecture,omitempty"`
}

type ArchitectureLink struct {
	URL       string   `json:"url,omitempty"`
	Hash      string   `json:"hash,omitempty"`
	Links     []string `json:"links,omitempty"`
	Launchers []string `json:"launchers,omitempty"`
}

type Architecture struct {
	Win64    *ArchitectureLink `json:"64bit,omitempty"`
	Win32    *ArchitectureLink `json:"32bit,omitempty"`
	WinARM64 *ArchitectureLink `json:"arm64,omitempty"`
}

type Environment struct {
	Category     string   `json:"category"`
	Paths        []string `json:"path,omitempty"`
	Includes     []string `json:"include,omitempty"`
	Libs         []string `json:"lib,omitempty"`
	NewDirs      []string `json:"mkdir,omitempty"`
	Env          []string `json:"env,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type Metadata struct {
	Version        string        `json:"version"`
	Description    string        `json:"description,omitempty"`
	Homepage       string        `json:"homepage,omitempty"`
	Extension      string        `json:"extension,omitempty"`
	Notes          string        `json:"notes,omitempty"`
	License        string        `json:"license,omitempty"`
	URL            string        `json:"url,omitempty"`
	Hash           string        `json:"url.hash,omitempty"`
	X64URL         string        `json:"url64,omitempty"`
	X64Hash        string        `json:"url64.hash,omitempty"`
	Arm64URL       string        `json:"urlarm64,omitempty"`
	Arm64Hash      string        `json:"urlarm64.hash,omitempty"`
	Links          []string      `json:"links,omitempty"`
	Launchers      []string      `json:"launchers,omitempty"`
	Links64        []string      `json:"links64,omitempty"`
	Launchers64    []string      `json:"launchers64,omitempty"`
	LinksArm64     []string      `json:"linksarm64,omitempty"`
	LaunchersArm64 []string      `json:"launchersarm64,omitempty"`
	Architecture   *Architecture `json:"architecture,omitempty"`
	Suggests       string        `json:"suggest,omitempty"`
	ForceDelete    []string      `json:"force_delete,omitempty"`
	VEnv           *Environment  `json:"venv,omitempty"`
}

func (m *Metadata) Compare(s *ScoopMetadata) int {
	return semver.Compare(m.Version, s.Version)
}

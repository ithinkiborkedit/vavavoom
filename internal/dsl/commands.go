package dsl

import (
	"nelly/internal/dslcore"
	"nelly/internal/gitcmd"
)

var Commands = map[string]dslcore.CommandFunc{}

func RegisterDefaultCommands() {
	Commands["clone"] = gitcmd.CloneCommand
	Commands["directory"] = nil
	Commands["init"] = gitcmd.InitCommand
	Commands["createBranch"] = gitcmd.CreateBranchCommand
	Commands["track"] = gitcmd.TrackCommand
}

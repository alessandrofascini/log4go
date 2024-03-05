package config

import "os"

const (
	DefaultChannelSize         = 160
	DefaultInternalChannelSize = 110
)

var defaultAppenderChannelSize = DefaultChannelSize

func SetDefaultAppenderChannelSize(newValue int) {
	defaultAppenderChannelSize = newValue
}

func GetDefaultAppenderChannelSize() int {
	return defaultAppenderChannelSize
}

var defaultInternalChannelSize = DefaultInternalChannelSize

func GetDefaultInternalChannelSize() int {
	return defaultInternalChannelSize
}

func SetDefaultInternalChannelSize(newValue int) {
	defaultInternalChannelSize = newValue
}

var createFolderPerm os.FileMode = 0755

func SetCreateFolderPerm(perm os.FileMode) {
	createFolderPerm = perm
}

func GetCreateFolderPerm() os.FileMode {
	return createFolderPerm
}

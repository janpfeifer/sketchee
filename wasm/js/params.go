package js

// Query parameters that can be added to the URL.
const (
	// Q_VIDEO can be set to 0 to skip opening cut-scene.
	Q_VIDEO = "video"

	// Q_AUDIO can be set to 0 for no audio.
	Q_AUDIO = "audio"

	// Q_LOGIN is the default user to login for auto-login.
	Q_LOGIN = "login"

	// Q_JOIN holds the game name to auto-join.
	Q_JOIN = "join"

	// Q_AUTOSTART triggers automatically starting the gamed joined.
	Q_AUTOSTART = "autostart"

	// Q_ANIMATION_STEP defines the number of milliseconds between a
	// refresh/animation step. It defines the (inverse of) the target
	// refresh rate.
	Q_ANIMATION_STEP = "animation_step"

	// Q_VMODULE defines the debugging level per module in glog.
	// It's equivalent to setting the --vmodule=... flag for glog.
	Q_VMODULE = "vmodule"

	// Q_PMODULE defines the list of files to enable profiling.
	Q_PMODULE = "pmodule"

	// Q_BASE_TEAM defines a number to be added to the team number when selecting the art. Used for debugging the art
	// for the various team numbers.
	Q_BASE_TEAM = "baseteam"
)

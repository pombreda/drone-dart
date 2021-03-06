package dart

// SDK represents the latest Dart version available
// for download, generated by Dart's Buildbot.
type SDK struct {
	Revision string `json:"revision"`
	Version  string `json:"version"`
	Date     string `json:"date"`
}

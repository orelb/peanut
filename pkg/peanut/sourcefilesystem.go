package peanut

// A SourceFilesystem is used to interact with source's files.
type SourceFilesystem interface {
	FetchAll(destDir string) error
}

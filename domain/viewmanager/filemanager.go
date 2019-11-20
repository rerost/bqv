package viewmanager

type FileManager interface {
	ViewReadWriter
}

func NewFileManager(dir string) FileManager {
	// TODO
	return nil
}

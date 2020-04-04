package peanut

import (
	"io"
)

//func CopyByMapping(sourceDir string, destinationDir string, mapping MappingRule) error {
//	sourcePath := filepath.Join(sourceDir, mapping.source)
//	destPath := filepath.Join(destinationDir, mapping.destination)
//
//	sourceInfo, err := fs.Stat(sourcePath)
//	if os.IsNotExist(err) {
//		errorMessage := fmt.Sprintf("File %s not found", sourcePath)
//		return errors.New(errorMessage)
//	}
//
//	if sourceInfo.Mode().IsDir() {
//		err = copyDirectory(sourcePath, destPath)
//		if err != nil {
//			return err
//		}
//	} else if sourceInfo.Mode().IsRegular() {
//		// if the source mapping is a file, copy to the destination path
//
//		destStat, err := fs.Stat(destPath)
//
//		// If the destination is a directory, set the destination path for our original filename inside the destination directory.
//		if err == nil && destStat.IsDir() {
//			destPath = filepath.Join(destPath, sourceInfo.Name())
//		}
//
//		err = copyFile(sourcePath, destPath)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func copyFile(sourcePath string, destinationPath string) error {
	destination, err := fs.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	source, err := fs.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}

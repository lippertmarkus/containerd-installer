package http

import (
	"archive/tar"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jfbus/httprs"
	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func DownloadAndExtract(url string, targetPath string, skipLevel int) error {
	logrus.Debugf("Downloading from %s", url)
	response, err := retryablehttp.Get(url)
	if err != nil {
		return errors.Wrapf(err, "Can't download %s", url)
	}
	if response.StatusCode != http.StatusOK {
		return errors.Errorf("Downloading failed with status %d %s", response.StatusCode, response.Status)
	}

	httpReadSeeker := httprs.NewHttpReadSeeker(response)
	defer httpReadSeeker.Close()

	if err := extract(path.Base(url), httpReadSeeker, response.ContentLength, targetPath, skipLevel); err != nil {
		return errors.WithMessage(err, "Unarchiving failed")
	}

	return nil
}

func extract(filename string, archiveFile io.Reader, archiveSize int64, targetDir string, skipLevel int) error {
	specificArchiver, err := archiver.ByExtension(filename)
	if err != nil {
		return err
	}

	archive := specificArchiver.(archiver.Reader)
	if err := archive.Open(archiveFile, archiveSize); err != nil {
		return err
	}
	defer archive.Close()

	for {
		file, err := archive.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := extractFile(file, targetDir, skipLevel); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file archiver.File, targetDir string, skipLevel int) error {
	defer file.Close()

	fileFullname, err := getFileFullNameFromHeader(file)
	if err != nil {
		return err
	}
	splittedArchivePath := strings.SplitN(fileFullname, "/", 1+skipLevel)
	targetFile := filepath.Join(targetDir, splittedArchivePath[len(splittedArchivePath)-1])

	if file.IsDir() {
		if err := os.MkdirAll(targetFile, file.Mode()); err != nil {
			return err
		}
	} else {
		targetFile, err := os.Create(targetFile)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, file); err != nil {
			return err
		}
	}

	return nil
}

func getFileFullNameFromHeader(file archiver.File) (string, error) {
	switch file.Header.(type) {
	case zip.FileHeader:
		return file.Header.(zip.FileHeader).Name, nil
	case *tar.Header:
		return file.Header.(*tar.Header).Name, nil
	default:
		return "", errors.New("unsupported header type")
	}
}

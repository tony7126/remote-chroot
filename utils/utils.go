package utils

import (
	"math/rand"
	"time"
	"fmt"
	"os"
	"io"
	"archive/tar"

)

const (
	RAND_DIR_NAME_LIMIT = 10000
)

// create directory with random name that doesn't exist
func CreateRandDest() (dest string) {
	dirExists := true
	for dirExists {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		n := r1.Intn(RAND_DIR_NAME_LIMIT)
		dest = fmt.Sprintf("%s/gs_%d", "/tmp", n)
		dirExists, _ = Exists(dest)
	}

	return
}

// Exists returns whether the given file or directory Exists or not
func Exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func Untar(r io.Reader, dest string) (err error) {
	tarBallReader := tar.NewReader(r)

     // Extracting tarred files

     for {
     		var header *tar.Header
             header, err = tarBallReader.Next()
             //fmt.Println(header, err)
             if err != nil {
                     if err == io.EOF {
                        err = nil
                     }
                     return
             }

             // get the individual filename and extract to the current directory
             filename := dest + "/" + header.Name
             switch header.Typeflag {
             case tar.TypeDir:
                     // handle directory
                     //fmt.Println("Creating directory :", filename)
                     err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

                     if err != nil {
                             fmt.Println(err)
                             return
                     }

             case tar.TypeReg:
                     // handle normal file
                     //fmt.Println("Untarring :", filename)
                     var writer *os.File
                     writer, err = os.Create(filename)

                     if err != nil {
                             fmt.Println(err)
                             return
                     }

                     io.Copy(writer, tarBallReader)

                     err = os.Chmod(filename, os.FileMode(header.Mode))

                     if err != nil {
                             fmt.Println(err)
                             return
                     }

                     writer.Close()

             case tar.TypeSymlink:
                err = os.Symlink(header.Linkname, filename)
                if err != nil {
                    fmt.Println(err)
                    return
                }
             default:
                     fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
             }
     }

     return

}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
    sfi, err := os.Stat(src)
    if err != nil {
        return
    }
    if !sfi.Mode().IsRegular() {
        // cannot copy non-regular files (e.g., directories,
        // symlinks, devices, etc.)
        return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
    }
    dfi, err := os.Stat(dst)
    if err != nil {
        if !os.IsNotExist(err) {
            return
        }
    } else {
        if !(dfi.Mode().IsRegular()) {
            return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
        }
        if os.SameFile(sfi, dfi) {
            return
        }
    }
    if err = os.Link(src, dst); err == nil {
        return
    }
    err = copyFileContents(src, dst)
    return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return
    }
    err = out.Sync()
    return
}
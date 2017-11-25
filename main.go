package main

import (
  "github.com/urfave/cli"
  "os"
  "math"
  "crypto/md5"
  "io"
  "encoding/hex"
  "fmt"
  "path"
  "errors"
)

/**
Calculate md5 of a file
 */
func calcMd5(filePath string) (string) {
  const fileChunk = 8192
  file, err := os.Open(filePath)

  if err != nil {
    panic(err.Error())
  }

  defer file.Close()

  // calculate the file size
  info, _ := file.Stat()

  fileSize := info.Size()

  blocks := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

  hash := md5.New()

  for i := uint64(0); i < blocks; i++ {
    blockSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
    buf := make([]byte, blockSize)

    file.Read(buf)
    io.WriteString(hash, string(buf)) // append into the hash
  }

  md5string := hex.EncodeToString(hash.Sum(nil))

  return md5string
}

func main() {
  app := cli.NewApp()

  app.Name = "md5-rename"
  app.Usage = "md5-rename" + " [path]"
  app.Version = "0.0.1"
  app.Description = "Rename the file with MD5 value"

  app.Action = func(c *cli.Context) (err error) {
    defer func() {
      if err != nil {
        fmt.Println(err.Error())
      }
    }()
    var (
      cwd      string
      fileStat os.FileInfo
      filePath = c.Args().Get(0)
      distPath string
    )

    if filePath == "" {
      err = errors.New("File name must be required!")
    }

    if cwd, err = os.Getwd(); err != nil {
      return
    }

    if path.IsAbs(filePath) {
      distPath = filePath
    } else {
      distPath = path.Join(cwd, filePath)
    }

    if fileStat, err = os.Stat(distPath); err != nil {
      return
    }

    if fileStat.IsDir() {
      err = errors.New("Can not rename a dir")
      return
    }

    md5 := calcMd5(distPath)

    newFilePath := path.Join(path.Dir(distPath), md5+path.Ext(distPath))

    if err = os.Rename(distPath, newFilePath); err != nil {
      return
    }

    fmt.Printf("%v > %v", distPath, newFilePath)
    return
  }

  app.Run(os.Args)
}

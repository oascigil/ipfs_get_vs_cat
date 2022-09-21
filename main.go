package main

import (
 "fmt"
 "log"
 "os"
 "io"
 "time"
 ipfsapi "github.com/ipfs/go-ipfs-api"
 "crypto/sha256"
)

func resolveUsingGet(shell *ipfsapi.Shell, cid string, outDir string) (error) {
    err := shell.Get(cid, outDir)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func resolveUsingCat(shell *ipfsapi.Shell, cid string) (io.ReadCloser, error) {
    f, err := shell.Cat(cid)
    if err != nil {
        log.Println(err)
        return nil, err
    }       
    return f, nil 
}        

func main() {
    cid := "QmTz3oc4gdpRMKP2sdGUPZTAGRngqjsi99BPoztyP53JMM"
    shell := ipfsapi.NewShell("localhost:5001")
    shell.SetTimeout(time.Duration(2 * time.Minute))
    err := resolveUsingGet(shell, cid, "./get/")
    if err != nil {
        log.Fatal(err)
    }

    res, err := resolveUsingCat(shell, cid)
    outFile, err := os.Create("./cat/" + cid)
    if err!= nil{
        log.Fatal(err)
    }
    defer outFile.Close()
    _, err = io.Copy(outFile, res)

    // Compute hash using io.Copy
    sh := sha256.New()
    if _, err := io.Copy(sh, res); err != nil {
        log.Println(err)
	}
    infile, err := os.Open("./cat/" + cid)
    if err!=nil {
        log.Fatal(err)
    }
    fmt.Printf("Hash computed using io.Copy: %x\n", sh.Sum(nil)[:])
    // Compute hash on the file
    sha_hasher := sha256.New()
    if _, err := io.Copy(sha_hasher, infile); err != nil {
        log.Fatal(err)
    }
    infile.Close()
    fmt.Printf("Hash computed on the file: %x\n", sha_hasher.Sum(nil)[:])
}

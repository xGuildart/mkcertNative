package main

import (
	"aouimar/flagf"
	"aouimar/mkcert"
	"flag"
	"fmt"
	"log"
	"runtime/debug"
)

import "C"

func main() {
	str := "127.0.0.1 example.com ::1"
	MkcertCommand(C.CString(str))
}

//export MkcertCommand
func MkcertCommand(arg *C.char) {
	var m *flagf.FFlag
	m = new(flagf.FFlag)
	m.Init(make([]string, 0, 10), make(map[string]*string), make(map[string]*bool), make(map[string]string))

	var (
		installFlag   = m.Bool("install", false, "")
		uninstallFlag = m.Bool("uninstall", false, "")
		pkcs12Flag    = m.Bool("pkcs12", false, "")
		ecdsaFlag     = m.Bool("ecdsa", false, "")
		clientFlag    = m.Bool("client", false, "")
		helpFlag      = m.Bool("help", false, "")
		carootFlag    = m.Bool("CAROOT", false, "")
		csrFlag       = m.String("csr", "", "")
		certFileFlag  = m.String("cert-file", "", "")
		keyFileFlag   = m.String("key-file", "", "")
		p12FileFlag   = m.String("p12-file", "", "")
		versionFlag   = m.Bool("version", false, "")
	)
	m.Parse(C.GoString(arg))

	if *helpFlag {
		fmt.Print(mkcert.ShortUsage)
		fmt.Print(mkcert.AdvancedUsage)
		return
	}
	if *versionFlag {
		if mkcert.Version != "" {
			fmt.Println(mkcert.Version)
			return
		}
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(buildInfo.Main.Version)
			return
		}
		fmt.Println("(unknown)")
		return
	}
	if *carootFlag {
		if *installFlag || *uninstallFlag {
			log.Fatalln("ERROR: you can't set -[un]install and -CAROOT at the same time")
		}
		fmt.Println(mkcert.GetCAROOT())
		return
	}
	if *installFlag && *uninstallFlag {
		log.Fatalln("ERROR: you can't set -install and -uninstall at the same time")
	}
	if *csrFlag != "" && (*pkcs12Flag || *ecdsaFlag || *clientFlag) {
		log.Fatalln("ERROR: can only combine -csr with -install and -cert-file")
	}
	if *csrFlag != "" && flag.NArg() != 0 {
		log.Fatalln("ERROR: can't specify extra arguments when using -csr")
	}

	(&mkcert.Mkcert{
		InstallMode: *installFlag, UninstallMode: *uninstallFlag, CsrPath: *csrFlag,
		Pkcs12: *pkcs12Flag, Ecdsa: *ecdsaFlag, Client: *clientFlag,
		CertFile: *certFileFlag, KeyFile: *keyFileFlag, P12File: *p12FileFlag,
	}).Run(m.Args())
}

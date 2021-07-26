package main

import (
	"bufio"
	"log"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func main() {
	wmname := ""
	xconn, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer xconn.Close()
	xroot := xproto.Setup(xconn).DefaultScreen(xconn).Root
	defer xproto.ChangeProperty(xconn, xproto.PropModeReplace, xroot, xproto.AtomWmName,
		xproto.AtomString, 8, uint32(len(wmname)), []byte(wmname))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			if len(scanner.Text()) == 0 {
				break
			}
			wmname = " " + scanner.Text() + " "
			xproto.ChangeProperty(xconn, xproto.PropModeReplace, xroot, xproto.AtomWmName,
				xproto.AtomString, 8, uint32(len(wmname)), []byte(wmname))
		}
		sigs <- os.Interrupt
	}()
	<-sigs
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
}

package network

import (
	"errors"
	"log"
	"net"
)

// Telnet IAC constants
const (
	SE   = byte(240)
	NOP  = byte(241)
	BRK  = byte(243)
	IP   = byte(244)
	AO   = byte(245)
	AYT  = byte(246)
	EC   = byte(247)
	EL   = byte(248)
	GA   = byte(249)
	SB   = byte(250)
	WILL = byte(251)
	WONT = byte(252)
	DO   = byte(253)
	DONT = byte(254)
	IAC  = byte(255)
)

const Escape = byte('\033')

// Telnet Options
const (
	ECHO     = byte(1)
	TTYPE    = byte(24)
	NAWS     = byte(31)
	LINEMODE = byte(34)
	EOR      = byte(239)
)

// RequestTerminalType sends the terminal type request
func RequestTerminalType(conn net.Conn) (string, error) {
	conn.Write([]byte{IAC, DO, 24})
	conn.Write([]byte{IAC, SB, 24, 1, IAC, SE})

	var buf []byte
	// Read until SE is read
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			return "", err
		}
		if b[0] != IAC && b[0] != SE && b[0] != WILL && b[0] != WONT && b[0] != DO && b[0] != DONT && b[0] != SB && b[0] != NAWS {
			//log.Println(b[0])
			buf = append(buf[:], b[0])
		}

		if b[0] == SE {
			return string(buf), nil
		}
	}
}

// RequestTerminalSize requests the terminal size
func RequestTerminalSize(conn net.Conn) (height, width int, err error) {
	conn.Write([]byte{IAC, DO, 31})
	conn.Write([]byte{IAC, SB, 31, 1, IAC, SE})

	var buf []byte
	// Read until SE is read
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			log.Println(err)
			return 0, 0, err
		}
		buf = append(buf[:], b[0])
		if b[0] == SE {
			if len(buf) == 12 {
				height = int(buf[9])
				width = int(buf[7])
			} else {
				return 0, 0, errors.New("client sent invalid terminal size response")
			}
			break
		}
	}
	return height, width, nil
}

// DisableEcho disables echo
func DisableEcho(conn net.Conn) error {
	conn.Write([]byte{IAC, WILL, ECHO})
	// We don't care about the response
	// Read until end of response
	var iac bool
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			//log.Println(err)
			return err
		}

		if b[0] == IAC {
			iac = true
		}
		if iac {
			if b[0] == WONT {
				// read next byte
				b := make([]byte, 1)
				_, err := conn.Read(b)
				if err != nil {
					//log.Println(err)
					return err
				}
				if b[0] == ECHO {
					return errors.New("client refused echo")
				}

			}
			if b[0] == DO {
				// read next byte
				b := make([]byte, 1)
				_, err := conn.Read(b)
				if err != nil {
					//log.Println(err)
					return err
				}
				if b[0] == ECHO {
					//log.Println("client accepted echo")
					return nil
				}
			}
		}
	}
}

// EnableLineMode disables line mode
func EnableLineMode(conn net.Conn) error {
	conn.Write([]byte{IAC, DO, LINEMODE})
	// We don't care about the response
	// Read until end of response
	var iac bool
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			log.Println(err)
			return err
		}

		if b[0] == IAC {
			iac = true
		}
		if iac {
			if b[0] == SB {
				//log.Println("Subnegotiation begun")
				// Read until SE is read because we don't care about the response
				for {
					b := make([]byte, 1)
					_, err := conn.Read(b)
					if err != nil {
						log.Println(err)
						return err
					}
					if b[0] == SE {
						//log.Println("Subnegotiation ended")
						return nil // Client accented linemode, and negotiated
					}
				}
				iac = false
			}
			if b[0] == WONT {
				// read next byte
				b := make([]byte, 1)
				_, err := conn.Read(b)
				if err != nil {
					//log.Println(err)
					return err
				}
				if b[0] == LINEMODE {
					return errors.New("client refused linemode")
				}
			}
			if b[0] == WILL {
				// read next byte
				b := make([]byte, 1)
				_, err := conn.Read(b)
				if err != nil {
					log.Println(err)
					return err
				}
				if b[0] == LINEMODE {
					continue
				}
			}
		}
	}
}

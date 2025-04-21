
package protocol

import (
	"bufio"
	"errors"
	"net"
	"strconv"
)

type Parser struct {
	conn net.Conn
	r    *bufio.Reader
	line []byte
	pos  int
}

type Command struct {
	Args []string
	Conn net.Conn
}

func NewParser(conn net.Conn) *Parser {
	return &Parser{
		conn: conn,
		r:    bufio.NewReader(conn),
		line: make([]byte, 0),
		pos:  0,
	}
}

func (p *Parser) Command() (Command, error) {
	b, err := p.r.ReadByte()
	if err != nil {
		return Command{}, err
	}
	if b == '*' {
		return p.respArray()
	} else {
		line, err := p.readLine()
		if err != nil {
			return Command{}, err
		}
		p.pos = 0
		p.line = append([]byte{b}, line...)
		return p.inline()
	}
}

func (p *Parser) readLine() ([]byte, error) {
	line, err := p.r.ReadBytes('\r')
	if err != nil {
		return nil, err
	}
	_, err = p.r.ReadByte() // Consume '\n'
	if err != nil {
		return nil, err
	}
	return line[:len(line)-1], nil
}

func (p *Parser) inline() (Command, error) {
	cmd := Command{Conn: p.conn}
	for p.pos < len(p.line) {
		arg, err := p.consumeArg()
		if err != nil {
			return cmd, err
		}
		if arg != "" {
			cmd.Args = append(cmd.Args, arg)
		}
	}
	return cmd, nil
}

func (p *Parser) consumeArg() (string, error) {
	for p.pos < len(p.line) && p.line[p.pos] == ' ' {
		p.pos++
	}
	start := p.pos
	for p.pos < len(p.line) && p.line[p.pos] != ' ' && p.line[p.pos] != '\r' {
		p.pos++
	}
	return string(p.line[start:p.pos]), nil
}


func (p *Parser) respArray() (Command, error) {
	cmd := Command{Conn: p.conn}
	szLine, err := p.readLine()
	if err != nil {
		return cmd, err
	}
	sz, _ := strconv.Atoi(string(szLine))
	for i := 0; i < sz; i++ {
		t, err := p.r.ReadByte()
		if err != nil {
			return cmd, err
		}
		if t != '$' {
			return cmd, errors.New("expected bulk string")
		}
		szLine, err := p.readLine()
		if err != nil {
			return cmd, err
		}
		sz, _ := strconv.Atoi(string(szLine))
		buf := make([]byte, sz+2)
		_, err = p.r.Read(buf)
		if err != nil {
			return cmd, err
		}
		cmd.Args = append(cmd.Args, string(buf[:sz]))
	}
	return cmd, nil
}

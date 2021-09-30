package postcss

import (
	"bytes"
	"strings"
	"text/scanner"
)

func Parse(rule string) Node {
	return ParseBytes([]byte(rule))
}

func ParseBytes(rule []byte) Node {
	s := &NodeScanner{}
	s.Scanner.Init(bytes.NewBuffer(rule))
	s.Scanner.Error = func(s *scanner.Scanner, msg string) {
		//fmt.Println(msg)
		// todo
	}
	return s.ScanNode()
}

type NodeAppendable interface {
	AppendNode(n Node)
}

type NodeScanner struct {
	Scanner scanner.Scanner
	temp    *bytes.Buffer
	stack   []NodeAppendable
	decl    *Declaration
}

func (s *NodeScanner) resetTemp() {
	s.temp = bytes.NewBuffer(nil)
}

func (s *NodeScanner) TempText() string {
	text := s.temp.String()
	s.resetTemp()
	return strings.TrimSpace(text)
}

func (s *NodeScanner) appendNode(n Node) {
	s.stack[len(s.stack)-1].AppendNode(n)

	if nodeAppendable, ok := n.(NodeAppendable); ok {
		s.stack = append(s.stack, nodeAppendable)
	}
}

func (s *NodeScanner) closeRule() {
	s.stack = s.stack[0 : len(s.stack)-1]
}

func (s *NodeScanner) openRule() {
	selector := s.TempText()

	if s.decl != nil {
		selector = s.decl.Prop + ":" + selector
		s.decl = nil
	}

	if selector[0] == '@' {
		parts := strings.SplitN(selector[1:], " ", 2)

		s.appendNode(&AtRule{
			Name: parts[0],
			Params: func() string {
				if len(parts) > 1 {
					return parts[1]
				}
				return ""
			}(),
		})
	} else {
		s.appendNode(&Rule{
			Selector: selector,
		})
	}
}

func (s *NodeScanner) openDecl() {
	s.decl = &Declaration{
		Prop: s.TempText(),
	}
}

func (s *NodeScanner) closeDeclIfNeed() {
	if s.decl != nil {
		s.decl.Value = s.TempText()
		s.appendNode(&*s.decl)
		s.decl = nil
	}
}

func (s *NodeScanner) ScanQuote() {
	s.Scanner.Scan()
	s.temp.WriteString(s.Scanner.TokenText())
}

func (s *NodeScanner) ScanParenthesis() rune {
	left := s.Scanner.Peek()

	tok := s.Scanner.Next()
	s.temp.WriteRune(tok)

	for tok != scanner.EOF {
		tok = s.Scanner.Peek()

		switch tok {
		case '"', '\'':
			s.ScanQuote()
			tok = s.Scanner.Peek()
		}

		if left == '(' && tok == ')' {
			s.temp.WriteRune(s.Scanner.Next())
			break
		}

		if left == '[' && tok == ']' {
			s.temp.WriteRune(s.Scanner.Next())
			break
		}

		s.temp.WriteRune(s.Scanner.Next())
	}

	return tok
}

func (s *NodeScanner) ScanNode() Node {
	root := &Root{}

	s.stack = []NodeAppendable{root}
	s.resetTemp()

	tok := s.Scanner.Peek()

	for tok != scanner.EOF {
		switch tok {
		case '"', '\'':
			tok = s.Scanner.Scan()
			s.temp.WriteString(s.Scanner.TokenText())
		case '\t', '\n', '\r':
			// ignore
			tok = s.Scanner.Next()
		case '(', '[':
			s.ScanParenthesis()
		case '&':
			s.temp.WriteString(HOLDER)
			s.Scanner.Next()
		case ':':
			s.openDecl()
			s.Scanner.Next()

			if s.decl.IsVariable() {
				for tok != scanner.EOF {
					tok = s.Scanner.Peek()
					if tok == '(' || tok == '[' {
						s.ScanParenthesis()
						tok = s.Scanner.Peek()
					}
					if tok == ';' {
						break
					}
					s.temp.WriteRune(s.Scanner.Next())
				}
				s.closeDeclIfNeed()
				s.Scanner.Next()
			}
		case ';':
			s.closeDeclIfNeed()
			s.Scanner.Next()
		case '{':
			s.openRule()
			s.Scanner.Next()
		case '}':
			s.closeDeclIfNeed()
			s.closeRule()
			s.Scanner.Next()
		default:
			s.temp.WriteRune(s.Scanner.Next())
		}

		tok = s.Scanner.Peek()
	}

	return root
}

// Copyright © 2021 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lexer

import "unicode"

func (l *lexer) consumeSpace() {
	for unicode.IsSpace(l.peek()) {
		l.consume()
	}

	l.ignore()
}

func (l *lexer) consumeWord() {
	for r := l.peek(); !unicode.IsSpace(r) && r != eof; r = l.peek() {
		l.consume()
	}
}

func (l *lexer) consumeIdent() {
	for isIdent(l.peek()) {
		l.consume()
	}
}

func (l *lexer) consumeComment() {
	for r := l.peek(); r != '\n' && r != eof; r = l.peek() {
		l.consume()

		if r == '\\' {
			if l.peek() == eof {
				l.error("unexpected EOF")
				break
			}

			l.consume()
		}
	}
}

func (l *lexer) consumeString() {
	for r := l.peek(); r != '"' && r != eof; r = l.peek() {
		l.consume()

		if r == '\\' {
			if l.peek() == eof {
				l.error("unexpected EOF")
				break
			}

			l.consume()
		}
	}

	if l.peek() == eof {
		l.error("unexpected EOF")
	}
}

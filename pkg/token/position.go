// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
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

package token

import "fmt"

// Position represents a specific line and column in a source string.
type Position struct {
	Line int
	Col  int
}

// String returns a string representation of p, in the format line:column.
func (p *Position) String() string {
	return fmt.Sprintf("%v:%v", p.Line, p.Col)
}

// NextLine emulates going to the next line from position p in a string by
// increasing line by 1 and setting column to 1, or the first column.
func (p *Position) NextLine() {
	p.Line++
	p.Col = 1
}

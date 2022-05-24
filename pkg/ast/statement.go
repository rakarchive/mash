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

package ast

type Statement interface {
	Node
	Statement()
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) Node()      {}
func (b *BlockStatement) Statement() {}

type IfStatement struct {
	Condition Expression
	BlockStmt *BlockStatement
	ElseBlock Statement
}

func (i *IfStatement) Node()      {}
func (i *IfStatement) Statement() {}

type ForStatement struct {
	Condition Expression
	BlockStmt *BlockStatement
}

func (f *ForStatement) Node()      {}
func (f *ForStatement) Statement() {}

type LetStatement struct {
	Expression Expression
}

func (l *LetStatement) Node()      {}
func (l *LetStatement) Statement() {}

type CmdStatement struct {
	Command Command
}

func (c *CmdStatement) Node()      {}
func (c *CmdStatement) Statement() {}

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

import "laptudirm.com/x/mash/pkg/token"

type Expression interface {
	Node
	Expression()
}

type Assignable interface {
	Expression
	Assignable()
}

type AssignExpression struct {
	Left     Assignable
	Operator token.Token
	Right    Expression
}

func (a *AssignExpression) Node()       {}
func (a *AssignExpression) Expression() {}

type LogicalExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (l *LogicalExpression) Node()       {}
func (l *LogicalExpression) Expression() {}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b *BinaryExpression) Node()       {}
func (b *BinaryExpression) Expression() {}

type UnaryExpression struct {
	Operator token.Token
	Right    Expression
}

func (u *UnaryExpression) Node()       {}
func (u *UnaryExpression) Expression() {}

type GroupExpression struct {
	Right Expression
}

func (g *GroupExpression) Node()       {}
func (g *GroupExpression) Expression() {}

type CallExpression struct {
	Callee      Expression
	Parenthesis token.Token
	Arguments   []Expression
}

func (c *CallExpression) Node()       {}
func (c *CallExpression) Expression() {}

type GetExpression struct {
	Name Expression
	Expr Expression
}

func (g *GetExpression) Node()       {}
func (g *GetExpression) Expression() {}
func (g *GetExpression) Assignable() {}

type SelectorExpression struct {
	Name  Expression
	Index token.Token
}

func (s *SelectorExpression) Node()       {}
func (s *SelectorExpression) Expression() {}
func (s *SelectorExpression) Assignable() {}

type VariableExpression struct {
	Name token.Token
}

func (v *VariableExpression) Node()       {}
func (v *VariableExpression) Expression() {}
func (v *VariableExpression) Assignable() {}

type LiteralExpression struct {
	Literal Literal
}

func (l *LiteralExpression) Node()       {}
func (l *LiteralExpression) Expression() {}

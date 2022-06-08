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

// NumberLiteral node represents a number constant.
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) Node()       {}
func (n *NumberLiteral) Expression() {}

// StringLiteral node represents a string constant.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (n *StringLiteral) Node()             {}
func (n *StringLiteral) Expression()       {}
func (n *StringLiteral) CommandComponent() {}

// FunctionLiteral node represents a function expression.
type FunctionLiteral struct {
	Token token.Token
	Block *BlockStatement
}

func (n *FunctionLiteral) Node()       {}
func (n *FunctionLiteral) Expression() {}

// ArrayLiteral node represents a array expression.
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) Node()       {}
func (a *ArrayLiteral) Expression() {}

// ObjectLiteral node represents an object expression.
type ObjectLiteral struct {
	Token    token.Token
	Elements map[Expression]Expression
}

func (o *ObjectLiteral) Node()       {}
func (o *ObjectLiteral) Expression() {}

// TemplateLiteral node represents a template string expression.
type TemplateLiteral struct {
	Expressions []Expression
	Components  []token.Token
}

func (t *TemplateLiteral) Node()             {}
func (t *TemplateLiteral) Expression()       {}
func (t *TemplateLiteral) CommandComponent() {}

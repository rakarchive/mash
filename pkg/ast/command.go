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

type Command interface {
	Node
	Command()
}

type CommandComponent interface {
	Node
	CommandComponent()
}

type LogicalCommand struct {
	Left     Command
	Operator token.Token
	Right    Command
}

func (l *LogicalCommand) Node()    {}
func (l *LogicalCommand) Command() {}

type UnaryCommand struct {
	Operator token.Token
	Right    Command
}

func (u *UnaryCommand) Node()    {}
func (u *UnaryCommand) Command() {}

type BinaryCommand struct {
	Left     Command
	Operator token.Token
	Right    Command
}

func (b *BinaryCommand) Node()    {}
func (b *BinaryCommand) Command() {}

type LiteralCommand struct {
	Components []CommandComponent
}

func (l *LiteralCommand) Node()    {}
func (l *LiteralCommand) Command() {}

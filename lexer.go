/* 
 * Copyright (C) 2015 Simon Schmidt
 *
 * This Source Code is subject to the terms of the Mozilla Public License, v. 2.0.
 * If a copy of the MPL was not distributed with this file, You can obtain one at:
 *
 *                 http://mozilla.org/MPL/2.0/
 */

package lexer

import "regexp"

type ParserError int
func (i ParserError) Error() string {
	switch i {
	case PE_NO_TOKEN_MATCHING: return "PE_NO_TOKEN_MATCHING"
	default: return "generic error (unknown)"
	}
}

const (
	PE_UNKNOWN ParserError = iota
	PE_NO_TOKEN_MATCHING
)

type Result struct{
	Number int
	Text   string
	Count  int
	Start  int
}

type Token struct{
	Number int
	Regexp string
}
type rule struct{
	number int
	expr   *regexp.Regexp
}

type Ruleset []rule
func Compile(t []Token) (rs Ruleset,err error) {
	rs = make(Ruleset,len(t))
	for i,tt := range t {
		rs[i].number = tt.Number
		rs[i].expr,err = regexp.Compile("^"+tt.Regexp)
		if err!=nil { rs = nil; return }
	}
	return
}
func (r Ruleset) Parse(txt string) (res []Result,err error) {
	var fi []int
	var num int
	lastnum := 0
	pos := 0
	for len(txt)>0 {
		fi = nil
		for _,ru := range r {
			num = ru.number
			if fi = ru.expr.FindStringIndex(txt); fi!=nil { break }
		}
		if fi==nil { err=PE_NO_TOKEN_MATCHING; return }
		cur := txt[:fi[1]]
		txt  = txt[fi[1]:]
		pos += fi[1]
		if num<0 {
			num =- num
			if lastnum==num {
				res[len(res)-1].Count++
				continue
			}
			lastnum = num
		} else if num==0 {
			continue
		} else {
			lastnum = 0
		}
		res = append(res,Result{num,cur,1,pos-fi[1]})
	}
	return
}

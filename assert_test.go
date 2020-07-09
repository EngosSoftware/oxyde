/*
 * MIT License
 *
 * Copyright (c) 2017-2020 Dariusz Depta Engos Software
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE. 
 */
 
package oxyde

import (
    "testing"
)

func TestEqualStrings(t *testing.T) {
    if !equalString("string", "string") {
        t.Error("strings are equal but test shows they are not")
    }
    if equalString("string", "string1") {
        t.Error("strings are not equal but test shows they are")
    }
    if equalString("Alfa", "Alfb") {
        t.Error("strings are not equal but test shows they are")
    }
    if equalString("string", "") {
        t.Error("strings are not equal but test shows they are")
    }
    if equalString("", "string") {
        t.Error("strings are not equal but test shows they are")
    }
}

func TestEqualStringsNullable(t *testing.T) {
    s1 := "string"
    s2 := "string"
    if !equalStringNullable(&s1, &s2) {
        t.Error("strings are equal but test shows they are not")
    }
    if !equalStringNullable(nil, nil) {
        t.Error("strings are equal but test shows they are not")
    }
    if equalStringNullable(&s1, nil) {
        t.Error("strings are not equal but test shows they are")
    }
    if equalStringNullable(nil, &s2) {
        t.Error("strings are not equal but test shows they are")
    }
    s1 = "Alfa"
    s2 = "Alfb"
    if equalStringNullable(&s1, &s2) {
        t.Error("strings are not equal but test shows they are")
    }
}

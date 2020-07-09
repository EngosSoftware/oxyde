/*
 * MIT License
 *
 * Copyright (c) 2017-2020 Engos Software
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

// JsonApiError is a JSON API implementation of an error.
type JsonApiError struct {
    Status *string `json:"status" api:"The HTTP status code applicable to reported problem."`
    Code   *string `json:"code"   api:"An application-specific error code."`
    Title  *string `json:"title"  api:"A short, human-readable summary of the problem that never changed from occurrence to occurrence of the problem."`
    Detail *string `json:"detail" api:"A human-readable explanation specific to the occurrence of the problem."`
}

// JsonApiErrors is an array of JSON API errors.
type JsonApiErrors = []JsonApiError
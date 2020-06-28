# Goto
> Goto is an interpreted programming language written in go.

[![License: MIT](http://img.shields.io/github/license/pandeykartikey/goto.svg)](https://github.com/pandeykartikey/goto/blob/master/LICENSE.md) [![Go Report Card](https://goreportcard.com/badge/github.com/pandeykartikey/goto)](https://goreportcard.com/report/github.com/pandeykartikey/goto) [![Release](http://img.shields.io/github/v/release/pandeykartikey/goto.svg)](https://github.com/pandeykartikey/goto/releases/latest)

## 1. Overview
Goto is a dynamically typed programming language written to support all the scripting requirements. It currently supports the following:
- Data Types: `integer`, `boolean`, `string`
- Data Structures: `list`
- Arithimetic Operations: `+`, `-`, `*`, `/`, `%`, `**`
- Comparisons: `==`, `!=`, `<`, `<=`, `>`, `>=` 
- Logical Operators:  `!`, `&&`, `||` 
- If-Else-If Statements
- For loops
- Control Flow Statements `continue`, `break`, `return`
- Multiple Assigments
- Operator Precedence Parsing
- Grouped Expressions
- Functions
- Scopes
- Error Handling
- Built in Functions: `append`, `print`, `len`

## 2. Table of Content
  - [1. Overview](#1-overview)
  - [2. Table of Content](#2-table-of-content)
  - [3. Installation](#3-installation)
    - [Source Installation](#source-installation)
    - [Binary Releases](#binary-releases)
  - [4. Usage](#4-usage)
  - [5. Syntax](#5-syntax)
    - [5.1 Definitions](#51-definitions)
      - [5.1.1 Multiple Assignments](#511-multiple-assignments)
      - [5.1.2 Scoping](#512-scoping)
    - [5.2 Arithmetic operations](#52-arithmetic-operations)
    - [5.3 Lists](#53-lists)
      - [5.3.1 Indexing](#531-indexing)
    - [5.4 Builtin functions](#54-builtin-functions)
    - [5.5 Functions](#55-functions)
      - [5.5.1 Local Functions](#551-local-functions)
    - [5.6 If-else statements](#56-if-else-statements)
    - [5.7 For-loop statements](#57-for-loop-statements)
    - [5.8 Control flow statements](#58-control-flow-statements)
  - [6. Contributing](#6-contributing)
  - [7. Acknowledgments](#7-acknowledgments)
  - [8. License](#8-license)
  - [9. Contact](#9-contact)


## 3. Installation

### Source Installation
To install `goto`, run the following command:   

    git clone https://github.com/pandeykartikey/goto
    cd goto
    go install


### Binary Releases

Alternatively, you could install a binary-release, from the [release page](https://github.com/pandeykartikey/goto/releases).

## 4. Usage
To execute a goto-script, pass the name to the interpreter:

    $ goto sample.to

To drop into goto-repl, type `goto`. 

## 5. Syntax

### 5.1 Definitions
Variables are defined using the `var` keyword, with each line ending with `;`.

    var a = 3;

A variable must be declared before its usage. It is not necessary to provide a value while declaring a variable.

    var a;

#### 5.1.1 Multiple Assignments
Goto supports multiple assignments and declarations.

    var a,b,c = 1,2,3;

The datatypes of all the variables need not be the same.

    a,b = 1,true;

#### 5.1.2 Scoping
Goto supports hiding of global variable in block constructs

    var a = 4;
    if true { var a = 5; print(a);} //prints 5
    print(a); //prints 4


### 5.2 Arithmetic operations
Goto supports all the basic arithmetic operations along with `**` operator for power. (Inspired from Python)

    square = b**2;
    remainder = b%2;

### 5.3 Lists
List is a data structure that organizes items by linear sequence. It can hold multiple types.

    var a = [1, true, "array"];

#### 5.3.1 Indexing
Lists are 0-indexed. Elements can be accessed using []. Similar indexing exists for strings.

    var a = [1, true, "array"];
    a[1] // returns true
    a[2][3] // returns "a"

### 5.4 Builtin functions
Goto currently supports 3 built-in functions:
1. `len`: Returns the length of string or a list.

    len("goto") //returns 4

2. `append`: appends a token at the end of an array

    var a = [1, 2];
    append(a, 3) // a becomes [1, 2, 3]

3. `print`: prints the content of parameters to STDOUT. It can take multiple arguments.

    print(1, 2, "goto")
    Output: 1
            2
            goto

### 5.5 Functions
Goto defines function using `func` followed by an identifier and a parameter list.

    func identity(x) {
	    return x;
    }

#### 5.5.1 Local Functions
You can define local functions inside a block statement with limited scope.

    func addTwo(x) {
      func addOne(x) { // addOne() is limited to addTwo()'s scope
        return x + 1;
      }
      x = addOne(x);
      return addOne(x);
    }



### 5.6 If-else statements
Goto supports if-else-if statements.
    
    var a,b,c = 1,2,3;
    if a > b { 
      c = a + b; 
    } else if a < c {
      c = 20; 
    } else {
      c = 30; 
    }
    print(c); //returns 20

### 5.7 For-loop statements
Goto has supports for-loop statements.

    for var i = 0; i < 10; i = i + 1 {
       i = i + 2;
       print(i);
    }

All the three initialization, condition, and update are optional in for loop.

    for ;; {
       i = i + 2;
       print(i);
    }

### 5.8 Control flow statements
There are three control flow statements in goto:

1. `continue`: It skips all the following statements in a for loop and moves on to the next iteration.

2. `break`: It is used to break a for loop.

3. `return`: It is used to terminate a function. It may also be used to return values from functions.  


## 6. Contributing
If you spot anything that seems wrong, please do [report an issue](https://github.com/pandeykartikey/goto/issues).

If you feel something is missing in goto, consider creating an [issue](https://github.com/pandeykartikey/goto/issues) or submitting a [pull request](https://github.com/pandeykartikey/goto/pulls).

## 7. Acknowledgments
This programming language takes inspiration from the Book [Write an Interpreter in Go](https://interpreterbook.com).

## 8. License
Goto is licensed under [MIT License](https://github.com/pandeykartikey/goto/blob/master/LICENSE.md).

## 9. Contact
If you have any queries regarding the project or just want to say hello, feel free to drop a mail at
 pandeykartikey99@gmail.com.
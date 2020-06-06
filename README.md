## Compiler for custom minimalistic grammar
I made up the following grammar and programmed a compiler named Compilang for it. This is a university project.

### Prerequisites
1. Install Golang by following the tutorial [here](https://golang.org/dl/)
2. Install git by following the tutorial [here](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Download and install Compilang compiler
1. `git clone git@github.com:nitrajka/compiler.git`
2. Go to `compiler/cmd/compilang`
3. Run `go install` 

### Write and compile a Copmilang program
1. Write a program in a working directory `program.compilang`
2. The command `compilang program.compilang` compiles your program. If not specified flag `-o output_file.go`, the program will be compiled to `a.go` by default. You do not need to create the `output_file.go` manually, the compiler creates it for you.
3. Run the compiled program with command `go run a.go` or `go run output_file.go` if you specified and output file in compilation step.


### Language syntax and examples
Examples of simple programs with their expected outputs can be found in [tests](https://github.com/nitrajka/compiler/tree/master/tests) folder.
Each example is written in its own file `example.compilang`. You can find its expected output in `example.output` file.
Here is an example of fibonacci algorithm written in Compilang.:
```
globals

endglobals

func fibonacci(int [n]): int {
	int [k, p];
    
    if n <= 0 {; print("invalid input") }
    else {;
        if n == 1 {; return 0}
        if n == 2 {; return 1}
        k = n-1
        p = n-2
        return call fibonacci(k) + call fibonacci(p)
    }
    return 0
}

main { 
    int [n];
    
    while n < 11 {;
        n = n+1
        print(call fibonacci(n))
    }
}
endmain
```


### Grammar of Compilang
The actual version of grammar can be found in [grammar.md](https://github.com/nitrajka/compiler/blob/master/pkg/grammar.md).
The structure of a Compilang program is as follows:

    1. Global variables
    2. Function definitions
    3. Main function

It is not allowed to write code anywhere else than inside functions.
#### Global variables and variable list
Global variables can be defined as follows:

```
globals
    int [i, j k]
    string [s, r] bool [b,a]
endglobals

...

globals
    int[i,j]bool[a,b]
endglobals
```
Notice, that there is no need for spaces between the type and list of variables,
as well as between list of variables and the following type. Variables of type `void` are not allowed to be created.

#### Functions

A function can be defined as follows:

    1. Keyword `func`
    2. Function's identifier
    3. Variable parameters declared as in `globals` section and enclosed in parenthesis.
    4. Colon
    5. Type 
    6. Body
```
func name(string [a]): string {; return ""}
```
Return clause is optional for functions of type `void`. However, if you want to add a return clause for such function, make sure you return `void` - `return void`.

Calling a function without assigning it to a variable is a valid Compilang call. Such call makes sense if the function modifies global variables. 

It is invalid to assign call of a void type function to a variable. However, since it is not possible to create a variable of `void` type, the mentioned assignment is not possible to make.

It is invalid to have 2 functions of the same name. 
It is not possible to call main function.

#### Body
Body consists of 2 parts divided by a semicolon:

    1. variable declarations
    2. statements
    
Since Compilang is compiled to Golang, I used the feature of assigning a base value for a variable of a type.
For instance, in the example below, 

    * variables `i` and `j` are assigned 0 during their declaration
    * variable `a` is assigned an empty string `""`
    * variable `b` is assigned `false`
```
{
    int [i, j] bool [b]
    string [a];

    if i == j {; print("i == j")}
    else {; print("i != j")}

    a = "hello world"
    print(a)
    
    j = 10
    while i < j {;
        print(i)
        print(var i)
        i = i+1
    }
}
```
The following four types of statement are supported in a body:

    * assignment
    * while loop
    * if statement
    * print statement

#### Expressions
Bool expression is expression which contains some of the bool operators - `&&`, `||`, `<`, `>`, `==`, `!=`, `<=`, `>=`.
Expression contains only arithmetic operators - `+`, `-`, `*`, `/`, `%`.

It is not possible to assign a bool expression, return a bool expression or print a bool expression.
For usage of expressions and bool expressions see examples [here](https://github.com/nitrajka/compiler/blob/master/tests/expressions.compilang) and [here](https://github.com/nitrajka/compiler/blob/master/tests/bool_expressions.compilang).

#### Assignment
Before assigning a value to a variable, this variable must have been declared in current or in one of the outer scopes.
According to the [grammar](https://github.com/nitrajka/compiler/blob/master/pkg/grammar.md) almost anything can be assigned to a variable but bool expression.
```
{ bool [b];
    b = true && false
}
```

#### Scopes
Scopes are nested. The outermost scope is global scope, which contains variables declared (actually defined with their base value as well) in `gloabls` section.
Exactly two scopes are defined before the first statement in main function:
 
    * globals cope
    * scope created from variables defined in its body.
Exactly two scopes are defined before the first statement in a function:

    * global scope
    * scope created from variables defined in its body and parameter variables


Each time a variable is used, the semantics subprogram checks whether the variable exists. 
If there are two variables with the same name in two different scopes, the variable defined in the closest scope on the way from current scope to the outermost scope is used.
Example:
```
globals
    int [i]
endglobals

func f(): int {int [i];
    i = 3
    print(i)
    return i
}
main {;
    i = 5
    print(call f())
    print(i)
} endmain
```
Outputs: 
```
3
3
5
```

It is not possible to have variables with the same name in the same scope.
Each time a new body is defined, new inner scope is created. Variables declared at the beginning of a function's body share scope with function's parameters.

#### Printing and returning
At most 1 variable is allowed in print statement. Empty print `print()` prints a newline.
Exactly 1 variable is allowed in return clause. It is not possible to return a bool expression.
For usage of `print()` and `return` see examples in `tests` folder - [here](https://github.com/nitrajka/compiler/blob/master/tests/printing.compilang) and [here](https://github.com/nitrajka/compiler/blob/master/tests/returning.compilang).


#### General notes
If you are not sure about where to put a whitespace, take a look into the [grammar](https://github.com/nitrajka/compiler/blob/master/pkg/grammar.peg).

Although the grammar supports arrays and maps, semantics and generator does not. Please, do not use these types yet.
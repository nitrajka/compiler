## Compiler for custom minimalistic grammar
I made up the following grammar and programmed a compiler for it. This is a university project.

### How to compile

1. Downlaod and install the compiler  by running `go get github.com/nitrajka/compiler`.
1. Compile program `compilang myprogram.col`.

### Grammar
```
ROOT -> 'globals' PARAMS_VARS 'endglobals' FUNCTIONS 'main' BODY 'endmain'

FUNCTIONS -> FUNCTION FUNCTIONS | ε

FUNCTION -> 'func' ID '(' PARAMS_VARS ')' ':' TYPE BODY

PARAMS_VARS -> TYPE '[' VAR_LIST ']' PARAMS_VARS | ε

----------------------------------------------
VAR_LIST -> ID | VAR_LIST ',' VAR_LIST
----------------------------------------------
VAR_LIST -> ID VAR_LIST'
VAR_LIST' -> ',' VAR_LIST VAR_LIST' | ε
----------------------------------------------

TYPE -> 'string' | 'int' | '[' [0-9]+ ']' TYPE | 'map[' TYPE ']' TYPE | 'bool' | 'void'

BODY -> '{' PARAMS_VARS ';' STATEMENTS RETURN_CLAUSE '}'

STATEMENTS -> STATEMENT STATEMENTS | ε

STATEMENT ->
    ASSIGNMENT
    | 'if' BOOL_EXPRESSION BODY ELSECLAUSE
    | 'while' BOOL_EXPRESSION BODY
    | FUNC_CALL

ELSECLAUSE -> 'else' BODY | ε

RETURN_CLAUSE -> 'return' VALUE

VALUE -> 'var' ID | STRING | INTEGER | BOOLEAN | BOOL_EXPRESSION | ARRAY | EXPRESSION | MAP | 'void'

----------------------------------------------
BOOL_EXPRESSION -> ID
            | BOOL_EXPRESSION '<' BOOL_EXPRESSION
            | BOOL_EXPRESSION '>' BOOL_EXPRESSION
            | BOOL_EXPRESSION '==' BOOL_EXPRESSION
            | BOOL_EXPRESSION '!=' BOOL_EXPRESSION
            | BOOL_EXPRESSION '<=' BOOL_EXPRESSION
            | BOOL_EXPRESSION '>=' BOOL_EXPRESSION
----------------------------------------------
BOOL_EXPRESSION -> ID BOOL_EXPRESSION'
BOOL_EXPRESSION' -> '<' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '>' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '==' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '!=' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '<=' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '>=' BOOL_EXPRESSION BOOL_EXPRESSION'
            | ε
----------------------------------------------

ASSIGNMENT -> 'var' ID TYPE '=' VALUE 
            | 'let' ID '=' VALUE
            | ID '=' VALUE

FUNC_CALL -> 'call' ID '(' VAR_LIST ')'

----------------------------------------------
EXPRESSION -> ID 
    | EXPRESSION '+' EXPRESSION
    | EXPRESSION '-' EXPRESSION
    | EXPRESSION '*' EXPRESSION
    | EXPRESSION '/' EXPRESSION
----------------------------------------------
EXPRESSION -> ID EXPRESSION'
EXPRESSION' -> '+' EXPRESSION EXPRESSION'
    | '-' EXPRESSION EXPRESSION'
    | '*' EXPRESSION EXPRESSION'
    | '/' EXPRESSION EXPRESSION'
    | ε
----------------------------------------------

ID -> [a-z]+
INTEGER -> INT | '-' INT
BOOLEAN -> 'true' | 'false'
STRING -> '"' [a-z]+ '"'
INT -> [0-9]+

ARRAY -> '[]' TYPE '[' VALUE ']'

----------------------------------------------------OPTIONAL: MAP-----------------------------------------------------------
MAP -> '(' KEY_VALUE_PAIRS ')'

----------------------------------------------
KEY_VALUE_PAIRS -> STRING ':' VALUE | KEY_VALUE_PAIRS ',' KEY_VALUE_PAIRS
----------------------------------------------
KEY_VALUE_PAIRS -> STRING ':' VALUE KEY_VALUE_PAIRS'
KEY_VALUE_PAIRS' -> ',' KEY_VALUE_PAIRS KEY_VALUE_PAIRS' |  ε
----------------------------------------------


SPECIAL_CHARS -> \r | \n | \t
```

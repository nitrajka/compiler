## Compiler for custom minimalistic grammar
The grammar is as follows:
```
PRORGAM -> STATEMENT+ '_|'
STATEMENT ->
	'if' PAREN_EXPRESSION STATEMENT
	| 'if' PAREN_EXPRESSION STATEMENT 'else' STATEMENT
	| 'while' PAREN_EXPRESSION STATEMENT
	| '{' STATEMENT '}'
	| EXPRESSION ';'
	| ';'

PAREN_EXPRESSION -> '(' EXPRESSION ')'
EXPRESSION -> TEST | id '=' EXPRESSION | 'var' ID TYPE '='

TEST ->
	SUM
	| SUM '<' SUM
	| SUM '>' SUM
	| SUM '==' SUM

SUM ->
	TERM
	| SUM '+' TERM
	| SUM '-' TERM
	| SUM '-=' TERM
	| SUM '+=' TERM

TERM -> ID | INTEGER | PAREN_EXPRESSION

ID -> STRING
INTEGER -> INT

STRING -> [a-z]+
INT -> [0-9]+

SPECIAL_CHARS -> \r | \n | \t
```
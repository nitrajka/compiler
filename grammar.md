```
PRORGAM -> STATEMENT+'_|'
FI1(PROGRAM) = FI1(STATEMENT)

STATEMENT ->
'if' PAREN_EXPRESSION STATEMENT
| 'if' PAREN_EXPRESSION STATEMENT 'else' STATEMENT
| 'while' PAREN_EXPRESSION STATEMENT
| '{' STATEMENT '}'
| EXPRESSION ';'
| ';'
FI(STATEMENT) = {'if', 'while', '{', ';'} + FI1(EXPRESSION) 
              = {'if', 'while', '{', ';'} + {'var'} + FI1(TEST) + FI1(ID) 
              = {'if', 'while', '{', ';'} + {'var'} + FI1(SUM) + FI1(STRING) 
              = {'if', 'while', '{', ';'} + {'var'} + FI1(TERM) + ([a-z]+) 
              = {'if', 'while', '{', ';'} + {'var'} + FI1(ID) + FI1(INTEGER) + FI1(PAREN_EXPRESSION) + ([a-z]+) 
              = {'if', 'while', '{', ';'} + {'var'} + FI1(STRING) + FI1(INT) + {'('} + ([a-z]+) 
              = {'if', 'while', '{', ';'} + {'var'} + {[a-z]+} + {[0-9]+} + {'('} + {[a-z]+} = ?conflict?

PAREN_EXPRESSION -> '(' EXPRESSION ')'
FI1(PAREN_EXPRESSION) = {'('}

EXPRESSION -> TEST | ID '=' EXPRESSION | 'var' ID TYPE '='
FI1(EXPRESSION) = FI1(SUM) + FI1(ID) + {'var'} 
                = FI1(TERM) + FI1(STRING) + {'var'} 
                = FI1(ID) + FI1(INTEGER) + FI1(PAREN_EXPRESSION) + ([a-z]+) + {'var'} 
                = FI1(STRING) + FI1(INT) + {'('} + ([a-z]+) + {'var'} 
                =([a-z]+) + ([0-9]+) + {'('} + ([a-z]+) + {'var'} = ?conflict?

TEST -> SUM | SUM '<' SUM | SUM '>' SUM | SUM '==' SUM
FI1(TEST) = FI1(SUM) = FI1(TERM) 
          = FI1(ID) + FI1(INTEGER) + FI1(PAREN_EXPRESSION)
          = FI1(STRING) + FI1(INT) + {'('}
          =  ([a-z]+) + ([0-9]+) + {'('} 
          = {([a-z]+), ([0-9]+), '('}


--------------------------------------------
SUM -> TERM | SUM '+' TERM | SUM '-' TERM | SUM '-=' TERM | SUM '+=' TERM
Left recursion removal:
---------------------------------------------
SUM -> TERM SUM' 
SUM' -> '+' TERM SUM' | '-' TERM SUM' | '-=' TERM SUM' | '+=' TERM SUM' | ε
FI1(SUM') = {'+', '-', '-=', '+='} + FO1(SUM')
          = {'+', '-', '-=', '+='} + FO1(SUM)
          = {'+', '-', '-=', '+='} + {'<', '>', '=='} + FO1(TEST)
          = {'+', '-', '-=', '+='} + {'<', '>', '=='} + FO1(EXPRESSION)
          = {'+', '-', '-=', '+='} + {'<', '>', '=='} + {')', ';'}
          = {'+', '-', '-=', '+=', '<', '>', '==', ')', ';'}

TERM -> ID | INTEGER | PAREN_EXPRESSION
FI1(TERM) = FI1(ID) + FI1(INTEGER) + FI1(PAREN_EXPRESSION)
          = FI1(STRING) + FI1(INT) + {'('}
          =  ([a-z]+) + ([0-9]+) + {'('} 
          = {([a-z]+), ([0-9]+), '('}

ID -> STRING
FI1(ID) = {([a-z]+)}

INTEGER -> INT
FI1(INTEGER) = {([0-9]+)}

STRING -> [a-z]+
INT -> [0-9]+

SPECIAL_CHARS -> \r | \n | \t
```
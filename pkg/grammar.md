```
ROOT -> 'globals' PARAMS_VARS 'endglobals' FUNCTIONS 'main' BODY 'endmain'
FI1(ROOT) = {'globals'}

FUNCTIONS -> FUNCTION FUNCTIONS | ε
FI1(FUNCTIONS) = FI1(FUNCTION) + FO1(FUNCTIONS) 
               = {'func'} + {'main'}

FUNCTION -> 'func' ID '(' PARAMS_VARS ')' ':' TYPE BODY
FI1(FUNCTION) = {'func'}

PARAMS_VARS -> TYPE '[' VAR_LIST ']' PARAMS_VARS | ε
FI1(PARAMS_VARS) = FI1(TYPE) + FO1(PARAMS_VARS)
                 = {'string', 'int', '[', 'map[', 'bool', 'void'} + {')', 'endglobals', ';'}

----------------------------------------------
VAR_LIST -> ID | VAR_LIST ',' VAR_LIST
----------------------------------------------
VAR_LIST -> ID VAR_LIST'
VAR_LIST' -> ',' VAR_LIST VAR_LIST' | ε
----------------------------------------------
FI1(VAR_LIST) = FI1(ID) = FI1([a-z]+)
FI1(VAR_LIST') = {','} + FO1(VAR_LIST') 
               = {','} + FO1(VAR_LIST)
               = {','} + {']', ')'}

TYPE -> 'string' | 'int' | '[' [0-9]+ ']' TYPE | 'map[' TYPE ']' TYPE | 'bool' | 'void'
FI1(TYPE) = {'string', 'int', '[', 'map[', 'bool', 'void'}

BODY -> '{' PARAMS_VARS ';' STATEMENTS RETURN_CLAUSE '}'
FI1(BODY) = {'{'}

STATEMENTS -> STATEMENT STATEMENTS | ε
FI1(STATEMENTS) = FI1(STATEMENT) + FO1(STATEMENTS)
                = {'var', 'let', [a-z]+, 'if', 'while', 'call'} + FI1(RETURN_CLAUSE)
                = {'var', 'let', [a-z]+, 'if', 'while', 'call'} + {'return'}

STATEMENT ->
    ASSIGNMENT
    | 'if' BOOL_EXPRESSION BODY ELSECLAUSE
    | 'while' BOOL_EXPRESSION BODY
    | FUNC_CALL
    | 'print(' VALUE? ')'
FI1(STATEMENT) = FI1(ASSIGNEMENT) + {'if', 'while'} + FI1(FUNC_CALL) + FI1('print(')
               = FI1(ASSIGNEMENT) + {'if', 'while'} + {'call'} + {'print('}
               = FI1(ASSIGNEMENT) + {'if', 'while'} + {'call'} + {'print('}
               = {'var', 'let', [a-z]+} + {'if', 'while'} + {'call'} + {'print('}

ELSECLAUSE -> 'else' BODY | ε
FI1(ELSECLAUSE) = {'else'} + FO1(ELSECLAUSE)
                = {'else'} + FO1(STATEMENT)
                = {'else'} + FI1(STATEMENTS)
                = {'else'} + {'var', 'let', [a-z]+, 'if', 'while', 'call', 'return'}

RETURN_CLAUSE -> 'return' VALUE
FI1(RETURN_CLAUSE) = {'return'}

VALUE -> 'var' ID | 'arr' ID INDEXED? | BOOLEAN | ARRAY | EXPRESSION | MAP | 'void'
FI1(VALUE) = {'var', '"', [0-9]+, '-', 'true', 'false', '[]', 'voidV'} + FI1(BOOL_EXPRESSION) + FI1(EXPRESSION) + FI1(MAP)
           = {'var', '"', [0-9]+, '-', 'true', 'false', '[]', 'voidV'} + {'<', '>', '==', '!=', '<=', '>=', '{'} + {[a-z]+} + {'(', ','}

EXPR_VALUE <- 'arr' ID INDEXED / FUNC_CALL / TEXT / INTEGER / ARRAY / ID

BOOL_EXPR_VALUE <- BOOLEAN / 'arr' ID INDEXED / ID / INTEGER / TEXT

INDEXED <- '[' INDEXABLE ']'

INDEXABLE <- 'var' ID / INT

----------------------------------------------
BOOL_EXPRESSION -> BOOL_EXPR_VALUE
            | BOOL_EXPRESSION '<' BOOL_EXPRESSION
            | BOOL_EXPRESSION '>' BOOL_EXPRESSION
            | BOOL_EXPRESSION '==' BOOL_EXPRESSION
            | BOOL_EXPRESSION '!=' BOOL_EXPRESSION
            | BOOL_EXPRESSION '<=' BOOL_EXPRESSION
            | BOOL_EXPRESSION '>=' BOOL_EXPRESSION
----------------------------------------------
BOOL_EXPRESSION -> BOOL_EXPR_VALUE BOOL_EXPRESSION'
BOOL_EXPRESSION' -> '<' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '>' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '==' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '!=' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '<=' BOOL_EXPRESSION BOOL_EXPRESSION'
            | '>=' BOOL_EXPRESSION BOOL_EXPRESSION'
            | ε
----------------------------------------------
FI1(BOOL_EXPRESSION) = {[a-z]+, 'true', 'false', 'arr', '-', [0-9], '"'}
FI1(BOOL_EXPRESSION') = {'<', '>', '==', '!=', '<=', '>='} + FO1(BOOL_EXPRESSION')
    = {'<', '>', '==', '!=', '<=', '>='} + FO1(BOOL_EXPRESSION)
    = {'<', '>', '==', '!=', '<=', '>='} + FI1(BODY) + FO1(VALUE)
    = {'<', '>', '==', '!=', '<=', '>='} + {'{'} + FO1(RETURN_CLAUSE) + FO1(ASSIGNMENT) + {']'}   //OPTIONAL: + FI1(KEY_VALUE_PAIRS')
    = {'<', '>', '==', '!=', '<=', '>='} + {'{'} + {'}'} + FO1(STATEMENT) + {']'}         //OPTIONAL: + {',', ')'}
    = {'<', '>', '==', '!=', '<=', '>='} + {'{'} + {'}'} + FI1(STATEMENTS) + {']'}      //OPTIONAL: + {',', ')'}
    = {'<', '>', '==', '!=', '<=', '>='} + {'{'} + {'}'} + {'var', 'let', [a-z]+, 'if', 'while', 'call', 'return'} + {']'} //OPTIONAL: + {',', ')'}

ASSIGNMENT -> ASSIGNABLE '=' VALUE
FI1(ASSIGNMENT) = {'arr'} + {[a-z]+}

ASSIGNABLE <- 'arr' ID '[' INDEXABLE ']' / ID

FUNC_CALL -> 'call' ID '(' VAR_LIST* ')'
FI1(FUNC_CALL) = {'call'}

----------------------------------------------
EXPRESSION -> EXPR_VALUE
    | EXPRESSION '+' EXPRESSION
    | EXPRESSION '-' EXPRESSION
    | EXPRESSION '*' EXPRESSION
    | EXPRESSION '/' EXPRESSION
----------------------------------------------
EXPRESSION -> EXPR_VALUE EXPRESSION'
EXPRESSION' -> '+' EXPRESSION EXPRESSION'
    | '-' EXPRESSION EXPRESSION'
    | '*' EXPRESSION EXPRESSION'
    | '/' EXPRESSION EXPRESSION'
    | ε
----------------------------------------------
FI1(EXPRESSION) = {[a-z]+, 'arr', 'call', '"', '-', [0-9]+, '[]'}
FI1(EXPRESSION') = {'+', '-', '*', '/'} + FO1(EXPPRESSION')
                 = {'+', '-', '*', '/'} + FO1(EXPRESSION)
                 = {'+', '-', '*', '/'} + FO1(VALUE)    // + FI1(KEY_VALUE_PAIRS')
                 = {'+', '-', '*', '/'} + FO1(RETURN_CLAUSE) + FO1(ASSIGNMENT) + {']'}   //OPTIONAL: + {',', ')'}
                 = {'+', '-', '*', '/'} + {'}'} + FO1(STATEMENT) + {']'}  //OPTIONAL: + {',', ')'}
                 = {'+', '-', '*', '/'} + {'}'} + FI1(STATEMENTS) + {']'}  //OPTIONAL: + {',', ')'}
                 = {'+', '-', '*', '/'} + {'}'} + {'var', 'let', [a-z]+, 'if', 'while', 'call', 'return'} + {']'} //OPTIONAL: + {',', ')'}

ID <- FIRST_CHAR ID1*
FIRST_CHAR <- [a-z] / '_' / [A-Z]
ID1 <- [a-z] / '_' / [A-Z] / [0-9]

INTEGER -> INT | '-' INT
BOOLEAN -> 'true' | 'false'
TEXT <- '"' STRING '"'
STRING <- CHARS_IN_STRING*
CHARS_IN_STRING <- ' ' / '!' / '#' / '$' / '%' / '&' / '(' / ')' / '*' / '+' / '-' / '<' / '>' / '=' / ',' / '.' / ':' / ';' / '?' / [0-9] / [a-z] / [A-Z]
INT -> [0-9]+

----------------------------------------------------TODO: MAP, ARRAY-----------------------------------------------------------
ARRAY -> '[]' TYPE '[' VALUE ']'
FI1(ARRAY) = {'[]'}

MAP -> '(' KEY_VALUE_PAIRS ')'
FI1(MAP) = {'('}

----------------------------------------------
KEY_VALUE_PAIRS -> STRING ':' VALUE | KEY_VALUE_PAIRS ',' KEY_VALUE_PAIRS
----------------------------------------------
KEY_VALUE_PAIRS -> STRING ':' VALUE KEY_VALUE_PAIRS'
KEY_VALUE_PAIRS' -> ',' KEY_VALUE_PAIRS KEY_VALUE_PAIRS' |  ε
----------------------------------------------
FI1(KEY_VALUE_PAIRS) = {'"'}
FI1(KEY_VALUE_PAIRS') = {','} + FO1(KEY_VALUE_PAIRS') = {','} + FO1(KEY_VALUE_PAIRS) = {','} + {')'}


SPECIAL_CHARS -> \r | \n | \t
```

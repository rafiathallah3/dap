# DAP

A friendly pseudocode that I implemented, inspired by my university professor, Mr. Jimmy. To prove that I can be like him. One day.

Example code:

```javascript
program ThisIsAProgram
dictionary
    n : integer
algorithm
    n <- 1

    for i = 1 to 10 do
        n <- n * i
    endfor

    print n
endprogram
```

<hr/>
The bug I can't fix occurs when you add a space at the end of a word before a newline; the newline token is ignored. The cause might be the regular expression in lexer.go
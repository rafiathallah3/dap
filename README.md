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
Bug that I can't fix is if you add space at the end of the word before the newline, the newline token would be ignored. The cause might be the Regular Expression in the lexer.go
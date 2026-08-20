# Control Flow

DAP supports comprehensive control flow structures: branching conditions, multiple loop formats, and pattern-matching switch blocks.

---

## 1. Conditional Statements (`if ... then ... else`)

### Multi-line If-Elif-Else
```javascript
program GradeEvaluation
dictionary
    score : integer
algorithm
    read score

    if score >= 90 then
        write "Grade A"
    elif score >= 80 then
        write "Grade B"
    elif score >= 70 then
        write "Grade C"
    else
        write "Grade F"
    endif
endprogram
```

### Single-line Short If
```javascript
if x < 0 then x <- -x
```

---

## 2. While Loop (`while ... do ... endwhile`)

Repeats statements as long as the condition evaluates to `true`.

```javascript
program WhileDemo
dictionary
    counter : integer
algorithm
    counter <- 1
    while counter <= 5 do
        print "Count: ", counter
        counter <- counter + 1
    endwhile
endprogram
```

---

## 3. For Loop (`for ... to / downto ... step ... do`)

The `for` loop in DAP is flexible and intuitive:

### Standard Increment (`to`)
```javascript
for i <- 1 to 10 do
    print i
endfor
```

### Custom Step (`step`)
```javascript
{{ Step by 2: 0, 2, 4, 6, 8, 10 }}
for i <- 0 to 10 step 2 do
    print i
endfor
```

### Decrementing Loop (`downto`)
```javascript
{{ Counts down: 10, 9, 8 ... 1 }}
for i <- 10 downto 1 do
    print "Countdown: ", i
endfor
```

---

## 4. Repeat-Until Loop (`repeat ... until`)

A post-tested loop that runs the body at least once and continues executing until the condition becomes `true`.

```javascript
program RepeatUntilDemo
dictionary
    num : integer
algorithm
    num <- 1
    repeat
        print "Value: ", num
        num <- num * 2
    until num > 50
endprogram
```

---

## 5. Switch-Case (`depend on ... enddependon`)

DAP provides two styles for `depend on` blocks:

### Style 1: Value Matching
```javascript
depend on (dayNumber)
    1 : print "Monday"
    2 : print "Tuesday"
    3 : print "Wednesday"
    default : print "Other day"
enddependon
```

### Style 2: Expression / Condition Matching
```javascript
depend on
    score >= 90 : print "Excellent"
    score >= 75 : print "Good"
    score >= 50 : print "Passing"
    default : print "Needs Improvement"
enddependon
```

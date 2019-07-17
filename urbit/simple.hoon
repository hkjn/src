:: create a gate to act as naked generator
::
:: the `|=` rune creates a gate, first child defines what kind of input it takes,
::   second input defines what to do with the input
::
:: first and second children are separated by "gaps", or double spaces
::
:: the first child defines that input is of any noun type and is referenced as `a`
|=  a=*
:: the second child produces that same input noun unchanged
a

:: create gate, and specify it takes atom type (`@`)
:: the `|=` rune creates a gate, first child defines what kind of input it takes,
::   second input defines what to do with the input, the input is of type atom
::   and is referenced as `end`
|=  end=@                                               ::  1
:: the `=/` rune stores a value with a name and specifies its type, taking three
::   children, the first child stores `1` as `count` and specifies that its of
::   type atom
=/  count=@  1                                          ::  2
:: the `|-` rune is a 'restart' point for recursion, used later
|-                                                      ::  3
:: the `^-` rune constrains output to list of atoms, and takes two children
^-  (list @)                                            ::  4
:: the `?:` rune evaluates whether the first child is true or false, and if
::   true branches to the second child, if false branches to the third child, it
::   takes three children, the `=(end count)` checks if the user's input equals
::   the `count` value that is being incremented to build the list, the
::   `=(end count)` is an irregular form of `.= end count`, and the `.=` rune
::   checks for the equality of its two children, and produces true/false
?:  =(end count)                                        ::  5
:: the `~` represents zod/null, and the L5 code branches here if `end == count`,
::   since we need lists to end with `~` in Hoon that's what we're adding here
  ~                                                     ::  6
:: the `:-` rune creates a cell, an ordered pair of two values, taking two children,
::   here creating a cell from whatever is stored in `count`, and with the product
::   of L8
:-  count                                               ::  7
:: the last line is a compact way of writing a rune expression, restarting the
::   program at L3, but with the value of `count` incremented by one, using the
::   `add` gate from the Hoon standard library
$(count (add 1 count))                                  ::  8
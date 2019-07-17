:: create a gate that takes two atoms and adds them
::
:: call this as `+add [3 9]` in the dojo, since naked generators
::   are limited to taking only one argument, so we define that the
::   gate takes a cell where we define "faces" `a` and `b`, which
::   we then add
|=  [a=@ud b=@ud]
(add a b)

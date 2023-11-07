# What is a string, in each language

## Definitions

A language has "Unicode strings" iff strings are sequences of Unicode scalar
values, and nothing else. I.e., you can represent valid Unicode text, and only
valid Unicode text, in the language's string type.

Really, you want a good interface that clearly distinguishes these cases:

* The length of the string in some encoded form (e.g., the byte length)
* The ability to iterate through scalar values (e.g., for parsing)
* The ability to operate on grapheme clusters (for text processing)

But many languages fail to provide scalar values, and usually grapheme clusters
are then out of the question. *Many* languages fail to provide basic type
safety, permitting hot garbage as strings.

Two common failures are worth setting out here:

* Sequence of code points: this permits surrogate code points to appear in the
  sequence.

  This does *not* correspond to WTF-8: such values will not necessarily
  round-trip: WTF-8 will encode code points outside the BMP and paired
  surrogates the same, meaning they cannot decode differently.

  Note that this model permits `[💩]` and `[U+d83d, U+dca9]` as distinct
  values.

* Sequence of UTF-16 code units: this permits lone surrogates. *This is
  different from a sequence of code points*, in that a sequence of code points
  permits `[💩]` and `[U+d83d, U+dca9]` as distinct values; sequence of UTF-16
  code units can only represent "💩" as a surrogate pair. However, this failure
  mode still permits *lone* surrogates.

  (To be clear: a language using UTF-16 as its internal encoding of strings
  does not fall into this category if it permits only UTF-16 strings.)

  WTF-8 corresponds roughly to this case. Any string from this case can be
  encoded into WTF-8. If one decodes the WTF-8, one can then take the resulting
  sequence of code points and transform them into UTF-16 code units by encoding
  any code point that is ≤0xffff into UTF-16 as-is, and encoding any code point
  \>0xffff into UTF-16 "normally". The result is the input.

Lastly, note that we're looking at what values are represented by the type. The
underlying encoding itself doesn't matter, so some of these might store these
values in memory as UTF-8, but only sometimes. We don't care about the memory
representation, rather, we care about the set of values representable by the
type.

## Language Overview

| Language | Strings Are | Well Specified? |
| -------- | ----------- | --------------- |
| Bash     | ❌ Even byte strings are hard | … |
| C        | ❌ Even byte strings are hard | If you like reading |
| C++      | ❌ Even byte strings are hard | If you like reading |
| Go       | ❌ Byte string | ❌ Not really |
| Haskell  | ❌ Sequence of code points | ✅ Yes |
| JavaScript | ❌ Sequence of UTF-16 code units | ✅ Yes |
| Java     | ❌ Sequence of UTF-16 code units | ❌ No |
| Lua      | ❌ Byte string | ❌ No |
| Objective-C | ❌ Unclear, probably sequence of UTF-16 code units | ❌ No |
| Python   | ❌ Sequence of code points | ✅ Yes |
| Rust     | ✅ Unicode  | ✅ Yes |
| Swift    | ❌ Unclear, probably sequence of UTF-16 code units | ❌ No |


## Notes & citations for each language

### Bash

Bash's string type can't represent the zero byte.

(Note this isn't true of all shells. `zsh`'s strings are byte strings.)

There's no Unicode anything to speak of, here. You're almost certainly better
served by at least Python.

### C

Just see C++, it has basically all the same problems. Yes, yes, they're
different languages, go away.

### C++

C++11 introduces `U'…'`, which can represent a code point, and `U"…"`, but
`U"…"` is a sequence of `char32_t`, which probably isn't what you want.

C++11 adds `u8"…"`, but until C++20, represents it as `char[]`. 🤦

C++17 introduces the `u8'…'` character literal, however, it is only capable of
representing code points < U+0080, and prior to C++20, represents them in the
wrong type `char`. 🤦

C++20 introduces the `char8_t` type, intended to hold UTF-8 code units. This
makes representing UTF-8 strings somewhat simpler.

Why not `char`? `char`'s signed-ness is implementation defined, making it hard
to work with correctly. Moreover, it is *signed* on common platforms such as
amd64. This makes doing things such as checking if the character is a
continuation byte or lead byte by doing `ch & 0x80` erroneous: 0x80 is not
representable in a `char`, if it is signed! But C++'s (i.e., C's) integer
promotion rules mean *that code is well-formed*, so no compiler errors, and it
is trivially 0 … obviously not the intention.

There's a lot more to cover here but … yeah C++ is 🦆ed.

TODO: the actual spec.


### Go

[String types](https://go.dev/ref/spec#String_types)

> A string value is a (possibly empty) sequence of bytes.

Clearly specified, if disappointing.

Go does include, in its standard library, tooling for *working* with UTF-8 byte
sequences, and through that, `string`s holding strings are generally expected
to be in UTF-8. E.g., you can quickly check if a `string` contains UTF-8 with
`utf8.ValidString()`.

That said, Go's `string` type squarely misses the mark, here.

… now, while the above spec would have earned an easy "✅" for being
well-specified … there are some other parts to Go's spec.

First, you can't do this:

```go
var s string = "\ud83d"
```

…or you'll get:

> prog.go:9:40: escape is invalid Unicode code point U+D83D

… but that *is* a valid code point. It's not a valid *scalar*.

The spec says …

> The escapes `\u` and `\U` represent Unicode code points so within them some
> values are illegal, in particular those above `0x10FFFF` and surrogate
> halves.

*sigh*.

So, you have,

> A rune literal represents a rune constant, an integer value identifying a
> Unicode code point.

So, given that `\u` seems to use "code point" to mean "scalar value", then a
`rune` is a scalar value?

```go
var r rune = 0xd83d
```

But … no… it's perfectly acceptable to assign a surrogate code point to a
`rune`.

Also, rune. *RUNE.* Go: why are we inventing words?! Just call it `code_point`.
Except…

```go
var r2 rune = 0xfffffff
```

No, it's not even that. It's just …

> alias for int32

Oh God, it's signed.

Anyways, Go's utter abuse of terminology and "rune" remove my ability to say
it's well specified.


### Haskell

* [String](https://hackage.haskell.org/package/base-4.19.0.0/docs/Data-String.html)
* [Char](https://hackage.haskell.org/package/base-4.19.0.0/docs/Data-Char.html#t:Char)

> `type String = [Char]`
>
> A `String` is a list of characters.

Clear enough; what's a `Char`?

> `data Char`
>
> The character type `Char` represents Unicode codespace and its elements are
> code points as in definitions [D9 and D10 of the Unicode
> Standard](https://www.unicode.org/versions/Unicode15.0.0/ch03.pdf#G2212).

Again, clear enough. Sequence of code points.


### Java

[String](https://docs.oracle.com/javase/8/docs/api/java/lang/String.html)

Defined as Unicode:

> A `String` represents a string in the UTF-16 format

Except in practice, permits lone surrogates, making this language "Sequence of
UTF-16 code units". See `counter-examples/java/Example.java`


### JavaScript

[The String Type](https://tc39.es/ecma262/multipage/ecmascript-data-types-and-values.html#sec-ecmascript-language-types-string-type)

> The *String type* is the set of all ordered sequences of zero or more 16-bit
> unsigned integer values (“elements”) up to a maximum length of 2⁵³ - 1
> elements. The String type is generally used to represent textual data in a
> running ECMAScript program, in which case each element in the String is
> treated as a UTF-16 code unit value.

Overly wordy, but precise: sequence of UTF-16 code units.

Example of an non-Unicode string in JavaScript:

```javascript
var s = '💩'[0];
console.log(s);
```


### Lua

[Values and Types](https://www.lua.org/manual/5.0/manual.html#2.2)

> String represents arrays of characters.

Without a definition of "character".

> Lua is 8-bit clean: Strings may contain any 8-bit character, including
> embedded zeros ('\0') (see 2.1).

One is left, I guess, to assume that "8-bit character" means "octet". This
seems to hold in practice.


### Objective-C

<https://developer.apple.com/documentation/foundation/nsstring>

Apple's documentation has the unique distinction of crashing, and failing to
display any documentation:

> Uncaught DOMException: The operation is insecure.

Yes, welcome to The Future, where documentation can crash. If you accept
cookies (which is what the developer is attempting to convey by dumping cryptic
errors into your browser's debugger) and venture onwards…

It seems like an `NSString` can be infallibly instantiated
from an array of UTF-16 code units, implying "sequence of UTF-16 code units".

TODO: counter-example?


### Python

[Text Sequence Type — `str`](https://docs.python.org/3/library/stdtypes.html#text-sequence-type-str)

> Strings are immutable sequences of Unicode code points.

Clearly written, though sadly sequence of code points.


Example showing non-Unicode strings, as well as non-UTF-16 code unit behavior:

```python
In [1]: '\ud83d\udca9' == '💩'
Out[1]: False
```


### Rust

[Primitive Type `str`](https://doc.rust-lang.org/std/primitive.str.html)

> String slices are always valid UTF-8.

Concise and to the point: Unicode!


## Swift

<https://developer.apple.com/documentation/foundation/nsstring>

Same bad documentation as Objective-C.

Even if you account for that … boy oh boy was some ink spent on this.

> Strings in Swift are Unicode correct

Implies they're Unicode.

But then…

> The `String` type bridges with the Objective-C class `NSString`

What does "bridge" mean here? It seems like you can infallibly create a
`String` from an `NSString`, so that implies that, iff `NSString` is Unicode,
then `String` is.

However, the Obj-C section has issues here. So presumably, a counter-example is
possible.


## Other

* JSON: Poorly specified: the grammar permits lone surrogates, so in a sense,
  Sequence of UTF-16 code units. The RFC says the behavior of strings that
  aren't valid Unicode "is unpredictable", but doesn't seem to *disallow* them.
* MySQL: MySQL has a litany of character encodings that aren't Unicode,
  including the deceptively named `utf8` which will reject valid UTF-8 strings.
  Just use PostgreSQL.
* SQLite: `text` is a truly horrific combination of octets *and* Unicode scalar
  values. Counter-example: `INSERT INTO test_table VALUES (x'ff' || 'bork');`
  will insert the sequence [the byte 0xff] || "bork" into a table.
  TODO: Surrogate behavior?
* WASM: (`stringref`) — a sequence of Unicode code points …
  [maybe?](https://github.com/WebAssembly/stringref/issues/65)
* YAML: Poorly specified: things are "Unicode characters" without ever giving a
  definition of such.


## The Real 💩

To me, this represents a real, festering problem within the industry of
software engineering: despite having a standard that defines text, and
captures, more or less, all the world's languages — quite a feat, and not
without debate about how well they did on that account — we as software
engineers fail to model strings.

If we can't model strings — what most languages call a "primitive" type! — …
what do you think our chances are of modelling anything more complex?

The adjacent conversation is that of much of the industry derides Computer
Science — the very science we allegedly engineer atop! — as being not useful.
Much of the above are painful to anyone with knowledge of type theory, as the
problems that arise from having, say, multiple representations for what, at a
human level, is really the same value, are obvious.

Lastly, the argument "but! but! UTF-16 happened *later*, and many of these are
really UTF-16!" is often just bull💩: for many languages above, the arrow of
causality would point the wrong way in that argument.

All of this means that interop is fraught with peril. The `stringref` type in
WASM is a good example of this, and I wonder if WASM will ever decide what a
"string" is, so long as it tries to please everyone in the cacophony of
indecision above.

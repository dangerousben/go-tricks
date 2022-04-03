# Techniques, patterns and hacks to make golang more bearable

I am not a great fan of go.  I have much material for a blog post on the subject which I will probably never get around to writing.  But in the meantime, I'm being paid to write the stuff, and while I will continue to protest that it is not the right tool for the job, I am unlikely to be able to escape it entirely. 

So I hereby present some ways of doing things that are otherwise difficult, ugly, or annoying to achieve in go.  The last language I remember feeling the need to develop a bunch of tricks like this for was PHP[^1], which I think says a great deal in itself.

[^1]: PHP 4 and 5 that is: I have no experience with versions after 7 and for all I know it's now a fantastic language

## General approach

I'm a functional ~~alcoholic~~ programmer, in both broad and narrow senses.  I don't have a hard position on the dynamic vs static typing debate, but I do believe that strong, featureful type systems are an enormous benefit to programmers, especially if checked statically.  The worst of both worlds is probably a weak static type system that offers a false sense of security while making the programmer's life more difficult for no obvious benefit.

The patterns aspire to the following qualities, in roughly descending order of importance:

- safety
- composability
- readability

These tricks will be based on go 1.18, released relatively recently at the time of writing.  This is a trade-off: I have no doubt that many sensible people are holding off an upgrade to 1.18 until it's had some time to bed in, but OTOH generics make such a huge difference to the langauge that avoiding their use would make several of these techniques obsolete right out the door.  I *may* choose to attempt alternative generic-free solutions to some problems, if there is a particularly good reason (such as doing it for the lulz).

## Patterns

(currently just list of random ideas)

### Channels

- fan out based on a tag extraction function
- fan in, with optional conversion (or should that be separate)
- priority
- functor and other typeclass support (except they're not really those typeclasses because the structures are mutable, maybe just use interfaces (maybe with the dawn of generics we'll get a standard or semi-standard set of such interfaces))
  - map
  - filter
  - folpd
  - zip

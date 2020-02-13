# go-workshop

Introductory workshop on syntax, structure, and concurrency in Go.

We will be making a simple look-up program that times
the user out after a set amount of time.

## Setup
  - Go can be downloaded [here](https://golang.org/dl/)
    - Go into the `skeleton` directory and make sure you can run `go run main.go`
  - Alternatively, you can also run Go online [here](https://repl.it/languages/Go)

## What is Go?
  - Statically-typed and compiled language similar to C or Java
  - Relatively young, started out at Google around 2007
  - Emphasis on speed and simplicity

## Hello World!

The entrypoint of any program is a `main` function in the
`main` package. We have also imported `fmt` which lets us
do some simple IO.

    package main
    import "fmt"
    
    func main() {
      fmt.Println("Hello world!")
    }

## Syntax Crash Course

Syntax is similar but slightly different from other languages.

Variables are declared with the type after the name. You can also
have the Go compiler infer the type.

    // declare and initialize
    var a string
    a = "Hello world"

    // or
    a := "Hello world"

One of the quirks with the Go compiler is that all declared variables
must be used somewhere in the program. We
can cheat this rule for now by assigning to the blank identifier.

    _ = a

Conditionals look very similar to C/Java but are more particular with
format.

    if 1 > 0 {
      fmt.Println("This is true")
    } else {
      fmt.Println("This is false")
    }

Looping in Go is a little unique.
There is only one loop in Go, the `for` loop.
    
    // infinite loop
    for {
      fmt.Println("looping")
    }

    // 'while' loop
    for a < 100 {
      fmt.Println(a)
      a++
    }

    // 'for' loop
    for i := 0; i < 100; i++ {
      fmt.Println(i)
    }

Functions are what you would expect as well.
    
    func multiplyByTwo(a int) int {
      return a * 2
    }

    multiplyByTwo(2)    // 4

## Data Structures

Let's make some simple structs for out lookup program!

    type human struct {
      name  string
      age   int
    }

    type dog struct {
      name  string
      owner *human
      age   int
    }

`*human` is a pointer type to the `human` struct which stores the owner as a reference instead of the whole object. We like this because it saves a lot of space.

Now let's initialize some structs.

    alice := human{"Alice", 56}
    bob := human{"Bob", 9}
    missile := dog{"Missile", &bob, 2}

We can identify our structs under a common type with an interface.
Since we are making a lookup table, let's make this interface
require a `getKey` method.

    type entry interface {
      getKey()  string
    }

Interfaces in Go are assigned implicitly. This means that all a struct needs
to do to implement an interface is implement all the methods inside it. We can
do this by writing function receivers for our structs.

    func (h human) getKey() string {
      return h.name
    }
    
    func (d dog) getKey() string {
      return d.name
    }

    alice.getKey()    // "Alice"

This code is kind of redundant though... Can we somehow condense it?

Nope!

Go isn't really object-oriented. It also doesn't support generic typing.

Let's create the actual lookup table as a map. A map can be initialized
like any other struct but we can use the `make` shorthand to *make* our lives
easier.

    lookup := make(map[string]entry)

    lookup[alice.getKey()] = alice
    lookup[bob.getKey()] = bob
    lookup[missile.getKey()] = missile

There's a bit of redundancy we can reduce here. Let's put these structs
in an array and loop through them instead. We can design array iteration
very quickly using the `range` keyword which will get us the index and
value of each element. Since we don't need the index, we can use the
blank identifier to indicate so.

    entryList := []entry{alice, bob, missile}
    for _, entry := range entryList {
      lookup[entry.getKey()] = entry
    }

Note that we could have also initialized an array with `make([]entry, 3)` for example, but we chose to do it manually to put our elements inside it.

From here, we can write a simple lookup loop that gives the user
10 queries and lets them find objects from our table. We will use
`fmt.Scanln` to read user input into a variable. We can also write
very concise map lookups in Go by taking advantage a second, optional
output which tells you whether the key was found or not.

    for i := 0; i < 10; i++ {
      fmt.Printf("You have %d queries left...\n", 10-i)
      fmt.Print("Enter: ")
      
      var input string
      fmt.Scanln(&input)
      
      if entry, ok := lookup[input]; ok {
        fmt.Println(entry)
      } else {
        fmt.Println("No entry found")
      }
    }

## Concurrency with Goroutines

What if we want to limit our user on the amount of time they
can spend using our program? This is a problem that can be
solved nicely by bringing some concurrency into the program.
If we were in another language and wanted to add
concurrency, it would take a while and a lot of screaming probably.

However, in Go, concurrency is built-in with Goroutines.
A Goroutine can be thought of as an independently running thread that's
able to share the resources in our program.
So far, we have been writing with a single Goroutine, the main function,
but we can spawn more Goroutines off of it.

We can make a Goroutine from any function call by putting `go` in front of it.

Yeah. It's really simple.

Let's turn our lookup loop into a Goroutine. We can do this with an
anonymous function as well since we only need to use it once.

    go func() {
      for i := 0; i < 10; i++ {
        // ...  
      }
    }()
    
We need another Goroutine to keep track of time. We can do
this with a `time.Sleep` from the `time` library.

    import "time"

    // ...

    go func() {
      time.Sleep(10 * time.Second)
    }()

If you run the program at this point, it looks like none of our
Goroutines ran. That's because our main thread never stopped running!

We need some way of blocking the main thread and also communicating
with it, and for that we turn to channels. A channel is like a tube
that you can send messages through. You can use `<-` to put a message
in a channel and `->` to pull it out.

Let's make a channel and
modify our two Goroutines to send different messages when they finish.
    
    c := make(chan string)

    // ...

    go func() {
      time.Sleep(10 * time.Second)
      c <- "\nTime up!"
    }()

    go func() {
      // ...
      c <- "Session ended"
    }()

We also need to make sure our main thread is blocking and waiting
for the message in the channel. Let's print it out as well.

    msg := <-c
    fmt.Println(msg)

With that, we have added concurrency in our lookup program with
very little effort!

## What's Next???

This is not the end!
If you think you like Go, there is more!

  - [Tour of Go](https://tour.golang.org/list): Basics of Go covered in greater detail
  - [Gophercises](https://gophercises.com/): Exercises in Go beyond the basics
  - [Go Time](https://changelog.com/gotime): Go-oriented podcast. Nice discussions about Go itself and the community.
  - [r/golang](https://www.reddit.com/r/golang/): Subreddit community
  - Start your own project!
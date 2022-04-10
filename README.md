# SSMutex

## Locks as container

ssmutex is a simple implementation for more ergonomic, and simpler to not misuse Locks in Go.

The main idea behind the package is making the Lock a **container** for the value that is being locked.
Making the lock a container allow to type check function boundaries.

## Motivation

While working with shared state, it is common to have functions that take as input some shared state.

It is important to know when the lock should be acquired and by whom.

Should the lock be acquired before calling the function by the caller, or the function itself should
acquire the lock? Using locks as containers allow to immediately answer this question, by looking only to the
function declaration.

If the argument of the function is a lock-container, the function itself need to acquire the lock (to get access
to the shared state.)

If the argument of the function is directly the shared state, it means that the lock was already acquired.

## How to release the lock

Using locks as container allows forcing the correct acquisition of locks.

However, it does nothing to force the correct release of locks.

Unfortunately I was not able to figure out a completely satisfactory answer to the problem, however the
current interface provides a reasonable solution.

The Lock() methods returns the object being locked AND a struct (a key) to use to unlock the mutex.

Only by using the key it is possible to unlock the mutex.

In order to obtain the value a second time, you are force to use the key once.

Each key can be only used once.

If a key is leaked (and the GC run) the software panics.

Enforcing the usage of the key could be done at LINT time, as we already do with context.Context and the
cancel function.

## Caveats

The package does not prohibit the leak of the shared state.

However, it should be immediately evident.
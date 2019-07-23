This library is being archived because it's naive in the following ways:

* Most logs are collected as metrics, so color doesn't make sense
* Logrotation (and in fact, any log storage at all) is usually handled by a service listening to stdout
* It's dependent on the stdlib's "log" as a backend, which may not be the best backend (isn't the best backend)

It's being replaced by my LoggerInterface (some undetermined library in github.com/ayjayt) that just requires two logging levels:
* Info
* Error

This is inspired by Dave Cheney's blog post. The basic idea is that:
1) There is no such thing as a warning. It is either an unhandlded error or error that _needs_ to be address, or it is just info.
2) I don't like a debug level: It adds lots of clutter and tons of lines of code that should be removed before deploying. Use a real debugger or something.

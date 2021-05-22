# 👾 Pron 👾
## Overview
An extensible personal cron tab, supporting external shell commands, as well as internal golang functions.

## Plans
I'm thinking about adding bindings for embedded language support; however, I may not do this as arbitrary code can be run externally anyhow - though this would allow for injection and interception on a finer grained level.

## Todo
* Fix issue where everything works except for the recurring arbitrary code execution... lol... the whole point of the program
* Enable smart recovery of panics triggered by the ticking mechanism
* Finish the logging module
* Write a daemon wrapper
* Write an external monitoring socket + client

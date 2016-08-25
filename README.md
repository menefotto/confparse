This is a small package the provides a ini style configuration parser. This is 
what is allowed:

- Comments, either with the "#" or ":" anything after it till newline is ignored
- Sections like [default]
- Key and values like "ip=192.168.10.1"

Empty line are ignored, whitespaces are ignored as well. It hasn't really been
tested yet, it might have some bugs.

Suggestion: do not use this software in a production enviroment, it's not even
completed yet and hasn't really been tested.

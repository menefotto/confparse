# YACP 

[![GoDoc](https://godoc.org/github.com/spf13/viper?status.svg)](https://godoc.org/github.com/spf13/viper)
[![Build Status](https://travis-ci.org/wind85/confparse.svg?branch=master)](https://travis-ci.org/wind85/confparse)
[![Coverage Status](https://coveralls.io/repos/github/wind85/confparse/badge.svg?branch=master)](https://coveralls.io/github/wind85/confparse?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
### Yet Another Configuration Parser
This is a small package the provides a ini style configuration parser. This is 
what is allowed:

- Comments start, either with the "#" or ":" anything after it, till newline is ignored
- Sections are written like the following [default] and contain a map of key values,
  anything between square brackets is a valid section.
- Key and values are like "ip=192.168.10.1" ,the separator is "=" otherwise will
  not be considered a key value.
- The Parser can handle bools, ints and floats (both 64bit), strings and string slices,
  as long as they are divided by columns.
- Empty line are ignored, white spaces are ignored as well. It hasn't really been
tested yet, it might have some bugs.

### How to use it
Pretty simple, there are only two methods to create a new configuration either call 
```
  ini := confparse.New()
  ini.Parse()
```
or
```
  ini, err := confparse.NewFromFile("config-name.whatever")
```
The second version does call Parse under the hood, and it isn't name sensible any valid name
can be passed. Then any of the valid supported values can be retrieved like so:
```
  value ,err := init.GetInt("sectionname","valuename")
  value ,err := init.GetFloat("sectionname","valuename")
  value ,err := init.GetSlice("sectionname","valuename")
  value ,err := init.GetString("sectionname","valuename")
  value ,err := init.GetBool("sectionname","valuename")
```
And that's pretty much it.

#### Disclaimer
This software in alpha quality, don't use it in a production environment, it's not even
completed yet and hasn't really been tested.

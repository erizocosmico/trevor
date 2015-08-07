# trevor [![Build Status](https://travis-ci.org/mvader/trevor.svg)](https://travis-ci.org/mvader/trevor) [![Coverage Status](https://coveralls.io/repos/mvader/trevor/badge.svg?branch=master&service=github)](https://coveralls.io/github/mvader/trevor?branch=master) [![GoDoc](https://godoc.org/github.com/mvader/trevor?status.svg)](http://godoc.org/github.com/mvader/trevor)
Trevor is an extensible framework to build [knowledge navigators](https://en.wikipedia.org/wiki/Knowledge_navigator) through custom plugins.

In brief, you could build a Siri-like or Google Now-like API (without the voice recognition layer) for you app. Trevor only provides the framework to build it.

## How does it work

* A request is made to an endpoint (configurable via [Config](http://godoc.org/github.com/mvader/trevor#Config)) with a JSON like:
```json
{
  "text": "recommend me a movie"
}
```
* The server collects that text and gives it to all available plugins. All plugins return a score.
* With the list of scores received and the preference of the plugins (you can add a number to represent the preference. Higher is better) it chooses the best candidate by sorting by exact match (the input received is an exact match of a rule in the plugin, meaning it's a perfect match), score and preference. That means that a plugin with preference 3 and a score of 5 will be selected over a plugin with preference 10 and score 0.
* With the best candidate selected the text will be given to that candidate and it will respond with data.
* The server will output the data in an output like:
```json
{
  "error": false,
  "type": "plugin name",
  "data": "<whatever, can be an array, an object, a string, ...>"
}
```

## Working with plugins

To create a plugin you just have to implement the [Plugin](http://godoc.org/github.com/mvader/trevor#Plugin) interface.

**Considerations:**
* Though there are no limits for the score (only that it has to be a float64) it is recommended to use a range of numbers consistent across all plugins. **The recommended range is [0, 10]**.
* Only return an exact match if it really is an exact match and you know no other plugin could have a better result. If the result is an exact match it will be on top of all other results. For example, if you search for a movie called "Lost" and you have a movie called "Lost puppies" don't return an exact match because the shows plugin will have a "Lost" show that really is an exact match.
But thing that there is also a search plugin. Neither of the aforementioned plugins could return an exact match for the text "lost" because maybe the user is not looking for shows.
**TL;DR:** think your plugins very well to work nice with other plugins.

Check out [trevor-plugins](https://github.com/mvader/trevor-plugins) for reference plugins.

## Example

Example implementation. We consider a fictional plugin `movie` that lives in `github.com/trevor/movie` (this plugin does not actually exist).

```go
package main

import (
  "github.com/mvader/trevor"
  "github.com/trevor/movie"
)

func main() {
  server := trevor.NewServer(trevor.Config{
    Plugins: []trevor.Plugin(NewMovie()),
    Port:    8888,
  })

  server.Run()
}
```

See? Easy peasy.

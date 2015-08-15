![trevor](https://raw.githubusercontent.com/mvader/trevor/master/trevor.png)

[![Build Status](https://travis-ci.org/mvader/trevor.svg)](https://travis-ci.org/mvader/trevor) [![Coverage Status](https://coveralls.io/repos/mvader/trevor/badge.svg?branch=master&service=github)](https://coveralls.io/github/mvader/trevor?branch=master) [![GoDoc](https://godoc.org/github.com/mvader/trevor?status.svg)](http://godoc.org/github.com/mvader/trevor)

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

## Create plugins

To create a plugin you just have to implement the [Plugin](http://godoc.org/github.com/mvader/trevor#Plugin) interface.

**Considerations:**
* Even though there are no limits for the score (only that it has to be a float64) it is recommended to use a range of numbers consistent across all plugins. **The recommended range is [0, 10]**.
* Only return an exact match if it really is an exact match and you know no other plugin could have a better result. If the result is an exact match it will be on top of all other results. For example, if you search for a movie called "Lost" and you have a movie called "Lost puppies" don't return an exact match because the shows plugin will have a "Lost" show that really is an exact match.
But thing that there is also a search plugin. Neither of the aforementioned plugins could return an exact match for the text "lost" because maybe the user is not looking for shows.
**TL;DR:** think your plugins very well to work nice with other plugins.

Check out [trevor-plugins](https://github.com/mvader/trevor-plugins) for reference plugins.

### Injectable pugins

A plugin may need some services to run. Maybe even two plugins need the same service. Instead of implementing each service in your own plugin trevor provides a mechanism to register services on the engine.
Using injectable plugins a plugin receives on startup time all the services it needs to work. To convert a plugin to an injectable plugin you just need to implement the [InjectablePlugin](http://godoc.org/github.com/mvader/trevor#InjectablePlugin) interface.

Basically, you implement the `NeededServices` method that returns a list with the names of all services needed by the plugin.
When the engine is started it will try to find those services and inject them to the plugin using the `SetPlugin` method of the plugin.

Note that only one service is injected at a time. If your plugin needs 4 services, the `SetPlugin` method will be called 4 times, every time with a different service.

Plugins receive a `trevor.Service`, you will have to cast them to the right type of the service you are using.
**Example:**

```go
func (p *myPlugin) SetService(name string, service trevor.Service) {
  switch name {
  case "redis_cache":
    p.redisService = service.(*MyRedisCacheService)
  break;

  case "memcached_cache":
    p.memcachedService = service.(*MyMemcachedCacheService)
  break;
  }
}
```

## Create services

To create a service you just have to implement the [Service](http://godoc.org/github.com/mvader/trevor#Service) interface.

All it is asked for a service to implement is `Name` and `SetName` methods. The rest is up to the developer of the service.

**Considerations:**
Use an unique name for the service. If you use the name "cache" it will sure clash with another service. Imagine a `RedisCacheService` and a `MemcachedCacheService`. If both are named `cache` only the last one added will be available in the engine. In that case, they should be named `redis_cache` and `memcached_cache`. Then, if the user wants to use them as `cache` they can be renamed with the `SetName` method.

## Memory service

The memory service is a special type of service, a service that also implements the [MemoryService](http://godoc.org/github.com/mvader/trevor#MemoryService) interface. The purpose of this kind of service is to give memory to the trevor engine. Not actual memory but the ability to "remember" an user. It works like regular authentication, given an HTTP request the service has a method to return a token based on that request. If that token is passed along in subsequent requests, trevor will be able to identify the user that is requesting information and then the trevor engine can give a more personalized response. For example, you could use that service to give better results based on what the user has previously requested.

How all of this is implemented is up to you, `MemoryService` is only defined as an interface in the core because it needs integration in the core and thus it needs a common interface.

If you need guidance on how to implement a memory service you can take a look at the [memory_service_test.go](https://github.com/mvader/trevor/blob/master/memory_service_test.go) file to see how the service is implemented for the tests.

### How it works
* The first time an user requests information no token is passed with the request.
* All the plugins receive the request.
* The plugin in charge of processing the request should assign a token to the [Request](http://godoc.org/github.com/mvader/trevor#Request) it received.
* The server will send the token previously assigned to the request with the response.
* In subsequent requests the user will pass the token with the request.

## Pokables

A `Pokable` in trevor is a component (a plugin or a service) that will be [poked](http://www.wanapesa.com/poke/img/94888877_o.png) in intervals defined by the same pokable.
To make a `Plugin` or a `Service` pokable the only thing you need to do is implement the [Pokable](http://godoc.org/github.com/mvader/trevor#Pokable) interface.

The motivation for the pokables is to have a centralized scheduler that will call the component every X time. That eliminates the need to have workers in most cases. For example, imagine you have a service that needs a configuration from a server but that configuration changes every 24 hours. You could implement a goroutine that fetches that configuration every 24 hours but that should not be a responsability of the service. Instead, you could make the service pokable and every 24 hours (if you define that interval in your implementation) the `Poke` method will be called. That way, your service does no longer have the responsability of spawning a goroutine to fetch periodically the configuration.

All poking goroutines are spawned when the server `Run` method is called.

## Custom analyzer

Maybe you want to ditch the default behavior of the trevor engine (iterate over all plugins to get the score returned of analysing the input and choosing the better match) and use your own analyzer function that decides which plugin should be used. You can do that by passing an [Analyzer](http://godoc.org/github.com/mvader/trevor#Analyzer) to the server on the configuration.

An `Analyzer` receives the input and returns the name of the plugin that will process that input and metadata just like the plugins `Analyze` method would.

#### Example

```go
func main() {
  config := trevor.Config{
    Plugins:  []trevor.Plugin{MyFancyPlugin(), MyOtherFancyPlugin()},
    Services: []trevor.Service{MyFancyService()},
    Port:     8888,
    Analyzer: func (text string) (string, interface{} {
      // wow such analysis
      return "plugin name", map[string]interface{}{
        "some": "values",
      }
    },
  }

  server := trevor.NewServer(config)
  server.Run()
}
```

## Example

Example implementation. We consider a fictional plugin `randomMovie` that lives in `github.com/trevor/movie` (this plugin does not actually exist). We also consider another fictional `cacheService` service that lives in `github.com/trevor/cache`.

**Note**: we use the service as it is. That means that its name is the one assigned by the person who created the service. If we know it might clash with another plugin you can always change it with the method `SetName` that every plugin has.

```go
package main

import (
  "github.com/mvader/trevor"
  "github.com/trevor/movie"
  "github.com/trevor/cache"
)

func main() {
  server := trevor.NewServer(trevor.Config{
    Plugins: []trevor.Plugin{movie.NewRandomMovie()},
    Services: []trevor.Service{cache.NewCacheService()},
    Port:    8888,
  })

  server.Run()
}
```

See? Easy peasy.

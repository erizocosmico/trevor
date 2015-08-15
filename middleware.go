package trevor

// Middleware is a function that wraps the process of the engine.
// Every middleware added is a new layer around the result. Its purpose
// is to be able to perform actions before and after the process
// has taken place.
// A middleware function receives three parameters:
// - The request
// - A function to retrieve services from the engine by their name
// - A function that will return the result of the next layer of middleware or, if there is none, the result of the process
type Middleware func(req *Request, getService func(string) Service, next func() (string, interface{}, error)) (string, interface{}, error)

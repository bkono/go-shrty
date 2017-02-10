## Shrty

Whenever I set out to properly learn a new language, I like to utilize a test-app with real world
use cases. A URL shortener, shrty, has become that go-to over the last few years. This iteration is
a pass utilizing mostly stdlib Go, but will be refactored to include gRPC and other libs, because I
can. 


### MVP

The minimum viable product (MVP) of shrty consists of:

1. Give a URL, get back a shortened version
1. Navigate to the shrt url, be redirected to the original

Thats it. Generally, I like to add a few other features like persistence through restarts (though
in-mem is often an acceptable first pass), tracking a count of views, etc. For this particular repo,
I'm going to iterate a bit. Pending tasks will be maintained in this readme.


### Additional Features (technical and product)

- [x] Persist with BoltDB, because Bolt!
- [x] Increment views counting per URL
- [x] Add cli flag for baseURL (flag parsing, gotta see what the stdlib offers)
- [x] Proto integ for serialized structs
- [ ] Migrate to something along the lines of [Standard Project Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1#.w26wk0yon)
- [ ] Don't recreate already known URLs
- [ ] gRPC endpoint for both create and expand
- [ ] gRPC client binary for launchbar integration
- [ ] JSON endpoint for fetching metrics
- [ ] Web views for metrics
- [ ] Swap to chi for routing. Do this one to see if chi is the right fit for other projects



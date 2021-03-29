# go-habitify

Github-like calendar heatmap (or whatever it's called) for Habitify data.

The root package exposes a very basic client for Habitify's REST API
currently exposing only a subset of their API.

* [x] `GET /journal`
* [x] `GET /habits` (Missing filters)
* [x] `GET /habits/:id`
* [ ] `GET /habits/:id/status`
* [ ] `PUT /habits/:id/status`
* [ ] `GET /habits/:id/logs`
* [x] `POST /habits/:id/logs`
* [ ] `DELETE /habits/:id/logs`
* [ ] `DELETE /habits/:id/logs/:id`
* [ ] `GET /habits/:id/notes`
* [ ] `POST /habits/:id/notes`
* [ ] `DELETE /habits/:id/notes`
* [ ] `DELETE /habits/:id/notes/:id`
* [ ] `GET /areas`

## Running

You can get the API key for your Habitify accoune from the apps' settings.

```sh
HABITIFY_API_KEY=<YOUR_API_KEY> go run ./cmd/main.go
```

![UI](README-ui.png)
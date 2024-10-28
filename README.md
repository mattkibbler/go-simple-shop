# go-simple-shop

## Demo
Demo available at: [https://go-simple-shop.mattkibbler.com](https://go-simple-shop.mattkibbler.com).

## About

A simple demo web storefront using product data from the [https://dummyjson.com/docs/products](https://dummyjson.com/docs/products) API.

Written using Go as I am enjoying writing more and more projects in Go and wanted to utilise this simple but powerful language to get as much as possible done in a very short time frame.

The application fetches all product data from the API and caches it locally so that we don't have the additional latency of making requests to the dummy server every time somebody visits our shop, there is also the possibility that the dummy API could be offline or have issues resulting in our own shop having issues.

A background goroutine refreshes the data every 10 minutes. Technically this could mean that our shop has out-dated information between refreshes, such as a product being removed or out of stock but this could be handled  when the user adds an item to their basket or checking out, at this point we can make a call to the dummy product server to verify that the product is available.

The product data is simply stored in a `sync.Map` in order to simply and efficiently handle concurrent reads and writes. The data is simply held in memory, which is fine for this demo application, but a more robust app would require storing using a real database or similar solution.

When the product data is refreshed we also serialize the data into a local JSON file, this is because our in-memory data store will be lost when the server is restarted and as we need to frequently restart the app during development this prevents us from hitting the product API every time, and it's handy to have product data available to us as soon as the server is online.

The storefront itself consists of a product listings page and a product details page. Each time a product or products are required we simply pull them from the in-memory store, which makes rendering pages very quick. In a real system we would probably want to implement some page caching as we would be pulling data from one or more sources such as a database, but thats out of scope for us here.

The frontend is very simple. We are using HTML and [TailwindCSS](https://tailwindcss.com) with some simple [Tailwind UI](https://tailwindui.com) components to display the products. Given the tight time constraints, there’s no advanced JavaScript beyond a small inline script on the product detail page for switching images, which is all we need for this lightweight storefront.

## Building / running

Clone the repo and copy the `.env.example` file to create your own local `.env` file. The files should have the following values:

| Key                | Value                                          |
| :------------------| :--------------------------------------------- |
| PORT               | The port number that the app should run on     |
| APP_ENV            | Should be "local" or "prod"                    |

Use the following commands to build/run the application.

| Command                | Action                                           |
| :--------------------- | :----------------------------------------------- |
| `make tailwind`        | Downloads the tailwind CLI binary. This should be ran just once.           |
| `make build`           | Builds the app and any associated assets.             |
| `make run`             | Runs the app via the `go run` command.       |
| `make watch`           | Uses `air` to watch source files and reload app when changes are made.       |
| `make docker-run`      | Builds and runs the app in a Docker container     |
| `make docker-down`     | Stops running Docker container     |

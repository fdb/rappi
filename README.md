Proxy around popular APIs that avoids the authentication mess.

## Development

Make sure you have a working Go environment set up, then use [gin](https://github.com/codegangsta/gin) for live reloading:

    go get github.com/codegangsta/gin
    gin

## Deploying

We deploy to Heroku.

    git commit -a
    git push heroku master


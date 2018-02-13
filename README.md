Proxy around popular APIs that avoids the authentication mess.

## Development

Make sure you have a working Go environment set up, then use [gin](https://github.com/codegangsta/gin) for live reloading:

    brew install golang dep
    cd $GOPATH
    mkdir -p src/bitbucket.org/fdb
    cd fdb
    git clone git@bitbucket.org:fdb/rappi.git
    cd rappi
    dep ensure

    go get github.com/codegangsta/gin
    gin

## Deploying

We deploy to Heroku.

    git commit -a
    git push heroku master


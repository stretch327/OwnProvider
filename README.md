# OwnProvider
An APNs Provider Server which based on JWT auth for iOS.  Build your own Apple Provider Server even you are not a developer.

# Env
go version go1.11 darwin/amd64

# Build & Install
```shell
go build -o ownprovider
```

```shell
/where/your/provider/is/installed/ownprovider
```

# Push a Message
```shell
curl -X POST "http://127.0.0.1:27953/api/notify" -d 'topic=YourBundleId&token=YourDeviceToken&payload=%7B%22aps%22:%7B%22alert%22:%22Hello%22%7D%7D'
```


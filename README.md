# OwnProvider
An APNs Provider Server which based on JWT auth for iOS.  Build your own Apple Provider Server even you are not a developer.

# Env
go version go1.11 darwin/amd64

# Build & Run
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

# URL Parameters
**env**        - Set this value to "sandbox" to use the development APNs server.<br>
**type**       - The type of message being sent. Can be one of "alert", "background", "voip", "complication", "fileprovider", "mdm".<br>
**token**      - The device token to send this message to.<br>
**payload**    - The APNs payload.<br>
**bundleid**   - The app's bundle identifier.<br>
**expiration** - The number of seconds from the current time when this notification expires and should no longer be delivered. Zero (the default) indicates that APNs should only attempt to deliver once. This value cannot be greater than 15777000 (6 months in seconds).<br>
**priority**   - The delivery priority. Can be any value from 0 (No Priority) to 10 (Highest Priority, Default).<br>
**collapseid** - An identifier for collapsing multiple push notifications in the notification center.<br>
**teamid**     - Your Apple Team Identifier.<br>
**keyid**      - Your APNs token key id.
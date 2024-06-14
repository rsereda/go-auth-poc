#Build 

```
cd auth
go build
```

#Using
```
./auth
```

test call

```
curl -X POST http://localhost:1323/auth   -H 'Content-Type: application/json'   -d @user.json
```
as e result should be new genereted token  


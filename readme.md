This is a simple HTTP web service implemented on Golang.

Web service is responsible for URL shortener duty. You can store new value by calling `PUT` and get a known record by calling `GET`.

To run this service type:
```
go run .\main.go
```

Sample 1. Save new value.
```
> curl -v -X PUT -d "addr=https://radio-t.com" http://localhost:8080

eyWT
```

Sample 2. Get previously saved value.
```
> curl -v -X GET http://localhost:8080/eyWT

<a href="https://radio-t.com">Moved Permanently</a>
```

For more information about URL shortening please visit [educative.io](https://www.educative.io/courses/grokking-the-system-design-interview/m2ygV4E81AR).
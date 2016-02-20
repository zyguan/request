# Usage

```go
urls := []string{"https://github.com/zyguan/request"}
for d := range request.Do(2, request.GetRequests(urls)) { ...
```

# Usage

```go
urls := []string{"https://github.com/zyguan/request"}
for d := range request.Do(request.GetRequests(urls), 2) { ...
```

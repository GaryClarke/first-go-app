<!-- File: docs/requests.md -->
### Check the app health
```bash
curl -i -X GET http://localhost:8080/healthz
```

### Get all books
```bash
curl -i -X GET http://localhost:8080/books
```

### Get a book by id
```bash
curl -i -X GET http://localhost:8080/books/999
```

### Create a new book
```bash
curl -i -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"The Go Workshop","author":"Delio D'\''Anna","year":0}'
```

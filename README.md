# Flix - Music trainer

This application is dedicated to a very good friend of mine, aka. "Flix".

## Development

```bash
# launch dev shell
nix develop

# apply code format
nix fmt
```

## Backend

### Test coverage report

After launching dev shell, just follow these steps to open the test coverage HTML report in your browser.

```bash
cd backend
cover-report

# FYI: 'cover-report' is just an alias in the dev shell, which combines the following commands
go test -coverprofile cover.out
go tool cover -html=cover.out -o cover.html
xdg-open cover.html
```

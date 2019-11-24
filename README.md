# XE
XE is a service that returns the value of 1 provided currency in euros. It also makes a recommendation
based on the last week rate history.

## Running the service
Running the application locally
```bash
make local_run
```

## Sending request to the service
Send a request with query param `currency`
```bash
curl -i localhost:3030/convert\?currency\=USD
```
Example response:
```json
{
  "from": "USD",
  "to": "EUR",
  "rate": 0.9043226623,
  "recommendation": "convert"
}
```
`rate` indicates the value of 1 USD in EUR, `recommendation` of "convert" means it's good to convert from USD to EUR.

## Checking test coverage
```bash
make cover && open coverage.html
```

## Running tests and linting
```bash
make all
```
which includes `make lint` and `make test`

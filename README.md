# Meisterwerk Quote

**Quote** is a web application to manage customer quotes. It lets you create, update, delete, and process quotes. You can also manage quote products, addresses, and payment information.

## Features

- Create, update, delete, and process customer quotes.
- Add, update, and remove products in quotes.
- Save customer addresses and payment details.
- Use external services for product catalog, tax calculation, and order processing.

## Installation

1. Clone the project:

   ```bash
   git clone https://github.com/yourusername/meisterwerk.git
   cd meisterwerk
   ```

2. Install dependencies:

   ```bash
   make deps
   ```

3. Run the web api application:

   ```bash
   make api
   ```

   The app runs on `http://localhost:8080`.

## Testing

Run tests with:

```bash
go test ./...
```

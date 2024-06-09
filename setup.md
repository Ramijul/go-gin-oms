## Running locally WITH docker-compose

- Ensure you have docker-daemon running
- Ensure you have docker-compose
- Build and run the images. This will also pre-populate the database with one user and one product:

  ```
  docker-compose up -- build
  ```

  _Note_ that it could take about 30 seconds for the API service to be active. It wait for all the services it dependents on to be active first!

- To take down all running containers, volumes and images

  ```
  docker-compose down -v --rmi all
  ```

## Running locally WITHOUT docker-compose

- Ensure you have postgres and rabbitmq running
- Generate `.env` files created in the root of `payments` and `orders` modules. Copy the items from their respective `env.example` file
- On terminal, navigate to `/orders` directory and run `go run ./migrate/migrate.go && go run .`, to get the API service up and running
- On another terminal, navigate to `/payments` directory and run `go run .`, to get the Payments Service up and consuming messages

## Testing the Order Management System

1. Grab the UserID from

   ```
   GET localhost:8080/users
   ```

   _This will list all the users. You will find only one for now._

2. Grab the ProductID from

   ```
   GET localhost:8080/products
   ```

   _This will list all the products. You will find only one for now._

3. Place an Order

   ```
   POST localhost:8080/order
   {
       "user_id": "<insert UserID>",
       "address": {
           "street": "Westdale",
           "zip": "L8S 1A8",
           "city": "Hamilton",
           "province": "Ontario"
       },
       "products": [{"id": "<insert ProductID>", "units": 15}]
   }
   ```

   _In Respone you will receive the newly created order record with a "PAYMENT PENDING" status_
   _Note_ that placing an order does not (yet) reduce the stock on the product.

4. Fetch the order

   ```
   GET localhost:8080/order/<insert OrderID>
   ```

   OR

   ```
   GET localhost:8080/orders
   ```

   _Upon processing the payment the Payment Status should reflect the confirmation_

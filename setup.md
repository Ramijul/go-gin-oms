## Running locally with docker-compose

- Ensure you have docker-daemon running
- Ensure you have docker-compose
- Build and run the images. This will also pre-populate the database with one user and one product:

  ```
  docker-compose up -- build
  ```

- To take down all running containers, volumes and images

  ```
  docker-compose down -v --rmi all
  ```

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

4. Fetch the order

   ```
   GET localhost:8080/order/<insert OrderID>
   ```

   OR

   ```
   GET localhost:8080/orders
   ```

   _This should show the expected payment status depending on the total order price_

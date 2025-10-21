# Mixie Melts V2

Yaaay! Everyone cheers!
This is my modern solution for the Mixie Melts site, now cloud-deployable on
any platoform Kubernetes is supported (Pretty Much all of them).

Or.. it will be when it's done! Right now you can run it with docker-compose:

```bash
docker-compose up --build
```

## Mixie Melts

Mixie Melts is a wax melts business that makes high-quality all-natural soy
wax melts. They sell the wax melts in 6-packs of single and multi-scents, as
well as do monthly subscription boxes that send customers wax melts that have
a unique theme for each month, like Mashed Potatoes and Gravy scented melts in
November. It also tracks inventory as customers place orders to ensure the
kitchen still has the scents in stock, and sends a ticket to the kitchen how
many wax melts of what scents need to be made.

## Back-End

This site's backend is a kubernetes application made up of the following
microservices.

1. Frontend Service: This controls the web interface that users interact with.
2. User Service: Handles customer registration, login, profiles, and auth.
3. Product Service: Handles wax melt products, details, pricing, and categories.
4. Inventory Service: A dedicated service to track stock levels. It gets
 updated when and order is place and when you add new stock.
5. Cart Service: Manages the shopping cart for each user.
6. Order Service: Processes completed carts, finalizes orders, and communicates
 and communicates with other services.
7. Subscription Service: Manages the logic for monthly subscription boxes,
 recurring billing, and customer subscription status.
8. Kitchen Notifier: A small service that listens for new order events and
 "dings" the kitchen.
9. Kitchen Frontend: The frontend for the kitchen ticket
 system, so the kitchen staff know what orders need
 to be made and can mark them off as completed in the
 system as they're made in batches.
10. API Gateway / Ingress: The single entry point for all web traffic, routing
 requests to the correct internal service.

### Frontend Service

This is a React based app themed with Tailwind CSS. It
handles all the user input as well as gives the user a
pleasing website to interact with.

### User Service

Ths applications handles all the customer data as well as
login and auth.

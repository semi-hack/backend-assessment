## Description

Backend-Assessment

## Documentation

https://documenter.getpostman.com/view/9562205/2sA2r9WicS

## DB Diagram

https://drive.google.com/file/d/1dNxAUIOJVEbFH84BZSbZ0T_HXNewNtzu/view?usp=drive_link

### Answers

 - when designing the db, i would have a merchants table and a product table, one merchant can have multiple products
   so a one to many relationship between merchant and products, id have my merchantId as a foreign key to a product
   using that i can map all products that belong to that merchant to the merchant using that foreign key

 - for this particular problem, speed, data complexity and consistency would inform me of what databse to use
   i would go for a sql db such as postgres, because i want to use of acid transactions to maintain consitency
   also the relationship between merchant and product is clear, the data is structured, so querying with a relational db
   such as postgres would be beneficial.







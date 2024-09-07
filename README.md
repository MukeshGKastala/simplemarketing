
# Simple Marketing

This project aims to enforce temporal uniqueness constraints at the API layer while ensuring data integrity for a RESTful promotions entity in a concurrent environment.

---

## Why This Project?

The main focus of this project is to explore:
- **Conditional Uniqueness**: When the uniqueness of a field (such as `promotion_code`) needs to be enforced conditionally, based on factors like whether a promotion is soft-deleted or has expired.
- **Use of Transactions**: Ensuring the proper use of transactions to maintain data integrity in the face of concurrent operations.
- **Transaction Isolation Levels**: Investigating how a database’s transaction isolation levels affect transactions when handling conflicting operations.

---

## Design and Architecture

This service relies on a MySQL database to store promotions.

```sql
CREATE TABLE promotions (
	id INT AUTO_INCREMENT PRIMARY KEY,
	promotion_code VARCHAR(255) NOT NULL,
	start_date DATETIME NOT NULL,
	end_date DATETIME NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at DATETIME DEFAULT NULL
);
```

## API Endpoints

- **Create a promotion**: 
  - Create a new promotion with a unique `promotion_code`
- **List promotions**:
  - Retrieve a list of all currently active (non-expired, non-deleted) promotions
- **Delete a promotion**:
  - Delete a promotion by `id`
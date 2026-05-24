-- name: CreateLoan :one
INSERT INTO loans (id, lender, borrower, book, lent_at)
VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    NOW()
) 
RETURNING *;

-- name: ReturnLoan :one
UPDATE loans
SET returned_at = NOW()
WHERE book = $1
RETURNING *;

-- name: GetLoanHistory :many
SELECT *
FROM loans
WHERE book = $1
ORDER BY lent_at DESC;
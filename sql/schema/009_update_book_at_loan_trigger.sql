-- +goose up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_book_availability()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE books
        SET is_available = FALSE,
            borrower = NEW.borrower,
            updated_at = NOW()
        WHERE id = NEW.book;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE books
        SET is_available = TRUE,
            borrower = NULL,
            updated_at = NOW()
        WHERE id = OLD.book;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_book_availability_trigger
AFTER INSERT OR UPDATE ON loans
FOR EACH ROW
EXECUTE FUNCTION update_book_availability();

-- +goose down
DROP TRIGGER IF EXISTS update_book_availability_trigger ON loans;
DROP FUNCTION IF EXISTS update_book_availability();
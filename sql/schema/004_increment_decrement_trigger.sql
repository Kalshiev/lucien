-- +goose up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION manage_book_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE collections
        SET book_count = book_count + 1
        WHERE id = NEW.collection_id;
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE collections
        SET book_count = book_count - 1
        WHERE id = OLD.collection_id;
        RETURN OLD;

    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.collection_id IS DISTINCT FROM OLD.collection_id THEN
            -- Decrement book_count for the old collection
            UPDATE collections
            SET book_count = book_count - 1
            WHERE id = OLD.collection_id;

            -- Increment book_count for the new collection
            UPDATE collections
            SET book_count = book_count + 1
            WHERE id = NEW.collection_id;
        END IF;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER book_count_sync
AFTER INSERT OR DELETE OR UPDATE ON books
FOR EACH ROW
EXECUTE FUNCTION manage_book_count();

-- +goose down
DROP TRIGGER IF EXISTS book_count_sync ON books;
DROP FUNCTION IF EXISTS manage_book_count();
CREATE OR REPLACE FUNCTION fn_audit() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    rec_old JSONB;
    rec_new JSONB;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        rec_old = to_jsonb(OLD);
        INSERT INTO audit_log(table_name, operation, row_id, old_data, performed_by, query)
        VALUES (TG_TABLE_NAME, 'D', OLD.id, rec_old, current_setting('app.current_user', true)::UUID, current_query());
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        rec_old = to_jsonb(OLD);
        rec_new = to_jsonb(NEW);
        INSERT INTO audit_log(table_name, operation, row_id, old_data, new_data, performed_by, query)
        VALUES (TG_TABLE_NAME, 'U', NEW.id, rec_old, rec_new, current_setting('app.current_user', true)::UUID, current_query());
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        rec_new = to_jsonb(NEW);
        INSERT INTO audit_log(table_name, operation, row_id, new_data, performed_by, query)
        VALUES (TG_TABLE_NAME, 'I', NEW.id, rec_new, current_setting('app.current_user', true)::UUID, current_query());
        RETURN NEW;
    END IF;
END;
$$;


DO $$
DECLARE t TEXT;
BEGIN
    FOR t IN SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename IN (
        'users','developers','games','game_licenses','transactions','purchases','library','reviews','achievements','user_achievements')
    LOOP
        EXECUTE format('CREATE TRIGGER audit_%I
            AFTER INSERT OR UPDATE OR DELETE ON %I
            FOR EACH ROW EXECUTE FUNCTION fn_audit();', t, t);
    END LOOP;
END;$$;


CREATE OR REPLACE FUNCTION fn_purchase_update_sales() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE games SET sales_count = sales_count + 1 WHERE id = NEW.game_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE games SET sales_count = GREATEST(sales_count - 1, 0) WHERE id = OLD.game_id;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_purchase_sales AFTER INSERT OR DELETE ON purchases FOR EACH ROW EXECUTE FUNCTION fn_purchase_update_sales();

CREATE OR REPLACE FUNCTION fn_update_avg_rating() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    avgr NUMERIC(3,2);
BEGIN
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
        SELECT ROUND(AVG(rating)::numeric,2) INTO avgr FROM reviews WHERE game_id = COALESCE(NEW.game_id, OLD.game_id);
        UPDATE games SET avg_rating = COALESCE(avgr, 0) WHERE id = COALESCE(NEW.game_id, OLD.game_id);
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_reviews_avg AFTER INSERT OR UPDATE OR DELETE ON reviews FOR EACH ROW EXECUTE FUNCTION fn_update_avg_rating();
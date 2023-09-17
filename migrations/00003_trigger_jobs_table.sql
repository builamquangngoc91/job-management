CREATE OR REPLACE FUNCTION jobs_table_inserted_updated()
  RETURNS TRIGGER 
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
    INSERT INTO job_logs(
        job_id,
        name,
        data,
        run_at,
        execute_at,
        status,
        ttl,
        times,
        executed_times,
        note,
        level,
        type,
        logs,
        created_at,
        updated_at,
        deleted_at
    )
    VALUES(
        NEW.id,
        NEW.name,
        NEW.data,
        NEW.run_at,
        NEW.execute_at,
        NEW.status,
        NEW.ttl,
        NEW.times,
        NEW.executed_times,
        NEW.note,
        NEW.level,
        NEW.type,
        NEW.logs,
        NEW.created_at,
        NEW.updated_at,
        NEW.deleted_at
    );

	RETURN NEW;
END;
$$;


CREATE TRIGGER jobs_table_changes
  AFTER UPDATE OR INSERT
  ON jobs
  FOR EACH ROW
  EXECUTE PROCEDURE jobs_table_inserted_updated();
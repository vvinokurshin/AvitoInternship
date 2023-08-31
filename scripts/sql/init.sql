CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS app.users
(
    user_id       bigserial		PRIMARY KEY,
    username      text 			UNIQUE NOT NULL,
    first_name    text			NOT NULL,
    last_name     text 			NOT NULL
);

CREATE TABLE app.segments
(
    segment_id	bigserial	PRIMARY KEY,
    slug 		text 		UNIQUE NOT NULL,
    percent 	int			DEFAULT NULL
);

CREATE TABLE app.users2segments
(
    user_id 		bigint		NOT NULL,
    segment_id		bigint		NOT NULL,
    until 			timestamptz DEFAULT NULL,

    PRIMARY KEY (user_id, segment_id),

    CONSTRAINT fk_u2s_user_id FOREIGN KEY (user_id)
        REFERENCES app.users ON DELETE CASCADE,
    CONSTRAINT fk_u2s_segment_id FOREIGN KEY (segment_id)
        REFERENCES app.segments ON DELETE CASCADE
);

CREATE TABLE app.history
(
    record_id 		bigserial 	PRIMARY KEY,
    user_id 		bigint		NOT NULL,
    segment_slug 	text		NOT NULL,
    operation 		TEXT 		NOT NULL,
    datetime 		timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT fk_history_user_id FOREIGN KEY (user_id)
        REFERENCES app.users ON DELETE CASCADE,

    CONSTRAINT fk_history_segment_slug FOREIGN KEY (segment_slug)
        REFERENCES app.segments(slug) ON DELETE CASCADE
);


-- триггер по добавлению записи истории (операция добавления)
CREATE OR REPLACE FUNCTION history_add()
RETURNS TRIGGER AS
$BODY$
    BEGIN
        INSERT INTO app.history(user_id, segment_slug, operation)
        SELECT NEW.user_id, (SELECT slug FROM app.segments WHERE segment_id = NEW.segment_id), 'ADD';

        RETURN NEW;
    END;
$BODY$
LANGUAGE plpgsql;

-- срабатывает после вставки записей в users2segments
CREATE TRIGGER trig_history_add
AFTER INSERT
ON app.users2segments
FOR EACH ROW
EXECUTE PROCEDURE history_add();


-- триггер по добавлению записи истории (операция удаления)
CREATE OR REPLACE FUNCTION history_del()
RETURNS TRIGGER AS
$BODY$
    BEGIN
        INSERT INTO app.history(user_id, segment_slug, operation)
        SELECT OLD.user_id, (SELECT slug FROM app.segments WHERE segment_id = OLD.segment_id), 'DEL';

        RETURN NEW;
    END;
$BODY$
LANGUAGE plpgsql;

-- срабатывает после удаления записей из users2segments
CREATE TRIGGER trig_history_del
AFTER DELETE
ON app.users2segments
FOR EACH ROW
EXECUTE PROCEDURE history_del();


-- триггер по обновлению времени в истории (при обновлении until в users2segments)
CREATE
OR REPLACE FUNCTION history_datetime_update()
RETURNS TRIGGER AS
$BODY$
    BEGIN
        UPDATE app.history SET datetime = current_timestamp
        WHERE record_id = (
            SELECT record_id
            FROM app.history
            WHERE datetime = (
                SELECT max(datetime)
                FROM app.history
                WHERE user_id = OLD.user_id AND
                        segment_slug = (SELECT slug FROM app.segments WHERE segment_id = OLD.segment_id) AND
                        operation = 'ADD'
            )
        );
        RETURN NEW;
    END;
$BODY$
LANGUAGE plpgsql;

-- срабатывает при обновлении времени в users2segments
CREATE TRIGGER trig_history_datetime_update
AFTER UPDATE
ON app.users2segments
FOR EACH ROW
WHEN (OLD.until IS DISTINCT FROM NEW.until)
EXECUTE PROCEDURE history_datetime_update();

-- удаление старых связей пользователя и сегмента
CREATE OR REPLACE FUNCTION delete_old_accesses()
RETURNS VOID AS
$$
    BEGIN
        DELETE FROM app.users2segments
        WHERE current_timestamp >= until;
    END;
$$
LANGUAGE plpgsql;

-- Drop the GOODS table
DROP TABLE IF EXISTS GOODS;

-- Drop the PROJECTS table
DROP TABLE IF EXISTS PROJECTS;

-- Create the PROJECTS table
CREATE TABLE PROJECTS
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Create the GOODS table
CREATE TABLE GOODS
(
    id          SERIAL PRIMARY KEY,
    project_id  INT          NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    priority    INT                   DEFAULT 0, -- Default value of 0
    removed     BOOLEAN               DEFAULT FALSE,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES PROJECTS (id)
);

-- Create a trigger function to set the priority
CREATE
OR REPLACE FUNCTION set_priority_auto_increment() RETURNS TRIGGER AS $$
BEGIN
    -- Find the maximum priority currently in the table
SELECT COALESCE(MAX(priority), 0) + 1
INTO NEW.priority
FROM GOODS;
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Create a trigger to call the trigger function before insert
CREATE TRIGGER trg_set_priority_auto_increment
    BEFORE INSERT
    ON GOODS
    FOR EACH ROW
    EXECUTE FUNCTION set_priority_auto_increment();

-- Insert a sample record into the PROJECTS table (assuming this is desired behavior based on the provided table description)
INSERT INTO PROJECTS (name)
VALUES ('Запись 1');

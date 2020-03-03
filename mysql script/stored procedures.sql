use sakila;

ALTER TABLE film_category DROP FOREIGN KEY fk_film_category_film;
ALTER TABLE film_category DROP FOREIGN KEY fk_film_category_category;
ALTER TABLE  film_category
    ADD CONSTRAINT fk_film_category_film
        foreign key (film_id) references film (film_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
ALTER TABLE  film_category
    ADD CONSTRAINT fk_film_category_category
        foreign key (category_id) references category (category_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE film_actor DROP FOREIGN KEY fk_film_actor_film;
ALTER TABLE film_actor DROP FOREIGN KEY fk_film_actor_actor;
ALTER TABLE  film_actor
    ADD CONSTRAINT fk_film_actor_film
        foreign key (film_id) references film (film_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
ALTER TABLE  film_actor
    ADD CONSTRAINT fk_film_actor_actor
        foreign key (actor_id) references actor (actor_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;


CREATE DATABASE IF NOT EXISTS sakila_crud;
USE sakila_crud;

CREATE OR REPLACE VIEW ProceduresAndViews AS
    SELECT TABLE_NAME nombre, TABLE_TYPE tipo FROM information_schema.TABLES
    WHERE TABLE_TYPE LIKE 'VIEW' AND TABLE_SCHEMA LIKE 'sakila_crud' AND TABLE_NAME NOT LIKE 'ProceduresAndViews' UNION
    SELECT ROUTINE_NAME nombre, ROUTINE_TYPE tipo FROM information_schema.routines
    WHERE routine_type = 'PROCEDURE' AND routine_schema = 'sakila_crud'
    ORDER BY tipo, nombre;

CREATE OR REPLACE VIEW filmData AS
    SELECT sf.film_id, sf.title, categories, sf.description, sf.release_year, lang.name AS language, olang.name AS original_language,
    sf.rental_duration, sf.rental_rate, sf.length, sf.replacement_cost, sf.rating, sf.special_features,
    sf.last_update, actors
	FROM sakila.film sf LEFT JOIN (
		SELECT film_actor.film_id, group_concat(concat_ws(' ', actor.first_name, actor.last_name) ORDER BY actor.last_name SEPARATOR ',') AS actors
		FROM sakila.film_actor film_actor INNER JOIN sakila.actor actor
        ON film_actor.actor_id = actor.actor_id
		GROUP BY film_actor.film_id
    ) actors ON sf.film_id = actors.film_id
    INNER JOIN sakila.language lang ON sf.language_id = lang.language_id
    LEFT JOIN sakila.language olang ON sf.original_language_id = olang.language_id
    LEFT JOIN (
        SELECT film_category.film_id, group_concat(category.name SEPARATOR ',') AS categories
        FROM sakila.film_category film_category INNER JOIN sakila.category category
        ON film_category.category_id = category.category_id
        GROUP BY film_category.film_id
    ) categories on sf.film_id = categories.film_id;

CREATE OR REPLACE VIEW categoryData AS
    SELECT category.category_id, name, film.film_id, title
    FROM sakila.category LEFT JOIN sakila.film_category ON category.category_id = film_category.category_id
    LEFT JOIN sakila.film ON film_category.film_id = film.film_id;

CREATE OR REPLACE VIEW actorData AS
    SELECT actor.actor_id, actor.first_name, actor.last_name, film.film_id, title
    FROM sakila.actor LEFT JOIN sakila.film_actor ON actor.actor_id = film_actor.actor_id
    LEFT JOIN sakila.film ON film_actor.film_id = film.film_id;

/**********************************************************************************************************************
READ
**********************************************************************************************************************/

DROP PROCEDURE IF EXISTS search_film;
DELIMITER $$
CREATE PROCEDURE search_film(
    p_film_id smallint,
    p_title varchar(255)    
    )
BEGIN		
    SELECT * FROM filmdata
	WHERE	p_film_id = film_id AND
			p_title LIKE title;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS search_films;
DELIMITER $$
CREATE PROCEDURE search_films(
    p_limit int unsigned,
    p_offset int unsigned,
    p_order_by enum ('id','title','year','rental duration','rental price','duration','replacement cost','last update'),
    p_order enum ('asc','desc'),
    p_film_id smallint,
    p_text text,
    p_release_year year(4),
    p_language char(20),
    p_original_language char(20),
    p_rental_duration tinyint,
    p_rental_rate decimal(4, 2),
    p_length smallint,
    p_replacement_cost decimal(5, 2),
    p_rating enum ('G','PG','PG-13','R','NC-17'))
BEGIN
    DECLARE totalLimit INT;
    DECLARE totalOffset INT;
    SET totalLimit =
            IFNULL(IF(p_limit > 0, p_limit, (SELECT COUNT(*) FROM sakila.film)), (SELECT COUNT(*) FROM sakila.film));
    SET totalOffset = IFNULL(p_offset, 0);

    SELECT *
    FROM filmdata
    WHERE (IF(p_film_id IS NULL, TRUE, p_film_id = film_id) AND
           IF(p_text IS NULL, TRUE,
               title LIKE CONCAT('%', p_text, '%') OR description LIKE CONCAT('%', p_text, '%')) AND
            IF(p_release_year IS NULL, TRUE, p_release_year = release_year) AND
           IF(p_language IS NULL, TRUE, p_language LIKE language) AND
           IF(p_original_language IS NULL, TRUE, p_original_language LIKE original_language) AND
           IF(p_rental_duration IS NULL, TRUE, p_rental_duration = rental_duration) AND
           IF(p_rental_rate IS NULL, TRUE, p_rental_rate = rental_rate) AND
           IF(p_length IS NULL, TRUE, p_length = length) AND
           IF(p_replacement_cost IS NULL, TRUE, p_replacement_cost = replacement_cost) AND
           IF(p_rating IS NULL, TRUE, p_rating = rating)
       )
    ORDER BY
        CASE WHEN (p_order_by = 'id' OR p_order_by IS NULL) AND p_order = 'asc' OR p_order IS NULL THEN film_id END ASC,
        CASE WHEN (p_order_by = 'id' OR p_order_by IS NULL) AND p_order = 'desc' THEN film_id END DESC,
        CASE WHEN p_order_by = 'title' AND p_order = 'asc' THEN title END ASC,
        CASE WHEN p_order_by = 'title' AND p_order = 'desc' THEN title END DESC,
        CASE WHEN p_order_by = 'year' AND p_order = 'asc' THEN release_year END ASC,
        CASE WHEN p_order_by = 'year' AND p_order = 'desc' THEN release_year END DESC,
        CASE WHEN p_order_by = 'rental duration' AND p_order = 'asc' THEN rental_duration END ASC,
        CASE WHEN p_order_by = 'rental duration' AND p_order = 'desc' THEN rental_duration END DESC,
        CASE WHEN p_order_by = 'rental price' AND p_order = 'asc' THEN rental_rate END ASC,
        CASE WHEN p_order_by = 'rental price' AND p_order = 'desc' THEN rental_rate END DESC,
        CASE WHEN p_order_by = 'duration' AND p_order = 'asc' THEN length END ASC,
        CASE WHEN p_order_by = 'duration' AND p_order = 'desc' THEN length END DESC,
        CASE WHEN p_order_by = 'replacement cost' AND p_order = 'asc' THEN replacement_cost END ASC,
        CASE WHEN p_order_by = 'replacement cost' AND p_order = 'desc' THEN replacement_cost END DESC,
        CASE WHEN p_order_by = 'last update' AND p_order = 'asc' THEN last_update END ASC,
        CASE WHEN p_order_by = 'last update' AND p_order = 'desc' THEN last_update END DESC,
        title ASC
    LIMIT totalLimit OFFSET totalOffset;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS search_films_count;
DELIMITER $$
CREATE PROCEDURE search_films_count(
    p_film_id smallint,
    p_text text,
    p_release_year year(4),
    p_language char(20),
    p_original_language char(20),
    p_rental_duration tinyint,
    p_rental_rate decimal(4, 2),
    p_length smallint,
    p_replacement_cost decimal(5, 2),
    p_rating enum ('G','PG','PG-13','R','NC-17'))
BEGIN
    SELECT COUNT(film_id) AS film_count
    FROM filmdata
    WHERE (IF(p_film_id IS NULL, TRUE, p_film_id = film_id) AND
           IF(p_text IS NULL, TRUE,
               title LIKE CONCAT('%', p_text, '%') OR description LIKE CONCAT('%', p_text, '%')) AND
            IF(p_release_year IS NULL, TRUE, p_release_year = release_year) AND
           IF(p_language IS NULL, TRUE, p_language LIKE language) AND
           IF(p_original_language IS NULL, TRUE, p_original_language LIKE original_language) AND
           IF(p_rental_duration IS NULL, TRUE, p_rental_duration = rental_duration) AND
           IF(p_rental_rate IS NULL, TRUE, p_rental_rate = rental_rate) AND
           IF(p_length IS NULL, TRUE, p_length = length) AND
           IF(p_replacement_cost IS NULL, TRUE, p_replacement_cost = replacement_cost) AND
           IF(p_rating IS NULL, TRUE, p_rating = rating)
       );
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS search_category;
DELIMITER $$
CREATE PROCEDURE search_category(p_category varchar(25))
BEGIN
    SELECT category_id, name, COUNT(film_id) AS 'film_count'
    FROM categoryData
	WHERE	IF(p_category IS NULL, TRUE, p_category LIKE name)
    GROUP BY category_id, name
    ORDER BY name ASC;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS search_actor;
DELIMITER $$
CREATE PROCEDURE search_actor(p_fname varchar(45), p_lname varchar(45))
BEGIN
    SELECT actor_id, first_name, last_name, group_concat(title ORDER BY title SEPARATOR ',') titles
    FROM actorData
	WHERE IF(p_fname IS NULL, TRUE, p_fname LIKE first_name) AND
	      IF(p_lname IS NULL, TRUE, p_lname LIKE last_name)
    GROUP BY actor_id, first_name
    ORDER BY first_name ASC;
END$$
DELIMITER ;

/**********************************************************************************************************************
CREATE
**********************************************************************************************************************/

DROP PROCEDURE IF EXISTS add_film;
DELIMITER $$
CREATE PROCEDURE add_film(
    p_title varchar(255),
    p_description text,
    p_categories json,
    p_actors json,
    p_release_year year(4),
    p_language char(20),
    p_original_language char(20),
    p_rental_duration tinyint unsigned,
    p_rental_rate decimal(4,2),
    p_length smallint unsigned,
	p_replacement_cost decimal(5,2),
	p_rating enum('G','PG','PG-13','R','NC-17'),
	p_special_features set('Trailers','Commentaries','Deleted Scenes','Behind the Scenes')
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE _category VARCHAR(25);
    DECLARE _actor_firstname VARCHAR(45);
    DECLARE _actor_lastname VARCHAR(45);
    DECLARE lang_id TINYINT;
    DECLARE o_lang_id TINYINT;
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    SET lang_id = (SELECT language_id FROM sakila.language WHERE name LIKE p_language LIMIT 1);
    SET o_lang_id = (SELECT language_id FROM sakila.language WHERE name LIKE p_original_language LIMIT 1);
    IF p_title IS NULL 
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'wrong title'; END IF;
    IF p_language IS NULL OR lang_id IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'wrong language'; END IF;
    START TRANSACTION;
    INSERT INTO sakila.film (title,
							description,
							release_year,
							language_id,
							original_language_id,
                            rental_duration,
                            rental_rate,
                            length,
                            replacement_cost,
                            rating,
                            special_features)
	VALUES (p_title,
			p_description,
            p_release_year,
            lang_id,
            o_lang_id,
            IFNULL(p_rental_duration, DEFAULT(rental_duration)),
            IFNULL(p_rental_rate, DEFAULT(rental_rate)),
            p_length,
            IFNULL(p_replacement_cost, DEFAULT(replacement_cost)),
            IFNULL(p_rating, DEFAULT(rating)),
            p_special_features);

    WHILE i < JSON_LENGTH(p_categories) DO
    SET _category = JSON_UNQUOTE(JSON_EXTRACT(p_categories, CONCAT('$[',i,']')));
    INSERT IGNORE INTO sakila.film_category(film_id, category_id)
        SELECT film_id, category_id FROM sakila.film, sakila.category
        WHERE p_title LIKE title AND
              p_release_year LIKE release_year AND
              p_length LIKE length AND
              category.name LIKE _category;
    SET i = i + 1;
    END WHILE;

    SET i = 0;
    WHILE i < JSON_LENGTH(p_actors) DO
    SET _actor_firstname = JSON_UNQUOTE(JSON_EXTRACT(p_actors, CONCAT('$[',i,'].firstname')));
    SET _actor_lastname = JSON_UNQUOTE(JSON_EXTRACT(p_actors, CONCAT('$[',i,'].lastname')));
    INSERT IGNORE INTO sakila.film_actor(actor_id, film_id)
        SELECT actor_id, film_id FROM sakila.film, sakila.actor
        WHERE p_title LIKE title AND
              p_release_year LIKE release_year AND
              p_length LIKE length AND
              actor.first_name LIKE _actor_firstname AND
              actor.last_name LIKE _actor_lastname;
    SET i = i + 1;
    END WHILE;

    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS add_film_category;
DELIMITER $$
CREATE PROCEDURE add_film_category(
    p_title varchar(255),
    p_category char(20)
)
BEGIN
    DECLARE categoryID TINYINT;
    DECLARE filmID SMALLINT;
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    SET categoryID = (SELECT category_id FROM sakila.category WHERE name LIKE p_category LIMIT 1);
    SET filmID = (SELECT film_id FROM sakila.film WHERE title LIKE p_title LIMIT 1);
    IF p_title IS NULL OR filmID IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'wrong film'; END IF;
    IF p_category IS NULL OR categoryID IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'wrong category'; END IF;
    START TRANSACTION;
    INSERT INTO sakila.film_category (film_id, category_id)
		VALUES (filmID, categoryID);
    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS add_category;
DELIMITER $$
CREATE PROCEDURE add_category(p_name varchar(25))
BEGIN
    INSERT INTO sakila.category (name)
		VALUES (p_name);
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS add_film_actor;
DELIMITER $$
CREATE PROCEDURE add_film_actor(
    p_title varchar(255),
    p_first_name varchar(45),
    p_last_name varchar(45)
)
BEGIN
    DECLARE actorID SMALLINT;
    DECLARE filmID SMALLINT;
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    SET actorID = (SELECT actor_id FROM sakila.actor WHERE first_name LIKE p_first_name AND last_name LIKE p_last_name LIMIT  1);
    SET filmID = (SELECT film_id FROM sakila.film WHERE title LIKE p_title LIMIT 1);
    IF p_title IS NULL OR filmID IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'wrong film'; END IF;
    IF p_first_name IS NULL OR p_last_name IS NULL OR actorID IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'wrong actor'; END IF;
    START TRANSACTION;
    INSERT INTO sakila.film_actor (film_id, actor_id)
		VALUES (filmID, actorID);
    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS add_actor;
DELIMITER $$
CREATE PROCEDURE add_actor(
    p_first_name varchar(45),
    p_last_name varchar(45)
)
BEGIN
    INSERT INTO sakila.actor (first_name, last_name)
		VALUES (p_first_name, p_last_name);
END$$
DELIMITER ;

/**********************************************************************************************************************
UPDATE
**********************************************************************************************************************/

DROP PROCEDURE IF EXISTS update_film;
DELIMITER $$
CREATE PROCEDURE update_film(
	p_film_id smallint,
	p_title varchar(255),
	p_categories json,
	p_actors json,
    p_description text,
    p_release_year year(4),
    p_language char(20),
    p_original_language char(20),
    p_rental_duration tinyint unsigned,
    p_rental_rate decimal(4,2),
    p_length smallint unsigned,
	p_replacement_cost decimal(5,2),
	p_rating enum('G','PG','PG-13','R','NC-17'),
	p_special_features set('Trailers','Commentaries','Deleted Scenes','Behind the Scenes')
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE _category VARCHAR(25);
    DECLARE _actor_firstname VARCHAR(45);
    DECLARE _actor_lastname VARCHAR(45);
    DECLARE lang_id TINYINT;
    DECLARE o_lang_id TINYINT;
    DECLARE categories CURSOR FOR SELECT * FROM sakila.film_category WHERE film_id = p_film_id;
    SET lang_id = (SELECT language_id FROM sakila.language WHERE name LIKE p_language);
    SET o_lang_id = (SELECT language_id FROM sakila.language WHERE name LIKE p_original_language);
    DELETE FROM sakila.film_category WHERE film_id = p_film_id;
    DELETE FROM sakila.film_actor WHERE film_id = p_film_id;
	UPDATE sakila.film
	SET title = IF(p_title NOT LIKE BINARY title, p_title, title),
		description = p_description,
		release_year = p_release_year,
		language_id = IF(lang_id <> language_id, lang_id, language_id),
		original_language_id = o_lang_id,
		rental_duration = IF(p_rental_duration <> rental_duration, p_rental_duration, rental_duration),
		rental_rate = IF(p_rental_rate <> rental_rate, p_rental_rate, rental_rate),
		length = p_length,
		replacement_cost = IF(p_replacement_cost <> replacement_cost, p_replacement_cost, replacement_cost),
		rating = p_rating,
		special_features = p_special_features,
        last_update = DEFAULT
	WHERE film_id = p_film_id;
    WHILE i < JSON_LENGTH(p_categories) DO
    SET _category = JSON_UNQUOTE(JSON_EXTRACT(p_categories, CONCAT('$[',i,']')));
    INSERT IGNORE INTO sakila.film_category(film_id, category_id)
        SELECT film_id, category_id FROM sakila.film, sakila.category
        WHERE p_film_id = film_id AND
              category.name LIKE _category;
    SET i = i + 1;
    END WHILE;

    SET i = 0;
    WHILE i < JSON_LENGTH(p_actors) DO
    SET _actor_firstname = JSON_UNQUOTE(JSON_EXTRACT(p_actors, CONCAT('$[',i,'].firstname')));
    SET _actor_lastname = JSON_UNQUOTE(JSON_EXTRACT(p_actors, CONCAT('$[',i,'].lastname')));
    INSERT IGNORE INTO sakila.film_actor(actor_id, film_id)
        SELECT actor_id, film_id FROM sakila.film, sakila.actor
        WHERE p_film_id = film_id AND
              actor.first_name LIKE _actor_firstname AND
              actor.last_name LIKE _actor_lastname;
    SET i = i + 1;
    END WHILE;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS update_category;
DELIMITER $$
CREATE PROCEDURE update_category(
	p_category_id tinyint,
	p_name varchar(25)
)
BEGIN
	UPDATE sakila.category
	SET name = p_name
	WHERE category_id = p_category_id;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS update_actor;
DELIMITER $$
CREATE PROCEDURE update_actor(
	p_actor_id smallint,
	p_first_name varchar(45),
	p_last_name varchar(45)
)
BEGIN
	UPDATE sakila.actor
	SET first_name = p_first_name,
	    last_name = p_last_name
	WHERE actor_id = p_actor_id;
END$$
DELIMITER ;


/**********************************************************************************************************************
DELETE
**********************************************************************************************************************/

DROP PROCEDURE IF EXISTS remove_film;
DELIMITER $$
CREATE PROCEDURE remove_film(p_film_id smallint)
BEGIN
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    IF p_film_id IS NULL OR (SELECT COUNT(*) FROM sakila.film WHERE film_id = p_film_id) = 0
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'id not found', MYSQL_ERRNO = 1; END IF;
    START TRANSACTION;    	    
    DELETE FROM sakila.film WHERE film_id = p_film_id;
    /* IN CASE THAT WE WANT TO MODIFY THE AUTO INCREMENT FIELD:*/
	-- set @val  = (SELECT MAX(film_id) + 1 FROM sakila.film);
	-- SET @sql = CONCAT('ALTER TABLE `sakila`.film AUTO_INCREMENT = ', @val);
	-- PREPARE st FROM @sql;
	-- EXECUTE st;
    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS remove_category_film;
DELIMITER $$
CREATE PROCEDURE remove_category_film(p_film_id smallint, p_category varchar(25))
BEGIN
    DECLARE categoryID TINYINT;
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    SET categoryID = (SELECT category_id FROM sakila.category WHERE name LIKE p_category);
    IF p_film_id IS NULL OR (SELECT COUNT(*) FROM sakila.film WHERE film_id = p_film_id) = 0
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'id not found', MYSQL_ERRNO = 1; END IF;
	IF p_category IS NULL OR categoryID IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'category not found', MYSQL_ERRNO = 1; END IF;
    START TRANSACTION;
    DELETE FROM sakila.film_category WHERE film_id = p_film_id AND category_id = categoryID;
    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS remove_category;
DELIMITER $$
CREATE PROCEDURE remove_category(p_category_id tinyint)
BEGIN
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    IF p_category_id IS NULL OR (SELECT COUNT(*) FROM sakila.category WHERE category_id = p_category_id) = 0
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'id not found', MYSQL_ERRNO = 1; END IF;
    START TRANSACTION;
    DELETE FROM sakila.category WHERE category_id = p_category_id;
    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS remove_actor_film;
DELIMITER $$
CREATE PROCEDURE remove_actor_film(p_film_id smallint, p_name varchar(90))
BEGIN
    DECLARE actorID SMALLINT;
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    SET actorID = (SELECT actor_id FROM sakila.actor WHERE concat_ws(' ', first_name, last_name) LIKE p_name);
    IF p_film_id IS NULL OR (SELECT COUNT(*) FROM sakila.film WHERE film_id = p_film_id) = 0
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'id not found', MYSQL_ERRNO = 1; END IF;
	IF p_name IS NULL OR actorID IS NULL
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'actor not found', MYSQL_ERRNO = 1; END IF;
    START TRANSACTION;
    DELETE FROM sakila.film_actor WHERE film_id = p_film_id AND actor_id = actorID;
    COMMIT;
END$$
DELIMITER ;

DROP PROCEDURE IF EXISTS remove_actor;
DELIMITER $$
CREATE PROCEDURE remove_actor(p_actor_id smallint)
BEGIN
	DECLARE err_arg CONDITION FOR SQLSTATE 'ERROR'; -- sqlstate must be a size 5 string
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    IF p_actor_id IS NULL OR (SELECT COUNT(*) FROM sakila.actor WHERE actor_id = p_actor_id) = 0
		THEN SIGNAL err_arg SET MESSAGE_TEXT = 'id not found', MYSQL_ERRNO = 1; END IF;
    START TRANSACTION;
    DELETE FROM sakila.actor WHERE actor_id = p_actor_id;
    COMMIT;
END$$
DELIMITER ;

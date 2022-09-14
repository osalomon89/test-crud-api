CREATE TABLE IF NOT EXISTS photos (
	id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	path longtext,
	item_id bigint(20) unsigned DEFAULT NULL,
	created_at datetime(3) DEFAULT NULL,
	updated_at datetime(3) DEFAULT NULL,
	PRIMARY KEY (id),
	KEY fk_items_photos (item_id),
	CONSTRAINT fk_items_photos FOREIGN KEY (item_id) REFERENCES items (id)
);
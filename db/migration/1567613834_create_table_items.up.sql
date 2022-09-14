CREATE TABLE IF NOT EXISTS items (
	id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	code varchar(191) DEFAULT NULL,
	title longtext,
	description longtext,
	price bigint(20) DEFAULT NULL,
	stock bigint(20) DEFAULT NULL,
	item_type longtext,
	leader tinyint(1) DEFAULT NULL,
	leader_level longtext,
	status longtext,
	created_at datetime(3) DEFAULT NULL,
	updated_at datetime(3) DEFAULT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY code (code)
);
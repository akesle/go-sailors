CREATE TABLE sailors.Sailor (
	FirstName varchar(100) NOT NULL,
	LastName varchar(100) NOT NULL,
	Age INT NOT NULL,
	ID BIGINT NOT NULL AUTO_INCREMENT,
	CONSTRAINT Sailor_PK PRIMARY KEY (ID)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;
CREATE INDEX Sailor_LastName_FirstName_Age_IDX USING BTREE ON sailors.Sailor (LastName,FirstName,Age) ;


GRANT Alter ON sailors.* TO 'api_user'@'%' ;
GRANT Create ON sailors.* TO 'api_user'@'%' ;
GRANT Create view ON sailors.* TO 'api_user'@'%' ;
GRANT Delete ON sailors.* TO 'api_user'@'%' ;
GRANT Drop ON sailors.* TO 'api_user'@'%' ;
GRANT Grant option ON sailors.* TO 'api_user'@'%' ;
GRANT Index ON sailors.* TO 'api_user'@'%' ;
GRANT Insert ON sailors.* TO 'api_user'@'%' ;
GRANT References ON sailors.* TO 'api_user'@'%' ;
GRANT Select ON sailors.* TO 'api_user'@'%' ;
GRANT Show view ON sailors.* TO 'api_user'@'%' ;
GRANT Trigger ON sailors.* TO 'api_user'@'%' ;
GRANT Update ON sailors.* TO 'api_user'@'%' ;
FLUSH PRIVILEGES;
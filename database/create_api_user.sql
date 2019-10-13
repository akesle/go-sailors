
CREATE USER 'api_user'@'%' ;
ALTER USER 'api_user'@'%'
IDENTIFIED BY 'example' ;
GRANT Create tablespace ON *.* TO 'api_user'@'%' ;
GRANT Show databases ON *.* TO 'api_user'@'%' ;
GRANT Reload ON *.* TO 'api_user'@'%' ;
GRANT Event ON *.* TO 'api_user'@'%' ;
GRANT File ON *.* TO 'api_user'@'%' ;
GRANT Process ON *.* TO 'api_user'@'%' ;
FLUSH PRIVILEGES;

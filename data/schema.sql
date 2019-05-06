CREATE TABLE IF NOT EXISTS players (
                                       playerId   VARCHAR(64) NOT NULL,
                                       birthDate DATE NOT NULL,
                                       birthCountry   VARCHAR(64) NOT NULL,
                                       birthState VARCHAR(64) NOT NULL,
                                       birthCity VARCHAR(64) NOT NULL,
                                       deathDate  DATE DEFAULT NULL,
                                       deathCountry   VARCHAR(64),
                                       deathState VARCHAR(64),
                                       deathCity VARCHAR(64),
                                       nameFirst VARCHAR(64) NOT NULL,
                                       nameLast VARCHAR(64) NOT NULL,
                                       nameGiven VARCHAR(64) NOT NULL,
                                       weight  INT NOT NULL,
                                       height  INT NOT NULL,
                                       bats  ENUM('R', 'L', 'B') NOT NULL,
                                       throws ENUM('R', 'L', 'B') NOT NULL,
                                       debut DATE NOT NULL,
                                       finalGame   DATE,
                                       retroID  VARCHAR(64) NOT NULL,
                                       bbrefID VARCHAR(64) NOT NULL,
                                       PRIMARY KEY (playerId)
)  ENGINE=INNODB;
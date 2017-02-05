CREATE DATABASE `RelayNovel` /*!40100 DEFAULT CHARACTER SET utf8 */;

-- マスター系のテーブルを作成
CREATE TABLE `RelayNovel`.`fan` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '読者ID',
  `name` varchar(45) NOT NULL COMMENT '読者名',
  `create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登録日時',
  `update` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='読者管理マスター';

CREATE TABLE `RelayNovel`.`genre` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ジャンルID',
  `genre` varchar(45) NOT NULL COMMENT 'ジャンル',
  `create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'デリートフラグ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `genre_UNIQUE` (`genre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='ジャンル管理マスター';

CREATE TABLE `RelayNovel`.`tag` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'タグID',
  `tag` varchar(10) NOT NULL COMMENT 'タグ',
  `create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'デリートフラグ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `tag_UNIQUE` (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='タグ管理マスター';

-- マスター系のテーブルに依存性のあるテーブルを作成
CREATE TABLE `RelayNovel`.`novelist` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '小説家ID',
  `name` varchar(45) NOT NULL COMMENT '小説家名',
  `create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `update` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='投稿者管理マスター';

CREATE TABLE `RelayNovel`.`novel` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ノベルID',
  `title` varchar(45) NOT NULL COMMENT 'タイトル',
  `genre-id` int(11) DEFAULT NULL COMMENT 'ジャンルID',
  `summary` varchar(200) NOT NULL COMMENT 'あらすじ',
  `put-gently-pen` tinyint(1) NOT NULL DEFAULT '0' COMMENT '完結フラグ',
  `relay-limmit` int(11) NOT NULL DEFAULT '0' COMMENT '作品が完結するまでのリレー回数',
  `novelist-id` int(11) NOT NULL COMMENT '作者ID',
  `first-edition` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '初回投稿日',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `genre-id_idx` (`genre-id`),
  KEY `tag-id_idx` (`id`),
  KEY `novelist-id_idx` (`novelist-id`),
  CONSTRAINT `novel-genre-id` FOREIGN KEY (`genre-id`) REFERENCES `genre` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `novel-novelist-id` FOREIGN KEY (`novelist-id`) REFERENCES `novelist` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='小説概要';

CREATE TABLE `RelayNovel`.`sentence` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '投稿文章ID',
  `novel-id` int(11) NOT NULL COMMENT 'ノベルID',
  `novelist-id` int(11) NOT NULL COMMENT '投稿者ID',
  `first-line` varchar(50) DEFAULT NULL COMMENT '一行目',
  `second-line` varchar(50) DEFAULT NULL COMMENT '二行目',
  `revision` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '投稿日',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `novel-id_idx` (`novel-id`),
  KEY `sentence-novelist-id_idx` (`novelist-id`),
  CONSTRAINT `sentence-novel-id` FOREIGN KEY (`novel-id`) REFERENCES `novel` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `sentence-novelist-id` FOREIGN KEY (`novelist-id`) REFERENCES `novelist` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='小説本文';

CREATE TABLE `RelayNovel`.`tag2novel` (
  `tag-id` int(11) NOT NULL COMMENT 'タグID',
  `novel-id` int(11) NOT NULL COMMENT 'ノベルID',
  PRIMARY KEY (`tag-id`,`novel-id`),
  KEY `novel-id_idx` (`novel-id`),
  CONSTRAINT `tag2novel-novel-id` FOREIGN KEY (`novel-id`) REFERENCES `novel` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `tag2novel-tag-id` FOREIGN KEY (`tag-id`) REFERENCES `tag` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='タグと小説概要を繋ぐ中間テーブル';

CREATE TABLE `RelayNovel`.`review` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'レビューID',
  `novel-id` INT NOT NULL COMMENT 'ノベルID',
  `fan-id` INT NOT NULL COMMENT '読者ID',
  `favorite` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'お気に入りフラグ',
  `funny` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '笑えるフラグ',
  `interesting` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '面白いフラグ',
  `sad` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '悲しいフラグ',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  INDEX `review-novel-id_idx` (`novel-id` ASC),
  INDEX `review-fan-id_idx` (`fan-id` ASC),
  CONSTRAINT `review-novel-id`
    FOREIGN KEY (`novel-id`)
    REFERENCES `RelayNovel`.`novel` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `review-fan-id`
    FOREIGN KEY (`fan-id`)
    REFERENCES `RelayNovel`.`fan` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
COMMENT = 'レビュー';

-- create test data
-- genre
INSERT INTO `RelayNovel`.`genre`(`genre`) values('genre');

-- tag
INSERT INTO `RelayNovel`.`tag`(`tag`) values('tag');

-- novelist
INSERT INTO `RelayNovel`.`novelist`(`name`) values('novelist');

-- fan
INSERT INTO `RelayNovel`.`fan`(`name`) values('fan');

-- novel
INSERT INTO `RelayNovel`.`novel`(`title`, `genre-id`, `summary`, `relay-limmit`, `novelist-id`) values('title', 1, 'summary', 10, 1);

-- sentence
INSERT INTO `RelayNovel`.`sentence`(`novel-id`, `novelist-id`, `first-line`, `second-line`) values(1, 1, 'first-line', 'second-line');

-- review 
INSERT INTO `RelayNovel`.`review`(`novel-id`, `fan-id`, `favorite`, `funny`, `interesting`, `sad`) values(1, 1, 1, 0, 0, 0);

-- tag2novel
INSERT INTO `RelayNovel`.`tag2novel`(`tag-id`, `novel-id`) values(1, 1);

CREATE TABLE `project` (
	`id` char(30) NOT NULL,
	`time_created` timestamp(3) NOT NULL DEFAULT (now()),
	`time_updated` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	`time_deleted` timestamp(3),
	`name` varchar(255) NOT NULL,
	`year` varchar(10) NOT NULL,
	`technologies` text NOT NULL,
	`description` text NOT NULL,
	`order` int,
	CONSTRAINT `project_id` PRIMARY KEY(`id`)
);

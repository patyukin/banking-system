CREATE KEYSPACE IF NOT EXISTS notifier WITH replication =
	{'class':'SimpleStrategy','replication_factor':'1'};

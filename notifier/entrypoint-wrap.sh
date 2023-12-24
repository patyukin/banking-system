#!/bin/bash

CREATE KEYSPACE IF NOT EXISTS notifier WITH replication =
{'class':'SimpleStrategy','replication_factor':'1'};

CREATE TABLE IF NOT EXISTS notifier.model
(
  timestamp timestamp,
  id varchar,
  PRIMARY KEY (id, timestamp)

) WITH CLUSTERING ORDER BY (timestamp DESC);

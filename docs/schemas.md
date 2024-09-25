# Schemas

Overview of tables in the database.

## Users

Table to store static info about users.

| user_id | region | lat | lon | joined_at | 
|--------|--------|-----|-----|--------|
|uuid    |varchar |real |real |timestamp|

This will be joined with the transaction table.

## Interactions

Every interaction will be logged here.

|event_id| user_id | event_type | sent_at |payload|
|--------|--------|-----------|-----------|-------|
|bigint  |uuid    | varchar   | timestamp |int

TODO: refactor EventType info enum




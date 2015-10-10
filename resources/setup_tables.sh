#!/bin/bash

sqlite3 bible.db < setup_tables.sql


for filename in HCSB/*.sql NIV/*.sql ESV/*.sql; do
	echo $filename
	sqlite3 bible.db < "$filename"
done
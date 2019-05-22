-- longterm goal: replace the 'identity_id' column with a 'username'
-- for now, just add the new `username` column as NULLABLE until we can safely
-- remove the `identity_id` column, which is not the case yet as we have multiple pods and
-- in case we need to rollback. 
ALTER TABLE audit_log add COLUMN username text;
-- also, make the `identity` column nullable from now on
ALTER TABLE audit_log ALTER COLUMN identity_id drop not null;

-- index to username on audit_log, to easily retrieve events for a user given her username
CREATE INDEX ix_auditlog_username ON audit_log USING btree (username);
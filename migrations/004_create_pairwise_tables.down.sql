-- Drop indexes first
DROP INDEX IF EXISTS idx_attendee_votes_attendee_id;
DROP INDEX IF EXISTS idx_attendee_votes_comparison_id;
DROP INDEX IF EXISTS idx_pairwise_comparisons_consensus;
DROP INDEX IF EXISTS idx_pairwise_comparisons_session_id;
DROP INDEX IF EXISTS idx_pairwise_sessions_status;
DROP INDEX IF EXISTS idx_pairwise_sessions_project_id;

-- Drop tables in reverse order due to foreign key dependencies
DROP TABLE IF EXISTS attendee_votes;
DROP TABLE IF EXISTS pairwise_comparisons;
DROP TABLE IF EXISTS pairwise_sessions;
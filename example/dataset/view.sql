/*[bqv:TEST]
- test_count
  - mock:
    - table: bigquery-public-data.stackoverflow.posts_answers
      sql: |
        SELECT 1 AS owner_user_id
        FROM UNNEST([1,2,3,4,5]) AS owner_user_id
  - target: |
    SELECT owner_user_id
    FROM dataset.view
    GROUP BY owner_user_id
  - expect: |
    [
      {
        "owner_user_id": "1"
      },
      {
        "owner_user_id": "2"
      },
      {
        "owner_user_id": "3"
      },
      {
        "owner_user_id": "4"
      },
      {
        "owner_user_id": "5"
      }
    ]
*/
SELECT owner_user_id, COUNT(1)
FROM `bigquery-public-data.stackoverflow.posts_answers`
GROUP BY owner_user_id

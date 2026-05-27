BEGIN;

WITH normalized AS (
  SELECT
    id,
    LOWER(BTRIM(COALESCE(category, ''))) AS raw_category,
    LOWER(COALESCE(title, '')) AS raw_title
  FROM jobs
)
UPDATE jobs j
SET category = CASE
  WHEN n.raw_category IN ('engineering', 'eng', 'developer', 'software', 'technology')
    OR n.raw_title ~* '(engineer|developer|software|backend|frontend|full stack|devops|platform|sre|mobile|ios|android)'
    THEN 'Engineering'

  WHEN n.raw_category IN ('data', 'analytics')
    OR n.raw_title ~* '(data|analytics|scientist|machine learning|ml|ai)'
    THEN 'Data'

  WHEN n.raw_category IN ('product')
    OR n.raw_title ~* '(product manager|product owner|product)'
    THEN 'Product'

  WHEN n.raw_category IN ('design')
    OR n.raw_title ~* '(designer|ux|ui|product design|visual design)'
    THEN 'Design'

  WHEN n.raw_category IN ('marketing', 'growth')
    OR n.raw_title ~* '(marketing|growth|seo|content|brand)'
    THEN 'Marketing'

  WHEN n.raw_category IN ('sales', 'business development')
    OR n.raw_title ~* '(sales|account executive|business development)'
    THEN 'Sales'

  WHEN n.raw_category IN ('hr', 'people', 'people ops', 'human resources', 'recruiting')
    OR n.raw_title ~* '(recruiter|talent|people operations|human resources|hr)'
    THEN 'People'

  WHEN n.raw_category IN ('finance', 'accounting')
    OR n.raw_title ~* '(finance|accounting|financial)'
    THEN 'Finance'

  WHEN n.raw_category IN ('security')
    OR n.raw_title ~* '(security|cybersecurity|infosec|application security)'
    THEN 'Security'

  WHEN n.raw_category IN ('operations', 'ops')
    OR n.raw_title ~* '(operations|program manager|technical program manager)'
    THEN 'Operations'

  WHEN n.raw_category IN ('customer success', 'support')
    OR n.raw_title ~* '(customer success|support|customer support|success manager)'
    THEN 'Customer Success'

  WHEN n.raw_category IN ('legal', 'compliance')
    OR n.raw_title ~* '(legal|counsel|paralegal|compliance)'
    THEN 'Legal'

  ELSE 'Other'
END
FROM normalized n
WHERE j.id = n.id;

COMMIT;

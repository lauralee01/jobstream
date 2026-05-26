INSERT INTO companies (name, slug, provider, enabled) VALUES
('Figma', 'figma', 'greenhouse', true),
('Stripe', 'stripe', 'greenhouse', true),
('Vercel', 'vercel', 'greenhouse', true),
('Mercury', 'mercury', 'greenhouse', true),
('Discord', 'discord', 'greenhouse', true),
('Pinterest', 'pinterest', 'greenhouse', true)
ON CONFLICT (slug) DO NOTHING;
